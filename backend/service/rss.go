package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"blog-server/config"
	"blog-server/entity"
	"blog-server/repository"
)

type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RssChannel `xml:"channel"`
	XMLNs   string     `xml:"xmlns:atom,attr"`
	Content string     `xml:"xmlns:content,attr"`
	DC      string     `xml:"xmlns:dc,attr"`
	FH      string     `xml:"xmlns:fh,attr,omitempty"`
}

type RssChannel struct {
	Title       string      `xml:"title"`
	Link        string      `xml:"link"`
	Description string      `xml:"description"`
	LastBuild   string      `xml:"lastBuildDate"`
	Items       []RssItem   `xml:"item"`
	AtomLinks   []AtomLink  `xml:"atom:link"`
	Complete    *FhComplete `xml:"fh:complete,omitempty"`
}

type AtomLink struct {
	Href  string `xml:"href,attr"`
	Rel   string `xml:"rel,attr"`
	Type  string `xml:"type,attr"`
	Title string `xml:"title,attr,omitempty"`
}

type FhComplete struct{}

type RssItemDescription struct {
	Value string `xml:",cdata"`
}

type RssItemContent struct {
	Value string `xml:",cdata"`
}

type RssCategory struct {
	Domain string `xml:"domain,attr,omitempty"`
	Value  string `xml:",chardata"`
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
	Content     *RssItemContent    `xml:"content:encoded,omitempty"`
	Author      string             `xml:"dc:creator,omitempty"`
	Categories  []RssCategory      `xml:"category,omitempty"`
}

type PaginationInfo struct {
	CurrentPage int
	TotalPages  int
	TotalItems  int
	HasNext     bool
	HasPrev     bool
}

type RssService interface {
	GenerateRSSFeedXML(ctx context.Context) ([]byte, error)
	GeneratePagedFeedXML(ctx context.Context, page, pageSize int) ([]byte, error)
	GenerateCompleteFeedXML(ctx context.Context) ([]byte, error)
}

type rssService struct {
	cfg      config.AppConfig
	postRepo repository.PostRepo
}

func NewRssService(cfg *config.Config, postRepo repository.PostRepo) RssService {
	return &rssService{cfg: cfg.App, postRepo: postRepo}
}

func (r *rssService) GenerateRSSFeedXML(ctx context.Context) ([]byte, error) {
	defaultPageSize := 100
	return r.GeneratePagedFeedXML(ctx, 1, defaultPageSize)
}

func (r *rssService) GenerateCompleteFeedXML(ctx context.Context) ([]byte, error) {
	total, err := r.postRepo.CountPublished(ctx)
	if err != nil || total == 0 {
		total = 500
	}

	ps, err := r.postRepo.ListPublishedForMeta(ctx, 1, total)
	if err != nil {
		return nil, err
	}

	feed := r.newBaseRSS(ctx)
	feed.Channel.Title = r.cfg.Name + " (Complete Archive)"
	feed.Channel.Description = "Full history archive of posts"
	feed.Channel.Items = r.convertPostsToItems(ps)
	feed.Channel.Complete = &FhComplete{}
	feed.Channel.AtomLinks = []AtomLink{
		{
			Href: r.cfg.Domain + "/api/v1/rss/complete",
			Rel:  "self",
			Type: "application/rss+xml",
		},
	}

	return xml.MarshalIndent(feed, "", "  ")
}

func (r *rssService) GeneratePagedFeedXML(ctx context.Context, page, pageSize int) ([]byte, error) {
	total, _ := r.postRepo.CountPublished(ctx)

	if page < 1 {
		page = 1
	}

	ps, err := r.postRepo.ListPublishedForMeta(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	pagination := &PaginationInfo{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}

	atomLinks := []AtomLink{
		{
			Href: fmt.Sprintf("%s/api/v1/rss?page=%d&pageSize=%d", r.cfg.Domain, page, pageSize),
			Rel:  "self",
			Type: "application/rss+xml",
		},
	}
	if pagination.HasNext {
		atomLinks = append(atomLinks, AtomLink{
			Href: fmt.Sprintf("%s/api/v1/rss?page=%d&pageSize=%d", r.cfg.Domain, page+1, pageSize),
			Rel:  "next",
			Type: "application/rss+xml",
		})
	}
	if pagination.HasPrev {
		atomLinks = append(atomLinks, AtomLink{
			Href: fmt.Sprintf("%s/api/v1/rss?page=%d&pageSize=%d", r.cfg.Domain, page-1, pageSize),
			Rel:  "previous",
			Type: "application/rss+xml",
		})
	}
	atomLinks = append(atomLinks,
		AtomLink{
			Href: fmt.Sprintf("%s/api/v1/rss?page=1&pageSize=%d", r.cfg.Domain, pageSize),
			Rel:  "first",
			Type: "application/rss+xml",
		},
		AtomLink{
			Href: fmt.Sprintf("%s/api/v1/rss?page=%d&pageSize=%d", r.cfg.Domain, totalPages, pageSize),
			Rel:  "last",
			Type: "application/rss+xml",
		},
	)

	feed := r.newBaseRSS(ctx)
	feed.Channel.AtomLinks = atomLinks
	feed.Channel.Description = fmt.Sprintf("Latest posts - Page %d of %d", page, totalPages)
	feed.Channel.Items = r.convertPostsToItems(ps)

	if page > 1 {
		feed.Channel.Complete = &FhComplete{}
	}

	return xml.MarshalIndent(feed, "", "  ")
}

func (r *rssService) newBaseRSS(ctx context.Context) RSS {
	return RSS{
		Version: "2.0",
		XMLNs:   "http://www.w3.org/2005/Atom",
		Content: "http://purl.org/rss/1.0/modules/content/",
		DC:      "http://purl.org/dc/elements/1.1/",
		FH:      "http://purl.org/syndication/history/1.0",
		Channel: RssChannel{
			Title:     r.cfg.Name,
			Link:      r.cfg.Domain,
			LastBuild: r.getBuildTime(ctx),
		},
	}
}

func (r *rssService) getBuildTime(ctx context.Context) string {
	latestPublishedAtPtr, err := r.postRepo.GetLatestPublishedAt(ctx)
	if err != nil || latestPublishedAtPtr == nil {
		return time.Now().Format(time.RFC1123Z)
	}
	return latestPublishedAtPtr.Format(time.RFC1123Z)
}

func (r *rssService) convertPostsToItems(ps []*entity.Post) []RssItem {
	items := make([]RssItem, len(ps))
	for i, post := range ps {
		link := r.cfg.Domain + "/blog/" + strconv.Itoa(int(post.ID))
		items[i] = RssItem{
			Title:       post.Title,
			Link:        link,
			GUID:        RssGUID{Value: link, IsPermaLink: true},
			PubDate:     post.PublishedAt.Format(time.RFC1123Z),
			Description: RssItemDescription{Value: *post.Summary},
		}
	}
	return items
}
