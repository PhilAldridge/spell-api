package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Group struct {
	ent.Schema
}

func (Group) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("join_code").Unique(),
		field.Time("join_code_valid_until_timestamp").Default(time.Now().Add(time.Hour*24*7)),
		field.Time("last_updated_at_timestamp").Default(time.Now()),
		field.Int("school_id"),
	}
}

func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("school", School.Type).Ref("groups").Unique().Required().Field("school_id"),
		edge.To("users", User.Type),
		edge.From("word_lists", WordList.Type).Ref("groups"),
		edge.From("competitions",Competition.Type).Ref("groups"),
	}
}