package mutation

import (
	"fmt"
	"go-graphql-mongo-server/auth"
	"go-graphql-mongo-server/common"
	"go-graphql-mongo-server/gqlhandler/schema"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
)

var CreateTokenMutation = &graphql.Field{
	Name:        "CreateToken",
	Type:        schema.TokenSchema,
	Description: "Create a long lived personal access token",
	Args: graphql.FieldConfigArgument{
		"tokenName": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"expiresAt": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		if common.IsValidUser(p) {
			return nil, common.ErrUnauthorized
		}

		_, err := common.Sanitize(p.Args)
		if err != nil {
			return nil, err
		}

		userName := common.GetUserName(p)
		logger.Log.Info("Mutation: Create Token called by " + userName)

		var token models.Token
		token.UserName = userName
		//Decode input to token
		mapstructure.Decode(p.Args, &token)

		err = auth.GenerateToken(&token)
		if err != nil {
			return nil, err
		}

		err = models.Insert(models.TokenCollection, token, p.Context)
		if err != nil && strings.Contains(err.Error(), "duplicate key error") {
			return nil, fmt.Errorf("a token with same name already exists")
		}

		return token, err

	},
}

var RevokeTokenMutation = &graphql.Field{
	Name:        "RevokeToken",
	Type:        graphql.Boolean,
	Description: "Revoke a long lived personal access token",
	Args: graphql.FieldConfigArgument{
		"tokenName": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		if common.IsValidUser(p) {
			return false, common.ErrUnauthorized
		}

		_, err := common.Sanitize(p.Args)
		if err != nil {
			return nil, err
		}

		userName := common.GetUserName(p)
		logger.Log.Info("Mutation: Revoke Token called by " + userName)

		err = models.Delete(
			models.TokenCollection,
			bson.M{"tokenName": p.Args["tokenName"].(string), "userName": userName},
			p.Context,
		)
		if err != nil {
			return false, err
		}
		return true, nil

	},
}
