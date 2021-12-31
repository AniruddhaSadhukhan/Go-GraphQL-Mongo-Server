package gqlhandler

import (
	"net/http"
)

func GraphiqlHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "resources/graphiql.html")
}
