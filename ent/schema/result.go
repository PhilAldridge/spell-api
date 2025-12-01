package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Result struct {
	ent.Schema
}

func (Result) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").NotEmpty(),
		field.Bool("correct").Default(false),
		field.Float("time_taken_in_seconds").Positive(),
		field.Time("tested_at_timestamp").Default(time.Now()),
	}
}

func (Result) Edges() []ent.Edge{
	return []ent.Edge{
		edge.To("user", User.Type).Unique().Immutable(),
		edge.To("word", Word.Type).Unique().Immutable(),
	}
}

func (Result) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tested_at_timestamp"),
		index.Fields("user"),
	}
}

func (Result) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Options: "using columnar",
		},
	}
}