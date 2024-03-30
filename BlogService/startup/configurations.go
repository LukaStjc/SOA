package configurations

import "os"

type Configurations struct {
	Port              string
	BlogDBHost        string
	BlogDBPort        string
	UserServiceDomain string
	UserServicePort   string
}

func NewConfigurations() *Configurations {
	configurations := &Configurations{
		Port:              os.Getenv("BLOG_SERVICE_PORT"),
		BlogDBHost:        os.Getenv("BLOG_DB_HOST"),
		BlogDBPort:        os.Getenv("BLOG_DB_PORT"),
		UserServiceDomain: os.Getenv("USER_SERVICE_DOMAIN"),
		UserServicePort:   os.Getenv("USER_SERVICE_PORT"),
	}

	configurations.initializeEnvironmentVariables()

	return configurations
}

func (configurations *Configurations) initializeEnvironmentVariables() {
	if configurations.Port == "" {
		configurations.Port = "8081"
	}
	if configurations.BlogDBHost == "" {
		configurations.BlogDBHost = "localhost"
	}
	if configurations.BlogDBPort == "" {
		configurations.BlogDBPort = "5432"
	}
	if configurations.UserServiceDomain == "" {
		configurations.UserServiceDomain = "localhost"
	}
	if configurations.UserServicePort == "" {
		configurations.UserServicePort = "3000"
	}
}
