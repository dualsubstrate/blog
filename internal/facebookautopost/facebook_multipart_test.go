package facebookautopost

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPostPhotoFromFileMultipartShape(t *testing.T) {
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "strip.webp")
	if err := os.WriteFile(imgPath, []byte("fake-webp-bytes"), 0o644); err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method: %s", r.Method)
		}
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("content-type: %s", r.Header.Get("Content-Type"))
		}
		if got := r.URL.Query().Get("access_token"); got != "tok-test" {
			t.Errorf("access_token query: %q", got)
		}
		_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			t.Fatal(err)
		}
		boundary, ok := params["boundary"]
		if !ok {
			t.Fatal("no boundary")
		}
		mr := multipart.NewReader(r.Body, boundary)
		var sawPublished, sawCaption, sawSource bool
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			slurp, _ := io.ReadAll(p)
			switch p.FormName() {
			case "published":
				sawPublished = true
				if string(slurp) != "true" {
					t.Errorf("published: %q", slurp)
				}
			case "caption":
				sawCaption = true
				if string(slurp) != "Line1\n\nhttps://blog.dualsubstrate.com/x/" {
					t.Errorf("caption: %q", slurp)
				}
			case "source":
				sawSource = true
				if string(slurp) != "fake-webp-bytes" {
					t.Errorf("source body: %q", slurp)
				}
				if fn := p.FileName(); fn != "strip.webp" {
					t.Errorf("filename: %q", fn)
				}
			default:
				t.Errorf("unexpected part %q", p.FormName())
			}
		}
		if !sawPublished || !sawCaption || !sawSource {
			t.Fatalf("parts published=%v caption=%v source=%v", sawPublished, sawCaption, sawSource)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"123"}`))
	}))
	t.Cleanup(srv.Close)

	c := &Client{
		HTTP:    srv.Client(),
		BaseURL: strings.TrimSuffix(srv.URL, "/"),
	}
	ctx := testCtx(t)
	if err := c.PostPhotoFromFile(ctx, "page1", "tok-test", imgPath, "Line1\n\nhttps://blog.dualsubstrate.com/x/"); err != nil {
		t.Fatal(err)
	}
}
