package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("password_hash").Sensitive().NotEmpty(),
		field.String("email").Unique(),
		field.Enum("account_type").Optional().Values("student","admin").Default("student"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("schools", School.Type).Ref("admins"),
		edge.To("owned_school", School.Type).Unique(),
		edge.From("groups", Group.Type).Ref("users"),
		edge.To("results", Result.Type),
		edge.To("refresh_tokens", RefreshToken.Type),
	}
}