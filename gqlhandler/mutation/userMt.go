package mutation

import (
	"go-graphql-mongo-server/common"
	"go-graphql-mongo-server/gqlhandler/schema"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"
	"go-graphql-mongo-server/telemetry"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
)

var UserMutation = &graphql.Field{
	Name:        "AddUsers",
	Type:        graphql.NewList(schema.UserSchema),
	Description: "Create multiple users",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(schema.UserInputSchema)),
		},
	},
	Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {

		if !common.IsInternalUser(p) {
			return nil, common.ErrUnauthorized
		}

		_, err := common.Sanitize(p.Args)
		if err != nil {
			return nil, err
		}

		defer telemetry.LogGraphQlCall(p, e)

		var userInput []models.User

		//Decode input to UserInput
		err = mapstructure.Decode(p.Args["input"], &userInput)
		if err != nil {
			logger.Log.Error(err)
		}

		// Create []interface{} from userInput
		var userInputInterface []interface{}
		for _, user := range userInput {
			userInputInterface = append(userInputInterface, user)
		}

		//Insert user
		err = models.InsertMany(p.Context, models.UserCollection, userInputInterface)
		return userInput, err

	},
}
