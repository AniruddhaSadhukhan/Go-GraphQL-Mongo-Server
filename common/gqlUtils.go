package common

import (
	"fmt"
	"go-graphql-mongo-server/models"
	"regexp"

	"github.com/graphql-go/graphql"
)

var ErrUnauthorized = fmt.Errorf("unauthorized access: you don't have permission for this operation")

const userNameRegex = "^[a-zA-Z]{2}\\d{5}$"

func GetUserName(p graphql.ResolveParams) string {
	return p.Context.Value(models.UserContextKey).(string)
}

func GetTokenUserName(p graphql.ResolveParams) string {
	customUserName, customUserNamePresent := p.Args["userName"].(string)
	if customUserNamePresent && IsInternalUser(p) && customUserName != "" {
		return "@" + customUserName + "@"
	}
	return p.Context.Value(models.UserContextKey).(string)
}

func IsValidUser(p graphql.ResolveParams) bool {
	return GetUserName(p) != ""
}

func IsInternalUser(p graphql.ResolveParams) bool {
	return GetUserName(p) == models.InternalUser
}

func GetUserType(p graphql.ResolveParams) string {
	isUser, _ := regexp.Match(userNameRegex, []byte(GetUserName(p)))
	if isUser {
		return "User"
	}

	if IsInternalUser(p) {
		return "Internal Service"
	}

	return "Service"
}
