package common

import (
	"encoding/json"
	"go-graphql-mongo-server/config"
	"go-graphql-mongo-server/logger"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")

	// Set HSTS header is HTTPS is enabled
	if config.ConfigManager.HttpsCert.HttpsEnabled {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}

	if config.ConfigManager.CORSAllowOrigins != "" {
		w.Header().Set("Access-Control-Allow-Origin", config.ConfigManager.CORSAllowOrigins)
	}
	w.WriteHeader(code)

	if code != http.StatusNoContent {
		_, err := w.Write(response)

		if err != nil {
			logger.Log.Errorf("Error in writing response %+v", err)
		}
	}

}
