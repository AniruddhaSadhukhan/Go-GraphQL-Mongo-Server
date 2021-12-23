package routes

import (
	"go-graphql-mongo-server/gqlhandler"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Method     string
	Pattern    string
	Handler    http.HandlerFunc
	Middleware mux.MiddlewareFunc
}

var routes []Route

func init() {
	setGraphQLRoutes()
}

func register(method, pattern string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	pattern = "/api/v1" + pattern
	routes = append(routes, Route{method, pattern, handler, middleware}, Route{"OPTIONS", pattern, handler, middleware})
}

func setGraphQLRoutes() {
	register("POST", "/graphql", gqlhandler.GraphqlHandler, nil)
	register("GET", "/graphiql", gqlhandler.GraphiqlHandler, nil)
}

func NewRouter() (router *mux.Router) {
	router = mux.NewRouter()
	for _, route := range routes {
		r := router.
			Methods(route.Method).
			Path(route.Pattern)

		if route.Middleware != nil {
			r.Handler(route.Middleware(route.Handler))
		} else {
			r.Handler(route.Handler)
		}
	}
	return
}
