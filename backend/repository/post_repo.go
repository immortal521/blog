package repository

import (
	"context"

	"blog-server/datastore"
	"blog-server/ent"
	"blog-server/ent/post"
	"blog-server/ent/user"
	"blog-server/entity"
	"blog-server/mapper"

	"entgo.io/ent/dialect/sql"
)

// PostRepo defines read and write operations for posts.
//
// Implementations are expected to apply visibility constraints consistently across all read methods.
type PostRepo interface {
	// GetAllPublished returns all posts that are published and not soft-deleted,
	// ordered by publication time in descending order.
	//
	// The returned slice is non-nil on success, even when no posts exist.
	// The result excludes unpublished or soft-deleted entries.
	GetAllPublished(ctx context.Context) ([]*entity.Post, error)

	// GetPublishedByID returns a single published, non-deleted post by its identifier.
	//
	// If no such post exists, an error is returned. Callers should treat not-found
	// errors distinctly from other failures if needed.
	GetPublishedByID(ctx context.Context, id uint) (*entity.Post, error)

	// UpdateViewCounts applies batched increments to post view counts.
	//
	// Each entry in updates represents a delta applied atomically per row.
	// An empty input results in a no-op and no error.
	//
	// Behavior is undefined if ids do not exist or refer to non-published posts.
	UpdateViewCounts(ctx context.Context, updates map[uint]int64) error

	GetPublishedSiteMap(ctx context.Context) ([]*entity.Post, error)

	GetPublishedMeta(ctx context.Context) ([]*entity.Post, error)
}

type postRepo struct {
	ds *datastore.DataStore
}

// NewPostRepo returns a PostRepo backed by the provided datastore.
//
// The returned implementation is safe for concurrent use provided the underlying
// datastore client is concurrency-safe.
func NewPostRepo(ds *datastore.DataStore) PostRepo {
	return &postRepo{ds: ds}
}

// GetAllPublished returns all published, non-deleted posts with selected fields
// and eagerly loaded relations required for listing views.
//
// The query avoids loading full entities to reduce I/O overhead. Ordering is by
// publication time descending. Related author, categories, and tags are preloaded.
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
			post.ByPublishedAt(
				sql.OrderDesc(),
			),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.ToPosts(ps), nil
}

// GetPublishedByID returns a single published, non-deleted post with selected
// fields and author information.
//
// Only a subset of fields is loaded to minimize query cost. If no matching post
// exists, the underlying query returns an error.
func (r *postRepo) GetPublishedByID(ctx context.Context, id uint) (*entity.Post, error) {
	entPost, err := r.ds.Client(ctx).Post.
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
		return nil, err
	}

	return mapper.ToPost(entPost), nil
}

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
		return nil, err
	}

	return mapper.ToPosts(ps), nil
}

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
			post.ByPublishedAt(
				sql.OrderDesc(),
			),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.ToPosts(ps), nil
}

// UpdateViewCounts performs a bulk, per-row increment of view counts using a
// single SQL statement.
//
// Each id is updated with its corresponding delta via a CASE expression, ensuring
// atomicity at the row level. This avoids multiple round-trips but generates a
// query proportional to the number of ids, which may impact performance for large
// batches.
//
// Passing an empty map results in a no-op. The method does not validate existence
// of ids; missing rows are silently ignored by the underlying UPDATE.
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
	return err
}
