package main

import (
	"context"
	"fmt"
	"go-graphql-mongo-server/config"
	"go-graphql-mongo-server/dbmigration"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"
	"go-graphql-mongo-server/routes"
	"net/http"
	"os"
	"strings"

	"github.com/adammck/venv"
	"github.com/gorilla/handlers"
)

type Service struct {
	HTTPServer http.Server
}

func main() {

	// Initialize config
	config.InitializeConfig(venv.OS())

	// Initialize Logger
	logger.Initialize()

	// Initialize Database
	models.InitializeDB()

	//Run DB Migration
	err := dbmigration.RunDbSchemaMigration("")
	if err != nil {
		panic(fmt.Errorf("error while running migration : %v", err))
	}

	// Configure Server
	service := &Service{}
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	// Get CORS Allowe Origins from config
	str := config.ConfigManager.CORSAllowOrigins
	CORSAllowOrigins := []string{}
	if str != "" {
		CORSAllowOrigins = strings.Split(str, ",")
	}
	origins := handlers.AllowedOrigins(CORSAllowOrigins)
	logger.Log.Info("CORS Allow Origins: ", CORSAllowOrigins)

	router := routes.NewRouter()
	service.HTTPServer = http.Server{
		Addr:    ":" + config.ConfigManager.ServicePort,
		Handler: handlers.CORS(origins, headers, methods)(router),
	}

	err = service.Run()
	if err != nil {
		logger.Log.Fatal("Error starting the server", err)
	}

}

func (s *Service) Run() error {
	logger.Log.Info("Starting the server on port ", config.ConfigManager.ServicePort)

	if checkIfHttpsCertExists() {
		logger.Log.Info("HTTPS Enabled")
		return s.HTTPServer.ListenAndServeTLS(config.ConfigManager.HttpsCert.CertFilePath, config.ConfigManager.HttpsCert.KeyFilePath)
	} else {
		logger.Log.Info("HTTPS Disabled")
		return s.HTTPServer.ListenAndServe()
	}
}

func checkIfHttpsCertExists() bool {
	if config.ConfigManager.HttpsCert.CertFilePath == "" || config.ConfigManager.HttpsCert.KeyFilePath == "" {
		return false
	}

	if _, err := os.Stat(config.ConfigManager.HttpsCert.CertFilePath); os.IsNotExist(err) {
		logger.Log.Error("HTTPS Cert File does not exist: ", config.ConfigManager.HttpsCert.CertFilePath)
		return false
	}

	if _, err := os.Stat(config.ConfigManager.HttpsCert.KeyFilePath); os.IsNotExist(err) {
		logger.Log.Error("HTTPS Key File does not exist: ", config.ConfigManager.HttpsCert.KeyFilePath)
		return false
	}

	config.ConfigManager.HttpsCert.HttpsEnabled = true
	return true
}

func (s *Service) Shutdown() error {
	if err := s.HTTPServer.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
