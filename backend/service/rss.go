package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"blog-server/config"
	"blog-server/entity"
	"blog-server/repository"
)

type RssService interface {
	GenerateRSSFeed(ctx context.Context) (*entity.RSS, error)
	GeneratePagedFeed(ctx context.Context, page, pageSize int) (*entity.RSS, error)
	GenerateCompleteFeed(ctx context.Context) (*entity.RSS, error)
}

type rssService struct {
	cfg      config.AppConfig
	postRepo repository.PostRepo
}

func NewRssService(cfg *config.Config, postRepo repository.PostRepo) RssService {
	return &rssService{cfg: cfg.App, postRepo: postRepo}
}

func (r *rssService) GenerateRSSFeed(
	ctx context.Context,
) (*entity.RSS, error) {
	defaultPageSize := 100
	return r.GeneratePagedFeed(ctx, 1, defaultPageSize)
}

func (r *rssService) GenerateCompleteFeed(ctx context.Context) (*entity.RSS, error) {
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

	feed.Channel.Complete = &entity.FhComplete{}
	feed.Channel.AtomLinks = []entity.AtomLink{
		{
			Href: r.cfg.Domain + "/api/v1/rss/complete",
			Rel:  "self",
			Type: "application/rss+xml",
		},
	}

	return feed, nil
}

func (r *rssService) GeneratePagedFeed(ctx context.Context, page, pageSize int) (*entity.RSS, error) {
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

	hasNext := page < totalPages
	hasPrev := page > 1

	atomLinks := []entity.AtomLink{
		{
			Href: fmt.Sprintf("%s/api/v1/rss?page=%d&pageSize=%d", r.cfg.Domain, page, pageSize),
			Rel:  "self",
			Type: "application/rss+xml",
		},
	}
	if hasNext {
		atomLinks = append(atomLinks, entity.AtomLink{
			Href: fmt.Sprintf("%s/api/v1/rss?page=%d&pageSize=%d", r.cfg.Domain, page+1, pageSize),
			Rel:  "next",
			Type: "application/rss+xml",
		})
	}
	if hasPrev {
		atomLinks = append(atomLinks, entity.AtomLink{
			Href: fmt.Sprintf("%s/api/v1/rss?page=%d&pageSize=%d", r.cfg.Domain, page-1, pageSize),
			Rel:  "previous",
			Type: "application/rss+xml",
		})
	}
	atomLinks = append(atomLinks,
		entity.AtomLink{
			Href: fmt.Sprintf("%s/api/v1/rss?page=1&pageSize=%d", r.cfg.Domain, pageSize),
			Rel:  "first",
			Type: "application/rss+xml",
		},
		entity.AtomLink{
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
		feed.Channel.Complete = &entity.FhComplete{}
	}

	return feed, nil
}

func (r *rssService) newBaseRSS(ctx context.Context) *entity.RSS {
	return &entity.RSS{
		Version: "2.0",
		XMLNs:   "http://www.w3.org/2005/Atom",
		Content: "http://purl.org/rss/1.0/modules/content/",
		DC:      "http://purl.org/dc/elements/1.1/",
		FH:      "http://purl.org/syndication/history/1.0",
		Channel: entity.RssChannel{
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

func (r *rssService) convertPostsToItems(ps []*entity.Post) []entity.RssItem {
	items := make([]entity.RssItem, len(ps))
	for i, post := range ps {
		link := r.cfg.Domain + "/blog/" + strconv.Itoa(int(post.ID))
		items[i] = entity.RssItem{
			Title:       post.Title,
			Link:        link,
			GUID:        entity.RssGUID{Value: link, IsPermaLink: true},
			PubDate:     post.PublishedAt.Format(time.RFC1123Z),
			Description: entity.RssItemDescription{Value: *post.Summary},
		}
	}
	return items
}
