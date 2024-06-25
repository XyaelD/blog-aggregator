package main

import "encoding/xml"

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel represents the channel in the RSS feed
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	Items       []Item `xml:"item"`
}

// Item represents an item in the RSS feed with additional optional fields
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate,omitempty"`
	Author      string `xml:"author,omitempty"`
	Content     string `xml:"encoded,omitempty" xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
}
