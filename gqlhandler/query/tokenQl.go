package query

import (
	"go-graphql-mongo-server/common"
	"go-graphql-mongo-server/gqlhandler/schema"
	"go-graphql-mongo-server/models"
	"go-graphql-mongo-server/telemetry"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
)

var TokenQuery = &graphql.Field{
	Name:        "Tokens",
	Type:        graphql.NewList(schema.TokenSchema),
	Description: "Get all tokens for current user",
	Args: graphql.FieldConfigArgument{
		"userName": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "This username will only be used when Admin runs this query. Otherwise this will be ignored.",
		},
	},
	Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {

		if !common.IsValidUser(p) {
			return nil, common.ErrUnauthorized
		}

		_, err := common.Sanitize(p.Args)
		if err != nil {
			return nil, err
		}

		defer telemetry.LogGraphQlCall(p, e)

		userName := common.GetTokenUserName(p)

		//Get Tokens from db
		var tokens []models.Token
		err = models.FindAll(p.Context, models.TokenCollection, bson.M{"userName": userName}, nil, &tokens)
		return tokens, err

	},
}
