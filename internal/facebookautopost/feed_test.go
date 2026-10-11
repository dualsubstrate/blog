package facebookautopost

import (
	"strings"
	"testing"
)

const sampleRSS = `<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:atom="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:content="http://purl.org/rss/1.0/modules/content/" version="2.0">
<channel>
<item>
<title>T1</title>
<link>https://example.com/a/</link>
<guid>https://example.com/a/</guid>
<pubDate>Thu, 07 May 2026 12:00:00 +0000</pubDate>
<description>Short</description>
<enclosure url="https://example.com/a/en.jpg" type="image/jpeg" length="0"/>
<media:thumbnail url="https://example.com/a/thumb.jpg"/>
<media:content url="https://example.com/a/content.jpg" type="image/jpeg" medium="image"/>
</item>
<item>
<title>No enclosure</title>
<link>https://example.com/b/</link>
<guid>b-guid</guid>
<description>D2</description>
<media:thumbnail url="https://example.com/b/thumb-only.jpg"/>
</item>
<item>
<title>Thumb only alt</title>
<link>https://example.com/c/</link>
<guid>c-guid</guid>
<description>D3</description>
<media:content url="https://example.com/c/mc.jpg" type="image/png" medium="image"/>
</item>
</channel>
</rss>`

func TestParseProductionFeedShape(t *testing.T) {
	one := `<rss xmlns:media="http://search.yahoo.com/mrss/" xmlns:content="http://purl.org/rss/1.0/modules/content/" version="2.0"><channel><item><title>T</title><link>https://blog.dualsubstrate.com/x/</link><guid>https://blog.dualsubstrate.com/x/</guid><description>D</description><enclosure url="https://blog.dualsubstrate.com/x/a.png" type="image/png" length="0"/><media:thumbnail url="https://blog.dualsubstrate.com/x/t.png"/><media:content url="https://blog.dualsubstrate.com/x/c.png" type="image/png" medium="image"/></item></channel></rss>`
	items, err := ParseFeed(strings.NewReader(one))
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("got %d items", len(items))
	}
	if PreferredImageURL(items[0]) != "https://blog.dualsubstrate.com/x/a.png" {
		t.Fatalf("image: %q", PreferredImageURL(items[0]))
	}
}

func TestParseSyncPageSocial(t *testing.T) {
	const syncRSS = `<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:sync="https://blog.dualsubstrate.com/xmlns/rss-sync/1.0" version="2.0">
<channel>
<item>
<title>T</title>
<link>https://example.com/p/</link>
<guid>https://example.com/p/</guid>
<description>Short card</description>
<sync:pageSocial><![CDATA[W16: Hub

✔️ TLDR: body line.

🔗 https://example.com/p/]]></sync:pageSocial>
</item>
</channel>
</rss>`
	items, err := ParseFeed(strings.NewReader(syncRSS))
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("items: %d", len(items))
	}
	if !strings.Contains(items[0].PageSocialPlain, "TLDR") {
		t.Fatalf("PageSocialPlain: %q", items[0].PageSocialPlain)
	}
	cap := SyncPostCaption(items[0])
	if !strings.Contains(cap, "TLDR") || strings.Contains(cap, "Short card") {
		t.Fatalf("SyncPostCaption should prefer page social: %q", cap)
	}
}

func TestParseFeedPreferredImage(t *testing.T) {
	items, err := ParseFeed(strings.NewReader(sampleRSS))
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 {
		t.Fatalf("items: got %d want 3", len(items))
	}
	if items[0].PublishedAt.IsZero() {
		t.Fatalf("expected PublishedAt for item0")
	}
	if PreferredImageURL(items[0]) != "https://example.com/a/en.jpg" {
		t.Errorf("item0 image: %q", PreferredImageURL(items[0]))
	}
	if PreferredImageURL(items[1]) != "https://example.com/b/thumb-only.jpg" {
		t.Errorf("item1 image: %q", PreferredImageURL(items[1]))
	}
	if PreferredImageURL(items[2]) != "https://example.com/c/mc.jpg" {
		t.Errorf("item2 image: %q", PreferredImageURL(items[2]))
	}
}
