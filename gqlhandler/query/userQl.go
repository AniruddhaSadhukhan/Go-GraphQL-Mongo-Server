package query

import (
	"go-graphql-mongo-server/common"
	"go-graphql-mongo-server/gqlhandler/schema"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"

	"github.com/graphql-go/graphql"
)

var UsersQuery = &graphql.Field{
	Name:        "Users",
	Type:        graphql.NewList(schema.UserSchema),
	Description: "Get all users",
	Args: graphql.FieldConfigArgument{
		// Add optional filters for users
		"id": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
		"name": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
		"dob": &graphql.ArgumentConfig{
			Type: graphql.DateTime,
		},
		"isVerified": &graphql.ArgumentConfig{
			Type: graphql.Boolean,
		},
		"subscription": &graphql.ArgumentConfig{
			Type: schema.SubscriptionTypeEnum,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		if !common.IsValidUser(p) {
			return nil, common.ErrUnauthorized
		}

		_, err := common.Sanitize(p.Args)
		if err != nil {
			return nil, err
		}

		userName := common.GetUserName(p)
		logger.Log.Info("Query: Users called by " + userName)

		//Get Users from db
		var users []models.User
		err = models.FindAll(models.UserCollection, p.Args, nil, &users, p.Context)
		return users, err

	},
}
