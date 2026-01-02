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

type RssItemDescription struct {
	Value *string `xml:",cdata"`
}

type RssGUID struct {
	Value       string `xml:",chardata"`
	IsPermaLink bool   `xml:"isPermaLink,attr"`
}

type RssItem struct {
	Title       string             `xml:"title"`
	Link        string             `xml:"link"`
	GUID        RssGUID            `xml:"guid"`
	PubDate     string             `xml:"pubDate"`
	Description RssItemDescription `xml:"description"`
}
