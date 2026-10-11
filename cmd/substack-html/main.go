// Command substack-html converts Hugo-style Markdown from a file into HTML
// suitable for pasting into Substack drafts. Browser automation is intentionally
// out of scope; this command only prints or writes HTML.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dualsubstrate/blog/internal/cliout"
	"github.com/dualsubstrate/blog/internal/substackbrowser"
	"github.com/dualsubstrate/blog/internal/substackhtml"
)

func main() {
	log.SetFlags(0)

	inPath := flag.String("in", "", "path to Markdown file (required)")
	outPath := flag.String("out", "", "output path (empty: write to stdout)")
	tables := flag.String("tables", "html", `table output: "html" (minimal table) or "list" (ordered list rows)`)
	paragraphs := flag.String("paragraphs", "", `paragraph output: "p" or "br" (empty: default)`)
	paraBreak := flag.Int("paragraph-break-br-count", 0, "when -paragraphs=br, how many <br> tags between paragraphs (0: default)")
	quotes := flag.String("quotes", "", `blockquote output: "blockquote", "monospace", or "pullquote_monospace" (empty: default)`)
	configPath := flag.String("config", substackbrowser.DefaultLocalConfigPath(), "optional local config path (default: substack.json at repo root)")
	configGlobal := flag.String("config-global", "", "optional shared base config JSON merged under -config (overlay wins); same grouped vs flat shape as -config")
	noDemote := flag.Bool("no-demote-headings", false, "keep heading levels (default demotes h1..h5 down one for editor paste)")
	includeLead := flag.Bool("include-frontmatter-lead", false, "prepend front matter intro (description/sowhat/youtube_id) before body")
	markdownLeadImageResolveOrigin := flag.String("markdown-lead-image-resolve-origin", "", `optional origin for featured image URLs (e.g. http://localhost:1313); overrides substack.json; same as substack-draft`)
	flag.Parse()

	if *inPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	localCfg, localCfgFound, localCfgErr := substackbrowser.LoadLocalConfigWithGlobal(*configGlobal, *configPath)
	if localCfgErr != nil {
		log.Fatalf("read config: %v", localCfgErr)
	}
	if localCfgFound {
		if strings.TrimSpace(*tables) == "html" && strings.TrimSpace(localCfg.TableMode) != "" {
			*tables = localCfg.TableMode
		}
		if strings.TrimSpace(*paragraphs) == "" && strings.TrimSpace(localCfg.ParagraphMode) != "" {
			*paragraphs = localCfg.ParagraphMode
		}
		if *paraBreak == 0 && localCfg.ParagraphBreakBRCount > 0 {
			*paraBreak = localCfg.ParagraphBreakBRCount
		}
		if strings.TrimSpace(*quotes) == "" && strings.TrimSpace(localCfg.QuoteMode) != "" {
			*quotes = localCfg.QuoteMode
		}
		if !*noDemote && localCfg.DemoteHeadings != nil && !*localCfg.DemoteHeadings {
			*noDemote = true
		}
		if !*includeLead && localCfg.IncludeFrontMatterLead {
			*includeLead = true
		}
	}

	raw, err := os.ReadFile(*inPath)
	if err != nil {
		log.Fatalf("read input: %v", err)
	}
	indexRaw := raw
	bodyRes, fbErr := substackhtml.ResolveSubstackBody(raw, *inPath)
	if fbErr != nil {
		log.Fatalf("substack body: %v", fbErr)
	}
	usedSubstackSidecar := bodyRes.Source == substackhtml.SubstackBodyFromSidecarMD
	raw = bodyRes.IndexRaw

	opt := substackhtml.DefaultOptions()
	switch *tables {
	case "html":
		opt.TableMode = substackhtml.TableHTML
	case "list":
		opt.TableMode = substackhtml.TableList
	default:
		log.Fatalf("unknown -tables value %q (use html or list)", *tables)
	}
	switch strings.TrimSpace(*paragraphs) {
	case "":
		// keep default
	case "p":
		opt.ParagraphMode = substackhtml.ParagraphP
	case "br":
		opt.ParagraphMode = substackhtml.ParagraphBR
	default:
		log.Fatalf("unknown -paragraphs value %q (use p or br)", *paragraphs)
	}
	if *paraBreak > 0 {
		opt.ParagraphBreakBRCount = *paraBreak
	}
	switch strings.TrimSpace(*quotes) {
	case "":
		// keep default
	case "blockquote":
		opt.QuoteMode = substackhtml.QuoteBlockquote
	case "monospace":
		opt.QuoteMode = substackhtml.QuoteMonospace
	case "pullquote_monospace":
		opt.QuoteMode = substackhtml.QuotePullquoteMonospace
	default:
		log.Fatalf("unknown -quotes value %q (use blockquote, monospace, or pullquote_monospace)", *quotes)
	}
	if *noDemote {
		opt.DemoteHeadings = false
	}
	if *includeLead && !usedSubstackSidecar {
		opt.IncludeFrontMatterLead = true
	}
	if usedSubstackSidecar {
		opt.IncludeFeaturedImageLead = true
	}

	opt.SourcePath = strings.TrimSpace(*inPath)
	if opt.SourcePath != "" {
		if perm, err := substackhtml.ResolvePagePermalinkForMarkdown(opt.SourcePath); err == nil && strings.TrimSpace(perm) != "" {
			opt.PagePermalink = strings.TrimSpace(perm)
		}
	}
	if strings.TrimSpace(*markdownLeadImageResolveOrigin) != "" {
		opt.ImageResolveOrigin = strings.TrimSpace(*markdownLeadImageResolveOrigin)
	} else {
		opt.ImageResolveOrigin = substackbrowser.EffectiveMarkdownLeadImageResolveOrigin(localCfg)
	}
	html, err := substackhtml.Convert(raw, opt)
	if err != nil {
		log.Fatalf("convert: %v", err)
	}

	meta := substackhtml.ExtractFrontMatterMeta(indexRaw)
	if localCfgFound && substackbrowser.EffectiveIncludeCognitiveMemeticsProjectAbout(localCfg) {
		spanish := substackhtml.IsSpanishSiteLocale(meta, opt.SourcePath)
		html, err = substackhtml.AppendCognitiveMemeticsProjectAboutHTML(html, opt.SourcePath, meta, spanish)
		if err != nil {
			log.Fatalf("cognitive-memetics project about: %v", err)
		}
	}

	if *outPath == "" {
		fmt.Print(html)
		return
	}
	if err := os.WriteFile(*outPath, []byte(html+"\n"), 0o644); err != nil {
		log.Fatalf("write output: %v", err)
	}
	cliout.PrintSubstackHTMLWritten(os.Stdout, *inPath, *outPath, len(html))
}
