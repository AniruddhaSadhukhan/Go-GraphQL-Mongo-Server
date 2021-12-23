package common

import (
	"encoding/json"
	"go-graphql-mongo-server/logger"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)

	if code != http.StatusNoContent {
		_, err := w.Write(response)

		if err != nil {
			logger.Log.Errorf("Error in writing response %+v", err)
		}
	}

}
