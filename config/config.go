package config

import "github.com/adammck/venv"

type Configurations struct {
	Database
	Auth
	HTTPSCert
	ProductionMode    bool
	CORSAllowOrigins  string
	ServicePort       string
	APILimitPerSecond string
	TelemetryURL      string
	PlatformName      string
	Env               string
	ComponentName     string
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
	// In House Token related configs
	JWTInHousePrivateKey string

	// Internal auth config
	SecretToken string

	// OIDC related configs
	OidcEnabled  bool
	OidcURL      string
	ClientID     string
	ClientSecret string
}

type HTTPSCert struct {
	HTTPSEnabled bool
	CertFilePath string
	KeyFilePath  string
}

var env venv.Env
var Store Configurations

// Set environment variables
func InitializeConfig(e venv.Env) {
	env = e
	Store = readConfigValues()
}

func readConfigValues() Configurations {
	//nolint:goconst
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
			JWTInHousePrivateKey: getEnvVariable("JWT_PRIVATE_KEY", ""),
			SecretToken:          getEnvVariable("SECRET_TOKEN", ""),
			OidcURL:              getEnvVariable("OIDC_URL", ""),
			OidcEnabled:          getEnvVariable("OIDC_URL", "") != "",
			ClientID:             getEnvVariable("CLIENT_ID", ""),
			ClientSecret:         getEnvVariable("CLIENT_SECRET", ""),
		},
		HTTPSCert: HTTPSCert{
			CertFilePath: getEnvVariable("HTTPS_CERT_FILE_PATH", ""),
			KeyFilePath:  getEnvVariable("HTTPS_KEY_FILE_PATH", ""),
			HTTPSEnabled: getEnvVariable("HTTPS_CERT_FILE_PATH", "") != "",
		},
		ProductionMode:    getEnvVariable("PRODUCTION_MODE", "true") == "true",
		CORSAllowOrigins:  getEnvVariable("CORS_ALLOW_ORIGINS", ""),
		ServicePort:       getEnvVariable("PORT", "8080"),
		APILimitPerSecond: getEnvVariable("API_LIMIT_PER_SECOND", "500"),
		TelemetryURL:      getEnvVariable("TELEMETRY_URL", ""),
		PlatformName:      getEnvVariable("PLATFORM_NAME", "Ani Platform"),
		Env:               getEnvVariable("ENV", ""),
		ComponentName:     getEnvVariable("COMPONENT_NAME", "Go GraphQl Mongo Server"),
	}
}

func getEnvVariable(key string, defaultValue string) string {
	value := env.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
