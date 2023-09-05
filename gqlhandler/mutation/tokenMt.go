package mutation

import (
	"fmt"
	"go-graphql-mongo-server/auth"
	"go-graphql-mongo-server/common"
	"go-graphql-mongo-server/gqlhandler/schema"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"
	"go-graphql-mongo-server/telemetry"
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
		"userName": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "This username will only be used when Admin runs this mutation. Otherwise this will be ignored.",
		},
	},
	Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {

		if common.IsValidUser(p) {
			return nil, common.ErrUnauthorized
		}

		_, err := common.Sanitize(p.Args)
		if err != nil {
			return nil, err
		}

		defer telemetry.LogGraphQlCall(p, e)

		userName := common.GetTokenUserName(p)

		var token models.Token
		//Decode input to token
		err = mapstructure.Decode(p.Args, &token)
		if err != nil {
			logger.Log.Error(err)
		}
		token.UserName = userName

		err = auth.GenerateToken(&token)
		if err != nil {
			return nil, err
		}

		err = models.Insert(p.Context, models.TokenCollection, token)
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
		"userName": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "This username will only be used when Admin runs this mutation. Otherwise this will be ignored.",
		},
	},
	Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {

		if common.IsValidUser(p) {
			return false, common.ErrUnauthorized
		}

		_, err := common.Sanitize(p.Args)
		if err != nil {
			return nil, err
		}

		defer telemetry.LogGraphQlCall(p, e)

		userName := common.GetUserName(p)

		err = models.Delete(
			p.Context,
			models.TokenCollection,
			bson.M{"tokenName": p.Args["tokenName"].(string), "userName": userName},
		)
		if err != nil {
			return false, err
		}
		return true, nil

	},
}
