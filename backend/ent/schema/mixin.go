// Package schema
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type TimeMixin struct {
	mixin.Schema
}

func boolPtr(b bool) *bool {
	return &b
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").
			Unique().
			Immutable().
			Annotations(
				entsql.Annotation{
					Incremental: boolPtr(true),
				}),
		field.
			Time("created_at").
			Default(time.Now),
		field.
			Time("updated_at").
			Default(time.Now),
		field.
			Time("deleted_at").
			Optional().
			Nillable(),
	}
}
