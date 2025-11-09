package response

import "encoding/xml"

type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RssChannel `xml:"channel"`
}

type RssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	LastBuild   string    `xml:"lastBuildDate"`
	Items       []RssItem `xml:"item"`
}

type RssItem struct {
	Title       string  `xml:"title"`
	Content     string  `xml:"content"`
	Link        string  `xml:"link"`
	GUID        string  `xml:"guid"`
	PubDate     string  `xml:"pubDate"`
	Description *string `xml:"description"`
}
