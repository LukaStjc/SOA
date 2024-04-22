package configurations

import "os"

type Configurations struct {
	Port                 string
	AuthenticationDBHost string
	AuthenticationDBPort string
	Secret               string
}

func NewConfigurations() *Configurations {
	configurations := &Configurations{
		Port:                 os.Getenv("AUTHENTICATION_SERVICE_PORT"),
		AuthenticationDBHost: os.Getenv("AUTHENTIFICATION_DB_HOST"),
		AuthenticationDBPort: os.Getenv("AUTHENTIFICATION_DB_PORT"),
		Secret:               os.Getenv("SECRET"),
	}

	configurations.initializeEnvironmentVariables()

	return configurations
}

func (configurations *Configurations) initializeEnvironmentVariables() {
	if configurations.Port == "" {
		configurations.Port = "3001"
	}
	if configurations.AuthenticationDBHost == "" {
		configurations.AuthenticationDBHost = "localhost"
	}
	if configurations.AuthenticationDBPort == "" {
		configurations.AuthenticationDBPort = "5432"
	}
	if configurations.Secret == "" {
		configurations.Secret = "SOAPROJEKAT"
	}
}
