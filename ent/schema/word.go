package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Word struct {
	ent.Schema
}

func (Word) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.Time("last_updated_at").Default(time.Now()),
	}
}

func (Word) Edges() []ent.Edge{
	return []ent.Edge{
		edge.From("word_lists", WordList.Type).Ref("words"),
		edge.From("school",School.Type).Ref("custom_words").Unique().Immutable(),
		edge.From("results",Result.Type).Ref("word"),
	}
}