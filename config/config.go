package config

import "github.com/adammck/venv"

type Configurations struct {
	Database
	Auth
	HttpsCert
	CORSAllowOrigins  string
	ServicePort       string
	ApiLimitPerSecond string
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

type Auth struct {
	JWT_PrivateKey string
	SecretToken    string
}

type HttpsCert struct {
	HttpsEnabled bool
	CertFilePath string
	KeyFilePath  string
}

var env venv.Env
var ConfigManager Configurations

//Set environment variables
func InitializeConfig(e venv.Env) {
	env = e
	ConfigManager = readConfigValues()
}

func readConfigValues() Configurations {
	return Configurations{
		Database: Database{
			Host:               getEnvVariable("DB_HOST", ""),
			Port:               getEnvVariable("DB_PORT", ""),
			Name:               getEnvVariable("DB_NAME", ""),
			Username:           getEnvVariable("DB_USERNAME", ""),
			Password:           getEnvVariable("DB_PASSWORD", ""),
			InsecureSkipVerify: getEnvVariable("DB_INSECURE_SKIP_VERIFY", "false") == "true",
		},
		Auth: Auth{
			JWT_PrivateKey: getEnvVariable("JWT_PRIVATE_KEY", ""),
			SecretToken:    getEnvVariable("SECRET_TOKEN", ""),
		},
		HttpsCert: HttpsCert{
			CertFilePath: getEnvVariable("HTTPS_CERT_FILE_PATH", ""),
			KeyFilePath:  getEnvVariable("HTTPS_KEY_FILE_PATH", ""),
		},
		CORSAllowOrigins:  getEnvVariable("CORS_ALLOW_ORIGINS", ""),
		ServicePort:       getEnvVariable("PORT", "8080"),
		ApiLimitPerSecond: getEnvVariable("API_LIMIT_PER_SECOND", "500"),
	}
}

func getEnvVariable(key string, defaultValue string) string {
	value := env.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
