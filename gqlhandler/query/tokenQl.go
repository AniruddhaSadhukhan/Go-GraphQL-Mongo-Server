package query

import (
	"go-graphql-mongo-server/common"
	"go-graphql-mongo-server/gqlhandler/schema"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
)

var TokenQuery = &graphql.Field{
	Name:        "Tokens",
	Type:        graphql.NewList(schema.TokenSchema),
	Description: "Get all tokens for current user",
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		userName := common.GetUserName(p)
		logger.Log.Info("Query: Tokens called by " + userName)

		//Get Tokens from db
		var tokens []models.Token
		err := models.FindAll(models.TokenCollection, bson.M{"userName":userName}, nil, &tokens, p.Context)
		return tokens, err

	},
}
