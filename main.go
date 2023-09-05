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
	"time"

	"github.com/adammck/venv"
	"github.com/gorilla/handlers"
	"github.com/robfig/cron/v3"
	"github.com/urfave/negroni"
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

	// Initialize Cron
	setUpCronJobs()

	// Configure Server
	service := &Service{}
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	// Get CORS Allowed Origins from config
	str := config.Store.CORSAllowOrigins
	CORSAllowOrigins := []string{}
	if str != "" {
		CORSAllowOrigins = strings.Split(str, ",")
	}
	origins := handlers.AllowedOrigins(CORSAllowOrigins)
	logger.Log.Info("CORS Allow Origins: ", CORSAllowOrigins)

	router := routes.NewRouter()
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(router)

	service.HTTPServer = http.Server{
		Addr:              ":" + config.Store.ServicePort,
		Handler:           handlers.CORS(origins, headers, methods)(n),
		ReadHeaderTimeout: 10 * time.Second,
		// Uncomment the next line to disable HTTP 2
		// TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	err = service.Run()
	if err != nil {
		logger.Log.Fatal("Error starting the server", err)
	}

}

func (s *Service) Run() error {
	logger.Log.Info("Starting the server on port ", config.Store.ServicePort)

	if checkIfHttpsCertExists() {
		logger.Log.Info("HTTPS Enabled")
		return s.HTTPServer.ListenAndServeTLS(config.Store.HTTPSCert.CertFilePath, config.Store.HTTPSCert.KeyFilePath)
	} else {
		logger.Log.Info("HTTPS Disabled")
		return s.HTTPServer.ListenAndServe()
	}
}

func checkIfHttpsCertExists() bool {
	if config.Store.HTTPSCert.CertFilePath == "" || config.Store.HTTPSCert.KeyFilePath == "" {
		return false
	}

	if _, err := os.Stat(config.Store.HTTPSCert.CertFilePath); os.IsNotExist(err) {
		logger.Log.Error("HTTPS Cert File does not exist: ", config.Store.HTTPSCert.CertFilePath)
		return false
	}

	if _, err := os.Stat(config.Store.HTTPSCert.KeyFilePath); os.IsNotExist(err) {
		logger.Log.Error("HTTPS Key File does not exist: ", config.Store.HTTPSCert.KeyFilePath)
		return false
	}

	config.Store.HTTPSCert.HTTPSEnabled = true
	return true
}

func (s *Service) Shutdown() error {
	if err := s.HTTPServer.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Set up all cron jobs
func setUpCronJobs() {
	cronJob := cron.New()
	_, err := cronJob.AddFunc("@every 5m", routes.CheckDbConnection)
	if err != nil {
		logger.Log.Error(err)
	}
	cronJob.Start()
}
