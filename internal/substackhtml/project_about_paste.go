package substackhtml

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ProjectAboutKind selects the same explainer blocks as Hugo partials under
// layouts/partials/*-project-about.html. Strings are read at runtime from i18n/en.toml and i18n/es.toml
// (tables sayingsProjectAbout*, cowsProjectAbout*, reptilocracyProjectAbout*, pawtropolisProjectAbout*) so Substack paste matches the site.
type ProjectAboutKind int

const (
	ProjectAboutNone ProjectAboutKind = iota
	ProjectAboutCubeCows
	ProjectAboutPorEstasCallesSayings
	ProjectAboutReptilocracy
	ProjectAboutPawtropolis
)

const reptilocracyPetitionURL = "https://www.aph.gov.au/e-petitions/petition/EN9806"

// DetectProjectAboutKinds returns explainer blocks to append for Substack HTML, in the same order
// as the site: sayings (Por-Estas-Calles then Reptilocracy), panel/cows (Cube-Cows then Reptilocracy).
func DetectProjectAboutKinds(mdPath string, meta FrontMatterMeta) []ProjectAboutKind {
	p := strings.ToLower(filepath.ToSlash(strings.TrimSpace(mdPath)))
	if !strings.Contains(p, "content/cognitive-memetics/") {
		return nil
	}
	cats := map[string]struct{}{}
	for _, c := range meta.Categories {
		x := strings.TrimSpace(strings.ToLower(c))
		if x != "" {
			cats[x] = struct{}{}
		}
	}
	has := func(name string) bool {
		_, ok := cats[strings.ToLower(strings.TrimSpace(name))]
		return ok
	}

	var out []ProjectAboutKind
	switch {
	case strings.Contains(p, "/sayings/"):
		if has("Por-Estas-Calles") {
			out = append(out, ProjectAboutPorEstasCallesSayings)
		}
		if has("Reptilocracy") {
			out = append(out, ProjectAboutReptilocracy)
		}
	case strings.Contains(p, "/cows/"):
		if has("Cube-Cows") {
			out = append(out, ProjectAboutCubeCows)
		}
		if has("Reptilocracy") {
			out = append(out, ProjectAboutReptilocracy)
		}
	case strings.Contains(p, "/reptilocracy/"):
		if has("Reptilocracy") {
			out = append(out, ProjectAboutReptilocracy)
		}
	case strings.Contains(p, "/pawtropolis/"):
		if has("Pawtropolis-Under-Fire") {
			out = append(out, ProjectAboutPawtropolis)
		}
	default:
		if has("Reptilocracy") {
			out = append(out, ProjectAboutReptilocracy)
		}
		if has("Pawtropolis-Under-Fire") {
			out = append(out, ProjectAboutPawtropolis)
		}
	}
	return out
}

// AppendCognitiveMemeticsProjectAboutHTML appends Substack-friendly explainer sections for cognitive-memetics
// posts that mirror the site "But why" partials. Spanish selects i18n/es copy when true.
func AppendCognitiveMemeticsProjectAboutHTML(html string, mdPath string, meta FrontMatterMeta, spanish bool) (string, error) {
	kinds := DetectProjectAboutKinds(mdPath, meta)
	if len(kinds) == 0 {
		return html, nil
	}
	if err := ensureProjectAboutI18n(); err != nil {
		return "", fmt.Errorf("load cognitive memetics project-about i18n: %w", err)
	}
	pc := projectAboutForLang(spanish)
	var b strings.Builder
	b.WriteString(strings.TrimSpace(html))
	for _, k := range kinds {
		frag, err := renderProjectAboutKind(k, pc)
		if err != nil {
			return "", err
		}
		b.WriteString(frag)
	}
	return b.String(), nil
}

func renderProjectAboutKind(k ProjectAboutKind, pc *projectAboutCopy) (string, error) {
	switch k {
	case ProjectAboutNone:
		return "", nil
	case ProjectAboutCubeCows:
		return renderCubeCowsProjectAbout(pc)
	case ProjectAboutPorEstasCallesSayings:
		return renderSayingsProjectAbout(pc)
	case ProjectAboutReptilocracy:
		return renderReptilocracyProjectAbout(pc)
	case ProjectAboutPawtropolis:
		return renderPawtropolisProjectAbout(pc)
	default:
		return "", nil
	}
}

// markdownProjectAboutBody converts Markdown body copy to HTML; label is used in errors (e.g. "cube-cows").
func markdownProjectAboutBody(bodyMD, label string) (string, error) {
	bodyHTML, err := markdownToHTML([]byte(bodyMD))
	if err != nil {
		return "", fmt.Errorf("%s project about: %w", label, err)
	}
	return bodyHTML, nil
}

// wrapMarkdownProjectAboutBlock converts body Markdown and wraps it in the shared explainer block.
func wrapMarkdownProjectAboutBlock(title, bodyMD, label string) (string, error) {
	bodyHTML, err := markdownProjectAboutBody(bodyMD, label)
	if err != nil {
		return "", err
	}
	return wrapProjectAboutBlock(title, bodyHTML), nil
}

func renderCubeCowsProjectAbout(pc *projectAboutCopy) (string, error) {
	return wrapMarkdownProjectAboutBlock(pc.CowsTitle, pc.CowsBody, "cube-cows")
}

func renderSayingsProjectAbout(pc *projectAboutCopy) (string, error) {
	title := pc.SayingsTitle
	p1 := htmlEscape(pc.SayingsP1)
	p2 := htmlEscape(pc.SayingsP2)
	body := "<p>" + p1 + "</p><p>" + p2 + " <span aria-hidden=\"true\">&#8274;</span></p>"
	return wrapProjectAboutBlock(title, body), nil
}

func renderReptilocracyProjectAbout(pc *projectAboutCopy) (string, error) {
	title := pc.ReptoTitle
	bodyHTML, err := markdownProjectAboutBody(pc.ReptoBody, "reptilocracy")
	if err != nil {
		return "", err
	}
	ctaTitle := htmlEscape(pc.ReptoCtaTitle)
	ctaBtn := htmlEscape(pc.ReptoCtaButton)
	cta := fmt.Sprintf(
		`<p><strong>%s</strong> <a href="%s" rel="noopener noreferrer" target="_blank">%s</a></p>`,
		ctaTitle, reptilocracyPetitionURL, ctaBtn,
	)
	return wrapProjectAboutBlock(title, bodyHTML+cta), nil
}

func renderPawtropolisProjectAbout(pc *projectAboutCopy) (string, error) {
	return wrapMarkdownProjectAboutBlock(pc.PawTitle, pc.PawBody, "pawtropolis")
}

func wrapProjectAboutBlock(title string, innerHTML string) string {
	titleEsc := htmlEscape(title)
	var sb strings.Builder
	sb.WriteString("<hr><blockquote><h3><span aria-hidden=\"true\">&#10067;</span> ")
	sb.WriteString(titleEsc)
	sb.WriteString("</h3>")
	sb.WriteString(strings.TrimSpace(innerHTML))
	sb.WriteString("</blockquote>")
	return sb.String()
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&#34;")
	return s
}
