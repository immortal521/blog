package repository

import (
	"context"
	"time"

	"blog-server/datastore"
	"blog-server/ent"
	"blog-server/ent/post"
	"blog-server/ent/user"
	"blog-server/entity"
	"blog-server/mapper"
	"blog-server/pkg/errx"

	"entgo.io/ent/dialect/sql"
)

// PostRepo defines persistence operations for the Post aggregate.
//
// All read operations implicitly apply:
// - published status filter
// - soft-delete filter (DeletedAt IS NULL)
type PostRepo interface {
	GetAllPublished(ctx context.Context) ([]*entity.Post, error)
	GetPublishedByID(ctx context.Context, id uint) (*entity.Post, error)
	UpdateViewCounts(ctx context.Context, updates map[uint]int64) error
	GetPublishedSiteMap(ctx context.Context) ([]*entity.Post, error)
	GetPublishedMeta(ctx context.Context) ([]*entity.Post, error)
	Create(ctx context.Context, post *entity.Post) (*entity.Post, error)
	AddTags(ctx context.Context, postID uint, tagIDs []uint) error
	ReplaceTags(ctx context.Context, postID uint, tagIDs []uint) error
	AddCategories(ctx context.Context, postID uint, categoryIDs []uint) error
	ReplaceCategories(ctx context.Context, postID uint, categoryIDs []uint) error
	IsOwner(ctx context.Context, userID uint, postID uint) (bool, error)
}

type postRepo struct {
	ds *datastore.DataStore
}

// NewPostRepo creates a PostRepo instance using the given datastore.
func NewPostRepo(ds *datastore.DataStore) PostRepo {
	return &postRepo{ds: ds}
}

// IsOwner checks whether a user is the author of a given post.
func (r *postRepo) IsOwner(ctx context.Context, userID uint, postID uint) (bool, error) {
	count, err := r.ds.Client(ctx).Post.
		Query().
		Where(
			post.IDEQ(postID),
			post.UserIDEQ(userID),
		).
		Count(ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Create inserts a new post record.
//
// Behavior:
// - Automatically sets created_at and updated_at timestamps
// - If status is Published, published_at is set automatically
// - Optional fields (summary, cover) are only set when provided
func (r *postRepo) Create(ctx context.Context, p *entity.Post) (*entity.Post, error) {
	now := time.Now()

	builder := r.ds.Client(ctx).Post.
		Create().
		SetTitle(p.Title).
		SetContent(p.Content).
		SetReadTimeMinutes(p.ReadTimeMinutes).
		SetStatus(p.Status).
		SetUserID(p.UserID).
		SetCreatedAt(now).
		SetUpdatedAt(now)

	if p.Summary != nil {
		builder.SetSummary(*p.Summary)
	}

	if p.Cover != nil {
		builder.SetCover(*p.Cover)
	}

	if p.Status == entity.PostStatusPublish {
		builder.SetPublishedAt(now)
	}

	ep, err := builder.Save(ctx)
	if err != nil {
		return nil, errx.New(errx.CodeInternalError, err)
	}

	return mapper.ToPost(ep), nil
}

// AddCategories attaches categories to a post without removing existing relations.
func (r *postRepo) AddCategories(ctx context.Context, postID uint, categoryIDs []uint) error {
	if len(categoryIDs) == 0 {
		return nil
	}

	err := r.ds.Client(ctx).Post.
		UpdateOneID(postID).
		AddCategoryIDs(categoryIDs...).
		Exec(ctx)
	if err != nil {
		return errx.New(errx.CodeInternalError, err)
	}
	return nil
}

// AddTags attaches tags to a post without removing existing relations.
func (r *postRepo) AddTags(ctx context.Context, postID uint, tagIDs []uint) error {
	if len(tagIDs) == 0 {
		return nil
	}

	err := r.ds.Client(ctx).Post.
		UpdateOneID(postID).
		AddTagIDs(tagIDs...).
		Exec(ctx)
	if err != nil {
		return errx.New(errx.CodeInternalError, err)
	}
	return nil
}

// ReplaceCategories replaces all categories of a post (clear then add).
func (r *postRepo) ReplaceCategories(ctx context.Context, postID uint, categoryIDs []uint) error {
	builder := r.ds.Client(ctx).Post.
		UpdateOneID(postID).
		ClearCategories()

	if len(categoryIDs) > 0 {
		builder.AddCategoryIDs(categoryIDs...)
	}

	err := builder.Exec(ctx)
	if err != nil {
		return errx.New(errx.CodeInternalError, err)
	}
	return nil
}

// ReplaceTags replaces all tags of a post (clear then add).
func (r *postRepo) ReplaceTags(ctx context.Context, postID uint, tagIDs []uint) error {
	builder := r.ds.Client(ctx).Post.
		UpdateOneID(postID).
		ClearTags()

	if len(tagIDs) > 0 {
		builder.AddTagIDs(tagIDs...)
	}

	err := builder.Exec(ctx)
	if err != nil {
		return errx.New(errx.CodeInternalError, err)
	}
	return nil
}

// GetAllPublished returns all published posts for list views.
//
// Optimizations:
// - Selects only fields required for listing (no content field)
// - Preloads author, categories, and tags
// - Ordered by published time descending
func (r *postRepo) GetAllPublished(ctx context.Context) ([]*entity.Post, error) {
	ps, err := r.ds.Client(ctx).Post.
		Query().
		Select(
			post.FieldID,
			post.FieldTitle,
			post.FieldSummary,
			post.FieldCover,
			post.FieldReadTimeMinutes,
			post.FieldViewCount,
			post.FieldPublishedAt,
			post.FieldCreatedAt,
			post.FieldUpdatedAt,
		).
		WithAuthor(func(q *ent.UserQuery) {
			q.Select(user.FieldUsername)
		}).
		WithCategories().
		WithTags().
		Where(
			post.StatusEQ(entity.PostStatusPublish),
			post.DeletedAtIsNil(),
		).
		Order(
			post.ByPublishedAt(sql.OrderDesc()),
		).
		All(ctx)
	if err != nil {
		return nil, errx.New(errx.CodeInternalError, err)
	}

	return mapper.ToPosts(ps), nil
}

// GetPublishedByID returns a single published post by ID.
//
// Notes:
// - Includes full content field
// - Preloads author information
// - Returns NotFound error if no matching record exists
func (r *postRepo) GetPublishedByID(ctx context.Context, id uint) (*entity.Post, error) {
	p, err := r.ds.Client(ctx).Post.
		Query().
		Select(
			post.FieldID,
			post.FieldTitle,
			post.FieldSummary,
			post.FieldCover,
			post.FieldContent,
			post.FieldReadTimeMinutes,
			post.FieldViewCount,
			post.FieldPublishedAt,
			post.FieldCreatedAt,
			post.FieldUpdatedAt,
		).
		WithAuthor(func(q *ent.UserQuery) {
			q.Select(user.FieldUsername)
		}).
		Where(
			post.IDEQ(id),
			post.StatusEQ(entity.PostStatusPublish),
			post.DeletedAtIsNil(),
		).
		First(ctx)
	if err != nil {
		return nil, errx.New(errx.CodeNotFound, err)
	}

	return mapper.ToPost(p), nil
}

// GetPublishedSiteMap returns minimal post data for sitemap generation.
//
// Only ID and UpdatedAt are selected for lightweight crawling support.
func (r *postRepo) GetPublishedSiteMap(ctx context.Context) ([]*entity.Post, error) {
	ps, err := r.ds.Client(ctx).Post.
		Query().
		Select(
			post.FieldID,
			post.FieldUpdatedAt,
		).
		Where(
			post.StatusEQ(entity.PostStatusPublish),
			post.DeletedAtIsNil(),
		).
		All(ctx)
	if err != nil {
		return nil, errx.New(errx.CodeInternalError, err)
	}

	return mapper.ToPosts(ps), nil
}

// GetPublishedMeta returns lightweight post metadata for SEO and previews.
//
// Ordered by published time descending.
func (r *postRepo) GetPublishedMeta(ctx context.Context) ([]*entity.Post, error) {
	ps, err := r.ds.Client(ctx).Post.
		Query().
		Select(
			post.FieldID,
			post.FieldTitle,
			post.FieldSummary,
			post.FieldUpdatedAt,
			post.FieldPublishedAt,
		).
		Where(
			post.StatusEQ(entity.PostStatusPublish),
			post.DeletedAtIsNil(),
		).
		Order(
			post.ByPublishedAt(sql.OrderDesc()),
		).
		All(ctx)
	if err != nil {
		return nil, errx.New(errx.CodeInternalError, err)
	}

	return mapper.ToPosts(ps), nil
}

// UpdateViewCounts performs a bulk increment of view counts.
//
// Implementation details:
// - Uses a single SQL UPDATE with CASE expression
// - Each post ID is updated atomically
// - Invalid (non-positive) deltas are ignored
// - Missing IDs are silently skipped by the database
//
// Note: query complexity grows linearly with batch size.
func (r *postRepo) UpdateViewCounts(ctx context.Context, updates map[uint]int64) error {
	for id, delta := range updates {
		if delta <= 0 {
			delete(updates, id)
		}
	}
	if len(updates) == 0 {
		return nil
	}

	ids := make([]uint, 0, len(updates))
	for id := range updates {
		ids = append(ids, id)
	}

	err := r.ds.Client(ctx).Post.
		Update().
		Where(post.IDIn(ids...)).
		Modify(func(u *sql.UpdateBuilder) {
			caseExpr := func(b *sql.Builder) {
				b.WriteString("CASE ")
				b.Ident(post.FieldID)
				b.WriteString(" ")

				for _, id := range ids {
					b.WriteString("WHEN ")
					b.Arg(id)
					b.WriteString(" THEN ")
					b.Ident(post.FieldViewCount)
					b.WriteString(" + ")
					b.Arg(updates[id])
					b.WriteString(" ")
				}

				b.WriteString("ELSE ")
				b.Ident(post.FieldViewCount)
				b.WriteString(" END")
			}
			u.Set(post.FieldViewCount, sql.ExprFunc(caseExpr))
		}).
		Exec(ctx)
	if err != nil {
		return errx.New(errx.CodeInternalError, err)
	}
	return nil
}
