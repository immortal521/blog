package handler

import (
	"blog-server/internal/dto/response"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IRssHandler interface {
	Subscribe(c *fiber.Ctx) error
}

type rssHandler struct {
	svc service.IPostService
}

func (r *rssHandler) Subscribe(c *fiber.Ctx) error {
	posts, err := r.svc.GetPosts(c.UserContext())
	if err != nil {
		return err
	}

	if len(posts) == 0 {
		return errs.NoContent("No posts available")
	}

	sort.Slice(posts, func(i int, j int) bool {
		return posts[i].UpdatedAt.After(posts[j].UpdatedAt)
	})

	items := make([]response.RssItem, len(posts))

	for i, post := range posts {
		items[i] = response.RssItem{
			Title:       post.Title,
			Link:        "https://blog.immortel.top/blog/" + strconv.Itoa(int(post.ID)),
			GUID:        uuid.NewString(),
			PubDate:     post.UpdatedAt.Format(time.RFC1123Z),
			Description: post.Summary,
		}
	}

	feed := response.RSS{
		Version: "2.0",
		Channel: response.RssChannel{
			Title:       "Example Blog",
			Link:        "https://blog.immortel.top",
			Description: "Latest posts",
			Items:       items,
			LastBuild:   posts[0].UpdatedAt.String(),
		},
	}

	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	c.Type("xml")
	c.Append("Content-Disposition", "attachment; filename=rss.xml")
	return c.Send(data)
}

func NewRssHandler(svc service.IPostService) IRssHandler {
	return &rssHandler{svc: svc}
}
