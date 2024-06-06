package configurations

import "os"

type Configurations struct {
	Port                        string
	UserDBHost                  string
	UserDBPort                  string
	UserGraphDBHost             string
	UserGraphDBPort             string
	UserGraphDBUsername         string
	UserGraphDBPassword         string
	Secret                      string
	AuthenticationServiceDomain string
	AuthenticationServicePort   string
}

func NewConfigurations() *Configurations {
	configurations := &Configurations{
		Port:                        os.Getenv("USER_SERVICE_PORT"),
		UserDBHost:                  os.Getenv("USER_DB_HOST"),
		UserDBPort:                  os.Getenv("USER_DB_PORT"),
		UserGraphDBHost:             os.Getenv("USER_GRAPH_DB_HOST"),
		UserGraphDBPort:             os.Getenv("USER_GRAPH_DB_PORT"),
		UserGraphDBUsername:         os.Getenv("USER_GRAPH_DB_USERNAME"),
		UserGraphDBPassword:         os.Getenv("USER_GRAPH_DB_PASS"),
		Secret:                      os.Getenv("SECRET"),
		AuthenticationServiceDomain: os.Getenv("AUTHENTICATION_SERVICE_DOMAIN"),
		AuthenticationServicePort:   os.Getenv("AUTHENTICATION_SERVICE_PORT"),
	}

	configurations.initializeEnvironmentVariables()

	return configurations
}

func (configurations *Configurations) initializeEnvironmentVariables() {
	if configurations.Port == "" {
		configurations.Port = "3002"
	}
	if configurations.UserDBHost == "" {
		configurations.UserDBHost = "localhost"
	}
	if configurations.UserDBPort == "" {
		configurations.UserDBPort = "5432"
	}
	if configurations.UserGraphDBHost == "" {
		configurations.UserGraphDBHost = "neo4j"
	}
	if configurations.UserGraphDBPort == "" {
		configurations.UserGraphDBPort = "7687"
	}
	if configurations.UserGraphDBUsername == "" {
		configurations.UserGraphDBUsername = "neo4j"
	}
	if configurations.UserGraphDBPassword == "" {
		configurations.UserGraphDBPassword = "nekaSifra"
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
