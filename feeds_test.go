package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
)

func TestFetchDataFromFeed(t *testing.T) {

	mockRSS := `
	<rss version="2.0">
	<channel>
	<title>Boot.dev Blog</title>
	<link>https://blog.boot.dev/</link>
	<description>Recent content on Boot.dev Blog</description>
	<generator>Hugo</generator>
	<language>en-us</language>
	<lastBuildDate>Wed, 05 Jun 2024 00:00:00 +0000</lastBuildDate>
	<atom:link href="https://blog.boot.dev/index.xml" rel="self" type="application/rss+xml"/>
	<item>
	<title>The Boot.dev Beat. June 2024</title>
	<link>https://blog.boot.dev/news/bootdev-beat-2024-06/</link>
	<pubDate>Wed, 05 Jun 2024 00:00:00 +0000</pubDate>
	<guid>https://blog.boot.dev/news/bootdev-beat-2024-06/</guid>
	<description>
	<![CDATA[ThePrimeagen&rsquo;s new Git course is live. A new boss battle is on the horizon, and we&rsquo;ve made massive speed improvements to the site.]]>
	</description>
	</item>
	</channel>
	</rss>
	`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, mockRSS)
	}))
	defer server.Close()

	rss, err := fetchDataFromFeed(server.URL)
	if err != nil {
		t.Fatalf("Failed to fetch RSS feed: %v", err)
	}

	fmt.Printf("%v", rss.Channel.Items[0].Description)

	expectedTitle := "Boot.dev Blog"
	expectedItemTitle := "The Boot.dev Beat. June 2024"

	if rss.Channel.Title != expectedTitle {
		t.Errorf("Expected title %q, but got %q", expectedTitle, rss.Channel.Title)
	}

	if len(rss.Channel.Items) == 0 {
		t.Fatalf("Expected at least one item in the RSS feed, but got none")
	}

	if rss.Channel.Items[0].Title != expectedItemTitle {
		t.Errorf("Expected item title %q, but got %q", expectedItemTitle, rss.Channel.Items[0].Title)
	}
}
