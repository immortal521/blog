package service

import (
	"blog-server/internal/config"
	"blog-server/internal/database"
	"blog-server/internal/response"
	"blog-server/internal/repo"
	"context"
	"encoding/xml"
	"sort"
	"strconv"
	"time"
)

type IRssService interface {
	GenerateRSSFeedXML(ctx context.Context) ([]byte, error)
}

type rssService struct {
	cfg      config.AppConfig
	db       database.DB
	postRepo repo.IPostRepo
}

func (r *rssService) GenerateRSSFeedXML(ctx context.Context) ([]byte, error) {
	posts, err := r.postRepo.GetAllPosts(ctx, r.db.Conn())
	if err != nil {
		return nil, err
	}

	sort.Slice(posts, func(i int, j int) bool {
		return posts[i].UpdatedAt.After(posts[j].UpdatedAt)
	})

	items := make([]response.RssItem, len(posts))
	var pubDate string

	if len(posts) == 0 {
		pubDate = time.Now().Format(time.RFC1123Z)
	} else {
		pubDate = posts[0].UpdatedAt.Format(time.RFC1123Z)
	}

	for i, post := range posts {
		link := r.cfg.Domain + "/blog/" + strconv.Itoa(int(post.ID))
		items[i] = response.RssItem{
			Title:       post.Title,
			Link:        link,
			GUID:        link,
			PubDate:     pubDate,
			Description: post.Summary,
		}
	}

	feed := response.RSS{
		Version: "2.0",
		Channel: response.RssChannel{
			Title:       r.cfg.Name,
			Link:        r.cfg.Domain,
			Description: "Latest posts",
			Items:       items,
			LastBuild:   posts[0].UpdatedAt.String(),
		},
	}

	return xml.MarshalIndent(feed, "", "  ")
}

func NewRssService(cfg *config.Config, db database.DB, postRepo repo.IPostRepo) IRssService {
	return &rssService{cfg: cfg.App, db: db, postRepo: postRepo}
}
