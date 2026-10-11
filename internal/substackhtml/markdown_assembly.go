package substackhtml

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed templates/lead_assembly.tmpl
var leadAssemblyFS embed.FS

var leadAssembly *template.Template

func init() {
	b, err := leadAssemblyFS.ReadFile("templates/lead_assembly.tmpl")
	if err != nil {
		panic("substackhtml: read lead assembly template: " + err.Error())
	}
	leadAssembly, err = template.New("lead_assembly").Parse(string(b))
	if err != nil {
		panic("substackhtml: parse lead assembly template: " + err.Error())
	}
}

type featuredImageTmplData struct {
	ImageURL        string
	YouTubeWatchURL string
	WatchOn         string
}

func featuredImageTmplDataFromMeta(meta FrontMatterMeta, opt Options) featuredImageTmplData {
	spanish := isSpanishLeadLocale(meta, opt)
	imgURL := strings.TrimSpace(meta.ImageURL)
	ytID := strings.TrimSpace(meta.YouTubeID)
	var ytWatchURL string
	if ytID != "" {
		ytWatchURL = "https://www.youtube.com/watch?v=" + ytID
	}
	watchOn := "Watch on "
	if spanish {
		watchOn = "Ver en "
	}
	return featuredImageTmplData{
		ImageURL:        imgURL,
		YouTubeWatchURL: ytWatchURL,
		WatchOn:         watchOn,
	}
}

func renderFeaturedImageMarkdown(meta FrontMatterMeta, opt Options) (string, error) {
	var buf bytes.Buffer
	if err := leadAssembly.ExecuteTemplate(&buf, "featured_image", featuredImageTmplDataFromMeta(meta, opt)); err != nil {
		return "", fmt.Errorf("featured_image template: %w", err)
	}
	return buf.String(), nil
}

func buildLeadMarkdown(meta FrontMatterMeta, opt Options) ([]byte, error) {
	feat, err := renderFeaturedImageMarkdown(meta, opt)
	if err != nil {
		return nil, err
	}
	isVideo := strings.EqualFold(strings.TrimSpace(meta.Type), "video")
	if isVideo {
		var buf bytes.Buffer
		data := struct {
			Spanish          bool
			HasDescription   bool
			Description      string
			HasSoWhat        bool
			SoWhat           string
			FeaturedMarkdown string
		}{
			Spanish:          isSpanishLeadLocale(meta, opt),
			HasDescription:   strings.TrimSpace(meta.Description) != "",
			Description:      meta.Description,
			HasSoWhat:        strings.TrimSpace(meta.SoWhat) != "",
			SoWhat:           meta.SoWhat,
			FeaturedMarkdown: feat,
		}
		if err := leadAssembly.ExecuteTemplate(&buf, "video_lead", data); err != nil {
			return nil, fmt.Errorf("video_lead template: %w", err)
		}
		return bytes.TrimSpace(buf.Bytes()), nil
	}
	var buf bytes.Buffer
	data := struct {
		HasDescription   bool
		Description      string
		HasSoWhat        bool
		SoWhat           string
		FeaturedMarkdown string
	}{
		HasDescription:   strings.TrimSpace(meta.Description) != "",
		Description:      meta.Description,
		HasSoWhat:        strings.TrimSpace(meta.SoWhat) != "",
		SoWhat:           meta.SoWhat,
		FeaturedMarkdown: feat,
	}
	if err := leadAssembly.ExecuteTemplate(&buf, "plain_lead", data); err != nil {
		return nil, fmt.Errorf("plain_lead template: %w", err)
	}
	return bytes.TrimSpace(buf.Bytes()), nil
}

func assembleMarkdownByType(meta FrontMatterMeta, body []byte, opt Options) ([]byte, error) {
	t := strings.ToLower(strings.TrimSpace(meta.Type))
	switch t {
	case "video":
		lead, err := buildLeadMarkdown(meta, opt)
		if err != nil {
			return nil, err
		}
		if len(bytes.TrimSpace(lead)) == 0 {
			return body, nil
		}
		return append(append(bytes.TrimSpace(lead), []byte("\n\n")...), body...), nil
	case "claims":
		feat, err := renderFeaturedImageMarkdown(meta, opt)
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		data := struct {
			HasClaim         bool
			Description      string
			FeaturedMarkdown string
		}{
			HasClaim:         strings.TrimSpace(meta.Description) != "",
			Description:      meta.Description,
			FeaturedMarkdown: feat,
		}
		if err := leadAssembly.ExecuteTemplate(&buf, "claims_prefix", data); err != nil {
			return nil, fmt.Errorf("claims_prefix template: %w", err)
		}
		out := append([]byte(strings.TrimSpace(buf.String())+"\n\n"), body...)
		out = bytes.TrimSpace(out)
		if strings.TrimSpace(meta.Grounding) != "" {
			var g bytes.Buffer
			if err := leadAssembly.ExecuteTemplate(&g, "claims_grounding_suffix", struct {
				Grounding string
			}{Grounding: strings.TrimSpace(meta.Grounding)}); err != nil {
				return nil, fmt.Errorf("claims_grounding_suffix template: %w", err)
			}
			out = append(out, g.Bytes()...)
		}
		return bytes.TrimSpace(out), nil
	case "sayings":
		feat, err := renderFeaturedImageMarkdown(meta, opt)
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		data := struct {
			HasTeaser        bool
			Description      string
			FeaturedMarkdown string
			HasTLDR          bool
			TLDR             string
			HasContext       bool
			Fluff            string
		}{
			HasTeaser:        strings.TrimSpace(meta.Description) != "",
			Description:      meta.Description,
			FeaturedMarkdown: feat,
			HasTLDR:          strings.TrimSpace(meta.TLDR) != "",
			TLDR:             meta.TLDR,
			HasContext:       strings.TrimSpace(meta.Fluff) != "",
			Fluff:            meta.Fluff,
		}
		if err := leadAssembly.ExecuteTemplate(&buf, "sayings_prefix", data); err != nil {
			return nil, fmt.Errorf("sayings_prefix template: %w", err)
		}
		prefix := strings.TrimSpace(buf.String())
		out := bytes.TrimSpace(append([]byte(prefix+"\n\n"), body...))
		if len(out) == 0 {
			return out, nil
		}
		if len(bytes.TrimSpace(body)) > 0 {
			if len(prefix) > 0 {
				joined := append([]byte(prefix), []byte("\n\n## 📄 Article\n\n")...)
				joined = append(joined, body...)
				return bytes.TrimSpace(joined), nil
			}
			joined := append([]byte("## 📄 Article\n\n"), body...)
			return bytes.TrimSpace(joined), nil
		}
		return out, nil
	case "panel", "scene", "cows":
		feat, err := renderFeaturedImageMarkdown(meta, opt)
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		data := struct {
			HasTeaser        bool
			Description      string
			FeaturedMarkdown string
		}{
			HasTeaser:        strings.TrimSpace(meta.Description) != "",
			Description:      meta.Description,
			FeaturedMarkdown: feat,
		}
		if err := leadAssembly.ExecuteTemplate(&buf, "panel_prefix", data); err != nil {
			return nil, fmt.Errorf("panel_prefix template: %w", err)
		}
		lead := strings.TrimSpace(buf.String())
		if lead == "" {
			return body, nil
		}
		if len(bytes.TrimSpace(body)) == 0 {
			return []byte(lead), nil
		}
		return bytes.TrimSpace(append([]byte(lead+"\n\n## 📄 Article\n\n"), body...)), nil
	default:
		lead, err := buildLeadMarkdown(meta, opt)
		if err != nil {
			return nil, err
		}
		if len(bytes.TrimSpace(lead)) == 0 {
			return body, nil
		}
		return append(append(bytes.TrimSpace(lead), []byte("\n\n")...), body...), nil
	}
}
