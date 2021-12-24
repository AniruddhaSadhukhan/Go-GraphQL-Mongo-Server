package mutation

import (
	"go-graphql-mongo-server/common"
	"go-graphql-mongo-server/gqlhandler/schema"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
)

var UserMutation = &graphql.Field{
	Name:        "AddUsers",
	Type:        graphql.NewList(schema.UserSchema),
	Description: "Create multiple users",
	Args: graphql.FieldConfigArgument{
		"input": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(schema.UserSchemaInput)),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		userName := common.GetUserName(p)
		logger.Log.Info("Mutation: User called by " + userName)

		var userInput []models.User

		//Decode input to UserInput
		mapstructure.Decode(p.Args["input"], &userInput)

		// Create []interface{} from userInput
		var userInputInterface []interface{}
		for _, user := range userInput {
			userInputInterface = append(userInputInterface, user)
		}

		//Insert user
		err := models.InsertMany(models.UserCollection, userInputInterface, p.Context)
		return userInput, err

	},
}
