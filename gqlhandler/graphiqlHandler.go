package gqlhandler

import (
	"go-graphql-mongo-server/config"
	"net/http"
)

func GraphiqlHandler(w http.ResponseWriter, r *http.Request) {
	// Set HSTS header is HTTPS is enabled
	if config.ConfigManager.HttpsCert.HttpsEnabled {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}
	http.ServeFile(w, r, "resources/graphiql.html")
}
