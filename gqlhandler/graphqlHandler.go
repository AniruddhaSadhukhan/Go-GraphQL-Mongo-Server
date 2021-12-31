package gqlhandler

import (
	"encoding/json"
	"go-graphql-mongo-server/common"
	"go-graphql-mongo-server/gqlhandler/mutation"
	"go-graphql-mongo-server/gqlhandler/query"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
)

var SchemaQl, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

var mutationMap = graphql.Fields{
	mutation.UserMutation.Name:        mutation.UserMutation,
	mutation.CreateTokenMutation.Name: mutation.CreateTokenMutation,
}
var queryMap = graphql.Fields{
	query.UsersQuery.Name: query.UsersQuery,
}

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name:   "Mutation",
	Fields: mutationMap,
})
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name:   "Query",
	Fields: queryMap,
})

func GraphqlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	queryBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError("Error in reading request body", err, w)
		return
	}

	requests, err := getRequest(queryBody)
	if err != nil {
		handleError("Error in parsing request body", err, w)
		return
	}

	var resultMap []*graphql.Result
	var errorCount int
	for _, request := range requests {
		result := graphql.Do(graphql.Params{
			Schema:         SchemaQl,
			RequestString:  request.Query,
			VariableValues: request.Variables,
			Context:        ctx,
		})

		resultMap = append(resultMap, result)
		if result.HasErrors() {
			errorCount++
		}

	}

	if len(requests) == errorCount {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if len(resultMap) > 0 {
		var response []byte
		if len(resultMap) == 1 {
			response, _ = json.Marshal(resultMap[0])
		} else {
			response, _ = json.Marshal(resultMap)
		}

		_, err = w.Write(response)
		if err != nil {
			logger.Log.Errorf("Error in writing response %+v", err)
		}
	}

}

func getRequest(queryBody []byte) ([]models.GQLRequestBody, error) {
	var requests []models.GQLRequestBody
	var err error
	queryBodyString := string(queryBody)
	if queryBodyString[0] == '[' {
		err = json.Unmarshal(queryBody, &requests)
		if err != nil {
			return nil, err
		}
	} else {
		jsonMap := make(map[string]interface{})

		err = json.Unmarshal(queryBody, &jsonMap)
		if err != nil {
			return nil, err
		}
		variables := make(map[string]interface{})
		if jsonMap["variables"] != nil {
			variables = jsonMap["variables"].(map[string]interface{})
		}

		requests = append(requests, models.GQLRequestBody{
			Query:     jsonMap["query"].(string),
			Variables: variables,
		})
	}

	return requests, nil
}

func handleError(text string, err error, w http.ResponseWriter) {
	logger.Log.Errorf("%v : %+v", text, err)
	common.RespondWithJSON(w, http.StatusBadRequest, `{"errors": [{"message": "`+err.Error()+`"}]}`)
}
