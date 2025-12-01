package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type WordList struct {
	ent.Schema
}

func (WordList) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.Time("last_updated_at").Default(time.Now()),
	}
}

func (WordList) Edges() []ent.Edge{
	return []ent.Edge{
		edge.To("words", Word.Type),
		edge.To("groups", Group.Type),
		edge.From("school", School.Type).Ref("custom_word_lists").Unique().Immutable(),
		edge.From("competitions", Competition.Type).Ref("word_lists"),
	}
}