package main

import (
	"go-graphql-mongo-server/config"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/routes"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
)

type Service struct {
	HTTPServer http.Server
}

func main() {
	config.InitializeConfig()
	// models.InitializeDB()

	service := &Service{}
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	// Get CORS Allowe Origins from config
	//If not present, allow all
	str := config.ConfigManager.CORSAllowOrigins
	CORSAllowOrigins := []string{"*"}
	if str != "" {
		CORSAllowOrigins = strings.Split(str, ",")
	}
	origins := handlers.AllowedOrigins(CORSAllowOrigins)
	logger.Log.Info("CORS Allow Origins: ", CORSAllowOrigins)

	router := routes.NewRouter()
	service.HTTPServer = http.Server{
		Addr:    ":" + config.ConfigManager.Port,
		Handler: handlers.CORS(origins, headers, methods)(router),
	}

	logger.Log.Info("Starting the server on port ", config.ConfigManager.Port)
	err := service.HTTPServer.ListenAndServe()
	if err != nil {
		logger.Log.Fatal("Error starting the server", err)
	}

}
