package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type RefreshToken struct {
	ent.Schema
}

func (RefreshToken) Fields() []ent.Field {
	return []ent.Field {
		field.String("token_hash").NotEmpty(),
		field.Time("expires_at"),
		field.Bool("revoked").Default(false),
		field.Time("created_at").Default(time.Now()),
		field.Int("user_id"),
	}
}

func (RefreshToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("refresh_tokens").Unique().Required().Field("user_id"),
	}
}