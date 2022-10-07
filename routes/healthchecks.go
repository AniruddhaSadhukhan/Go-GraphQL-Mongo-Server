package routes

import (
	"encoding/json"
	"go-graphql-mongo-server/logger"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(map[string]any{"service": "go-graphql-mongo-server-api", "ok": true}); err != nil {
		logger.Log.Fatal("encoding failed : %v", err)
	}
}
