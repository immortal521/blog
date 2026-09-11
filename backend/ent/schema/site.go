package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Site holds the schema definition for the Site entity.
type Site struct {
	ent.Schema
}

// Fields of the Site.
func (Site) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Default(1).
			Immutable(),
		field.String("name"),
		field.String("logo").
			Optional(),
		field.String("greeting").
			Optional(),
		field.String("description").
			Optional(),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now),
	}
}

// Edges of the Site.
func (Site) Edges() []ent.Edge {
	return nil
}
