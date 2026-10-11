package linkedinapi

import "strings"

// ArticleContent is LinkedIn Posts API content.article (link card; API does not scrape URLs).
type ArticleContent struct {
	Source       string // destination URL (e.g. YouTube watch link)
	ThumbnailURN string // urn:li:image:... from Images API
	Title        string
	Description  string
}

// PostOptions selects text-only, single image media, or article link card.
type PostOptions struct {
	ImageURN string
	AltText  string
	Article  *ArticleContent
}

func (o PostOptions) hasArticle() bool {
	return o.Article != nil && strings.TrimSpace(o.Article.Source) != ""
}

func (o PostOptions) hasMedia() bool {
	return strings.TrimSpace(o.ImageURN) != "" && !o.hasArticle()
}
