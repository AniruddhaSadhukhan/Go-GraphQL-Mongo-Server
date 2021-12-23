package config

import "github.com/adammck/venv"

type Configurations struct {
	Database
	CORSAllowOrigins string
	Port             string
}

// Database configuration
type Database struct {
	Host     string
	Name     string
	User     string
	Password string
}

var env venv.Env
var ConfigManager Configurations

//Set environment variables
func InitializeConfig() {
	env = venv.OS()
	ConfigManager = readConfigValues()
}

func readConfigValues() Configurations {
	return Configurations{
		Database: Database{
			Host:     env.Getenv("DB_HOST"),
			Name:     env.Getenv("DB_NAME"),
			User:     env.Getenv("DB_USER"),
			Password: env.Getenv("DB_PASSWORD"),
		},
		CORSAllowOrigins: env.Getenv("CORS_ALLOW_ORIGINS"),
		Port:             env.Getenv("PORT"),
	}
}
