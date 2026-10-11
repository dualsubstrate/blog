package facebookautopost

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"
)

// NamespaceRSSSync is the xmlns URI for sync:* elements emitted in index.xml (plain social bundle text).
const NamespaceRSSSync = "https://blog.dualsubstrate.com/xmlns/rss-sync/1.0"

// Item is a normalized RSS entry for posting.
type Item struct {
	Title                string
	Link                 string
	Description          string
	GUID                 string
	PublishedAt          time.Time
	PageSocialPlain      string
	EnclosureImageURL    string
	MediaContentImageURL string
	MediaThumbnailURL    string
}

// SyncPostCaption returns sync:pageSocial text when present, otherwise title, description, and link (plain).
func SyncPostCaption(it Item) string {
	if s := strings.TrimSpace(it.PageSocialPlain); s != "" {
		return strings.TrimRight(s, "\n")
	}
	return strings.TrimRight(strings.TrimSpace(fmt.Sprintf("%s\n\n%s\n\n%s",
		strings.TrimSpace(it.Title),
		strings.TrimSpace(it.Description),
		strings.TrimSpace(it.Link),
	)), "\n")
}

// PreferredImageURL picks the first non-empty URL in priority order:
// enclosure, media:content, media:thumbnail.
func PreferredImageURL(it Item) string {
	for _, u := range []string{
		strings.TrimSpace(it.EnclosureImageURL),
		strings.TrimSpace(it.MediaContentImageURL),
		strings.TrimSpace(it.MediaThumbnailURL),
	} {
		if u != "" {
			return u
		}
	}
	return ""
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Enclosure   *struct {
		URL    string `xml:"url,attr"`
		Type   string `xml:"type,attr"`
		Length string `xml:"length,attr"`
	} `xml:"enclosure"`
	MediaThumbnail *struct {
		URL string `xml:"url,attr"`
	} `xml:"http://search.yahoo.com/mrss/ thumbnail"`
	MediaContent *struct {
		URL    string `xml:"url,attr"`
		Type   string `xml:"type,attr"`
		Medium string `xml:"medium,attr"`
	} `xml:"http://search.yahoo.com/mrss/ content"`
	ContentEncoded string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	PageSocial     string `xml:"https://blog.dualsubstrate.com/xmlns/rss-sync/1.0 pageSocial"`
}

// ParseFeed reads RSS 2.0 XML (with media RSS and content modules) and returns items in document order.
func ParseFeed(r io.Reader) ([]Item, error) {
	var raw rssFeed
	dec := xml.NewDecoder(r)
	dec.Strict = false
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf("parse rss: %w", err)
	}
	out := make([]Item, 0, len(raw.Channel.Items))
	for _, ri := range raw.Channel.Items {
		it := Item{
			Title:           strings.TrimSpace(ri.Title),
			Link:            strings.TrimSpace(ri.Link),
			Description:     strings.TrimSpace(ri.Description),
			GUID:            strings.TrimSpace(ri.GUID),
			PageSocialPlain: strings.TrimSpace(ri.PageSocial),
		}
		if s := strings.TrimSpace(ri.PubDate); s != "" {
			// RSS 2.0 pubDate is typically RFC1123Z.
			if t, err := time.Parse(time.RFC1123Z, s); err == nil {
				it.PublishedAt = t
			} else if t, err := time.Parse(time.RFC1123, s); err == nil {
				it.PublishedAt = t
			}
		}
		if it.GUID == "" {
			it.GUID = it.Link
		}
		if ri.Enclosure != nil {
			it.EnclosureImageURL = strings.TrimSpace(ri.Enclosure.URL)
		}
		if ri.MediaContent != nil {
			it.MediaContentImageURL = strings.TrimSpace(ri.MediaContent.URL)
		}
		if ri.MediaThumbnail != nil {
			it.MediaThumbnailURL = strings.TrimSpace(ri.MediaThumbnail.URL)
		}
		out = append(out, it)
	}
	return out, nil
}
