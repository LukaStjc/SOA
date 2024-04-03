package configurations

import "os"

type Configurations struct {
	Port                        string
	UserDBHost                  string
	UserDBPort                  string
	Secret                      string
	AuthenticationServiceDomain string
	AuthenticationServicePort   string
}

func NewConfigurations() *Configurations {
	configurations := &Configurations{
		Port:                        os.Getenv("USER_SERVICE_PORT"),
		UserDBHost:                  os.Getenv("USER_DB_HOST"),
		UserDBPort:                  os.Getenv("USER_DB_PORT"),
		Secret:                      os.Getenv("SECRET"),
		AuthenticationServiceDomain: os.Getenv("AUTHENTICATION_SERVICE_DOMAIN"),
		AuthenticationServicePort:   os.Getenv("AUTHENTICATION_SERVICE_PORT"),
	}

	configurations.initializeEnvironmentVariables()

	return configurations
}

func (configurations *Configurations) initializeEnvironmentVariables() {
	if configurations.Port == "" {
		configurations.Port = "3000"
	}
	if configurations.UserDBHost == "" {
		configurations.UserDBHost = "localhost"
	}
	if configurations.UserDBPort == "" {
		configurations.UserDBPort = "5432"
	}
	if configurations.Secret == "" {
		configurations.Secret = "SOAPROJEKAT"
	}
	if configurations.AuthenticationServiceDomain == "" {
		configurations.AuthenticationServiceDomain = "localhost"
	}
	if configurations.AuthenticationServicePort == "" {
		configurations.AuthenticationServicePort = "3001"
	}
}
