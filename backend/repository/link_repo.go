package repository

import (
	"context"
	"strings"

	"blog-server/datastore"
	"blog-server/ent/link"
	"blog-server/entity"
	"blog-server/mapper"

	"entgo.io/ent/dialect/sql"
)

type LinkRepo interface {
	GetAll(ctx context.Context) ([]*entity.Link, error)
	GetAllEnabled(ctx context.Context) ([]*entity.Link, error)
	Create(ctx context.Context, link *entity.Link) (*entity.Link, error)
	UpdateStatusBatch(ctx context.Context, updates map[uint]entity.LinkStatus) error
}

type linkRepo struct {
	ds *datastore.DataStore
}

func NewLinkRepo(ds *datastore.DataStore) LinkRepo {
	return &linkRepo{ds: ds}
}

func (r *linkRepo) GetAll(ctx context.Context) ([]*entity.Link, error) {
	links, err := r.ds.Client(ctx).Link.
		Query().
		Where(
			link.DeletedAtIsNil(),
		).
		Order(
			link.ByID(
				sql.OrderDesc(),
			),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.ToLinks(links), nil
}

func (r *linkRepo) GetAllEnabled(ctx context.Context) ([]*entity.Link, error) {
	links, err := r.ds.Client(ctx).Link.
		Query().
		Where(
			link.EnabledEQ(true),
			link.DeletedAtIsNil(),
		).
		Order(
			link.ByID(
				sql.OrderDesc(),
			),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.ToLinks(links), nil
}

func (r *linkRepo) Create(ctx context.Context, l *entity.Link) (*entity.Link, error) {
	c := r.ds.Client(ctx).Link.
		Create().
		SetName(l.Name).
		SetURL(l.URL)

	if l.Description != nil && strings.TrimSpace(*l.Description) != "" {
		c.SetDescription(*l.Description)
	}
	if l.Avatar != nil && strings.TrimSpace(*l.Avatar) != "" {
		c.SetAvatar(*l.Avatar)
	}

	link, err := c.Save(ctx)
	if err != nil {
		return nil, err
	}

	return mapper.ToLink(link), nil
}

func (r *linkRepo) UpdateStatusBatch(ctx context.Context, updates map[uint]entity.LinkStatus) error {
	if len(updates) == 0 {
		return nil
	}

	ids := make([]uint, 0, len(updates))
	for id := range updates {
		ids = append(ids, id)
	}

	err := r.ds.Client(ctx).Link.
		Update().
		Where(link.IDIn(ids...)).
		Modify(func(u *sql.UpdateBuilder) {
			caseExpr := func(b *sql.Builder) {
				b.WriteString("CASE ")
				b.Ident(link.FieldID)
				b.WriteString(" ")

				for _, id := range ids {
					b.WriteString("WHEN ")
					b.Arg(id)
					b.WriteString(" THEN ")
					b.Arg(updates[id])
					b.WriteString(" ")
				}
				b.WriteString("ELSE ")
				b.Ident(link.FieldStatus)
				b.WriteString(" END")
			}
			u.Set(link.FieldStatus, sql.ExprFunc(caseExpr))
		}).
		Exec(ctx)

	return err
}
