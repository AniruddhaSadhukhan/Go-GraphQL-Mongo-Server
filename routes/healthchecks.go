package routes

import (
	"context"
	"encoding/json"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"
	"net/http"
	"strings"
)

var isDbConnOk bool = true

func CheckDbConnection() {
	err := models.PingDatabase(context.TODO())
	if err != nil {
		logger.Log.Error(err)
		if strings.Contains(err.Error(), "auth error") {
			// If DB creds are rotated while this service is running
			logger.Log.Warn("MongoDB Auth error, trying to restart")
			isDbConnOk = false
		}
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if !isDbConnOk {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(map[string]any{"service": "go-graphql-mongo-server-api", "ok": isDbConnOk}); err != nil {
		logger.Log.Fatal("encoding failed : %v", err)
	}
}
