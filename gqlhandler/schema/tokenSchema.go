package schema

import "github.com/graphql-go/graphql"

var TokenSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Token",
		Fields: graphql.Fields{
			"tokenName": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"expiresAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"token": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
