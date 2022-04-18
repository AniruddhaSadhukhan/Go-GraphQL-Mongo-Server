package routes

import (
	"go-graphql-mongo-server/auth"
	"go-graphql-mongo-server/gqlhandler"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Method     string
	Pattern    string
	Handler    http.HandlerFunc
	Middleware []mux.MiddlewareFunc
}

var routes []Route

func initializeRoutes() {
	createLimiterMiddleware()
	setGraphQLRoutes()
}

func register(method, pattern string, handler http.HandlerFunc, middleware []mux.MiddlewareFunc) {
	pattern = "/api/v1" + pattern
	routes = append(routes, Route{method, pattern, handler, middleware}, Route{"OPTIONS", pattern, handler, middleware})
}

func setGraphQLRoutes() {

	register(
		"POST",
		"/graphql",
		gqlhandler.GraphqlHandler,
		[]mux.MiddlewareFunc{
			limiterMiddleware.Handle,
			auth.AuthMiddleware,
		})

	register(
		"GET",
		"/graphiql",
		gqlhandler.GraphiqlHandler,
		nil,
	)
}

func NewRouter() (router *mux.Router) {
	initializeRoutes()

	router = mux.NewRouter()

	for _, route := range routes {
		r := router.
			Methods(route.Method).
			Path(route.Pattern)

		var handler http.Handler = route.Handler
		for i := len(route.Middleware) - 1; i >= 0; i-- {
			handler = route.Middleware[i](handler)
		}

		r.Handler(route.Handler)

	}
	return
}
