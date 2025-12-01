package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type School struct {
	ent.Schema
}

func (School) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.Time("last_updated_at").Default(time.Now()),
	}
}

func (School) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("admins",User.Type),
		edge.From("owner", User.Type).Ref("owned_school").Unique(),
		edge.To("groups", Group.Type),
		edge.To("custom_words", Word.Type),
		edge.To("custom_word_lists", WordList.Type),
	}
}