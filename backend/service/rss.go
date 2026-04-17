package service

import (
	"context"
	"encoding/xml"
	"sort"
	"strconv"
	"time"

	"blog-server/config"
	"blog-server/repository"
)

type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RssChannel `xml:"channel"`
	XMLNs   string     `xml:"xmlns:atom,attr"`
	Content string     `xml:"xmlns:content,attr"`
	DC      string     `xml:"xmlns:dc,attr"`
}

type RssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	LastBuild   string    `xml:"lastBuildDate"`
	Items       []RssItem `xml:"item"`
	AtomLink    AtomLink  `xml:"atom:link"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
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

type RssService interface {
	GenerateRSSFeedXML(ctx context.Context) ([]byte, error)
}

type rssService struct {
	cfg      config.AppConfig
	postRepo repository.PostRepo
}

func (r *rssService) GenerateRSSFeedXML(ctx context.Context) ([]byte, error) {
	ps, err := r.postRepo.GetPublishedMeta(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(ps, func(i int, j int) bool {
		return ps[i].UpdatedAt.After(ps[j].UpdatedAt)
	})

	items := make([]RssItem, len(ps))
	var pubDate string

	if len(ps) == 0 {
		pubDate = time.Now().Format(time.RFC1123Z)
	} else {
		pubDate = ps[0].PublishedAt.Format(time.RFC1123Z)
	}

	for i, post := range ps {
		link := r.cfg.Domain + "/blog/" + strconv.Itoa(int(post.ID))
		items[i] = RssItem{
			Title:       post.Title,
			Link:        link,
			GUID:        RssGUID{Value: link, IsPermaLink: true},
			PubDate:     post.PublishedAt.Format(time.RFC1123Z),
			Description: RssItemDescription{Value: post.Summary},
		}
	}

	feed := RSS{
		Version: "2.0",
		XMLNs:   "http://www.w3.org/2005/Atom",
		Content: "http://purl.org/rss/1.0/modules/content/",
		DC:      "http://purl.org/dc/elements/1.1/",
		Channel: RssChannel{
			Title:       r.cfg.Name,
			Link:        r.cfg.Domain,
			Description: "Latest posts",
			Items:       items,
			LastBuild:   pubDate,
			AtomLink: AtomLink{
				Href: r.cfg.Domain + "/api/v1/rss",
				Rel:  "self",
				Type: "application/rss+xml",
			},
		},
	}

	return xml.MarshalIndent(feed, "", "  ")
}

func NewRssService(cfg *config.Config, postRepo repository.PostRepo) RssService {
	return &rssService{cfg: cfg.App, postRepo: postRepo}
}
