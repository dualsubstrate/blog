package substackhtml

import (
	"strings"
	"testing"
)

func TestDetectProjectAboutKindsCubeCows(t *testing.T) {
	meta := FrontMatterMeta{Categories: []string{"Cognitive-Memetics", "Cube-Cows"}}
	p := "content/cognitive-memetics/cows/2026-02-26-cow-w01/index.md"
	got := DetectProjectAboutKinds(p, meta)
	if len(got) != 1 || got[0] != ProjectAboutCubeCows {
		t.Fatalf("got %#v", got)
	}
}

func TestDetectProjectAboutKindsSayingsPorEstasThenReptilocracy(t *testing.T) {
	meta := FrontMatterMeta{Categories: []string{"Por-Estas-Calles", "Reptilocracy"}}
	p := "content/cognitive-memetics/sayings/2026-04-27-saying-16/index.md"
	got := DetectProjectAboutKinds(p, meta)
	if len(got) != 2 || got[0] != ProjectAboutPorEstasCallesSayings || got[1] != ProjectAboutReptilocracy {
		t.Fatalf("got %#v", got)
	}
}

func TestDetectProjectAboutKindsReptilocracyPath(t *testing.T) {
	meta := FrontMatterMeta{Categories: []string{"Reptilocracy"}}
	p := "content/cognitive-memetics/reptilocracy/2026-04-12-not-in-our-term/index.md"
	got := DetectProjectAboutKinds(p, meta)
	if len(got) != 1 || got[0] != ProjectAboutReptilocracy {
		t.Fatalf("got %#v", got)
	}
}

func TestDetectProjectAboutKindsPawtropolisPath(t *testing.T) {
	meta := FrontMatterMeta{Categories: []string{"Pawtropolis-Under-Fire"}}
	p := "content/cognitive-memetics/pawtropolis/2026-05-13-01-i-want-this-to-be-over/index.md"
	got := DetectProjectAboutKinds(p, meta)
	if len(got) != 1 || got[0] != ProjectAboutPawtropolis {
		t.Fatalf("got %#v", got)
	}
}

func TestAppendCognitiveMemeticsProjectAboutHTMLPawtropolisEN(t *testing.T) {
	meta := FrontMatterMeta{Categories: []string{"Cognitive-Memetics", "Pawtropolis-Under-Fire"}}
	p := "content/cognitive-memetics/pawtropolis/2026-05-13-01-i-want-this-to-be-over/index.md"
	out, err := AppendCognitiveMemeticsProjectAboutHTML("<p>body</p>", p, meta, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "quiet protest") {
		t.Fatalf("missing pawtropolis framing: %s", out)
	}
	if !strings.Contains(out, "<blockquote>") || !strings.Contains(out, "BUT WHY:") {
		t.Fatalf("expected blockquote heading: %s", out)
	}
}

func TestDetectProjectAboutKindsNonCognitive(t *testing.T) {
	meta := FrontMatterMeta{Categories: []string{"Cube-Cows"}}
	p := "content/other/2026/x/index.md"
	if got := DetectProjectAboutKinds(p, meta); len(got) != 0 {
		t.Fatalf("got %#v", got)
	}
}

func TestAppendCognitiveMemeticsProjectAboutHTMLCubeCowsEN(t *testing.T) {
	meta := FrontMatterMeta{Categories: []string{"Cube-Cows", "Cognitive-Memetics"}}
	p := "content/cognitive-memetics/cows/2026-02-26-cow-w01/index.md"
	out, err := AppendCognitiveMemeticsProjectAboutHTML("<p>body</p>", p, meta, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Tales from the Cube Farm") {
		t.Fatalf("missing cube framing: %s", out)
	}
	if !strings.Contains(out, "<blockquote>") || !strings.Contains(out, "But why:") {
		t.Fatalf("expected blockquote heading: %s", out)
	}
}

func TestAppendCognitiveMemeticsProjectAboutHTMLSayingsES(t *testing.T) {
	meta := FrontMatterMeta{Categories: []string{"Por-Estas-Calles"}}
	p := "content/cognitive-memetics/sayings/2026-04-27-saying-16/index.es.md"
	out, err := AppendCognitiveMemeticsProjectAboutHTML("<p>x</p>", p, meta, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "cronómetro cultural") {
		t.Fatalf("missing Spanish sayings copy: %s", out)
	}
}
