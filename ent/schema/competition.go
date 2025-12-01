package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Competition struct {
	ent.Schema
}

func (Competition) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("type").Default("basic"),
		field.Time("start_timestamp").Default(time.Now()),
		field.Time("end_timestamp").Default(time.Now().Add(time.Hour*24*7)),
		field.Time("last_updated_at_timestamp").Default(time.Now()),
	}
}

func (Competition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("groups", Group.Type),
		edge.To("word_lists", WordList.Type),
	}
}