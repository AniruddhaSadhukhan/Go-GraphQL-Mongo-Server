package config

import "github.com/adammck/venv"

type Configurations struct {
	Database
	CORSAllowOrigins string
	Port             string
	SecretToken      string
	JWT_PrivateKey   string
}

// Database configuration
type Database struct {
	Host               string
	Port               string
	Name               string
	Username           string
	Password           string
	InsecureSkipVerify bool
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
			Host:               env.Getenv("DB_HOST"),
			Port:               env.Getenv("DB_PORT"),
			Name:               env.Getenv("DB_NAME"),
			Username:           env.Getenv("DB_USERNAME"),
			Password:           env.Getenv("DB_PASSWORD"),
			InsecureSkipVerify: env.Getenv("DB_INSECURE_SKIP_VERIFY") == "true",
		},
		CORSAllowOrigins: env.Getenv("CORS_ALLOW_ORIGINS"),
		Port:             env.Getenv("PORT"),
		SecretToken:      env.Getenv("SECRET_TOKEN"),
		JWT_PrivateKey:   env.Getenv("JWT_PRIVATE_KEY"),
	}
}
