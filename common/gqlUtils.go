package common

import (
	"fmt"
	"go-graphql-mongo-server/models"

	"github.com/graphql-go/graphql"
)

var ErrUnauthorized = fmt.Errorf("unauthorized access: you don't have permission for this operation")

func GetUserName(p graphql.ResolveParams) string {
	return p.Context.Value("User").(string)
}

func IsValidUser(p graphql.ResolveParams) bool {
	return GetUserName(p) != ""
}

func IsInternalUser(p graphql.ResolveParams) bool {
	return GetUserName(p) == models.InternalUser
}
