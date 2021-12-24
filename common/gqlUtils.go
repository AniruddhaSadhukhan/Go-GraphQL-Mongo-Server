package common

import "github.com/graphql-go/graphql"

func GetUserName(p graphql.ResolveParams) string {
	return p.Context.Value("User").(string)
}