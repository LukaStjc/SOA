package configurations

import "os"

type Configurations struct {
	Port              string
	TourDBHost        string
	TourDBPort        string
	UserServiceDomain string
	UserServicePort   string
}

func NewConfigurations() *Configurations {
	configurations := &Configurations{
		Port:              os.Getenv("TOUR_SERVICE_PORT"),
		TourDBHost:        os.Getenv("TOUR_DB_HOST"),
		TourDBPort:        os.Getenv("TOUR_DB_PORT"),
		UserServiceDomain: os.Getenv("USER_SERVICE_DOMAIN"),
		UserServicePort:   os.Getenv("USER_SERVICE_PORT"),
	}

	configurations.initializeEnvironmentVariables()

	return configurations
}

func (configurations *Configurations) initializeEnvironmentVariables() {
	if configurations.Port == "" {
		configurations.Port = "3003"
	}
	if configurations.TourDBHost == "" {
		configurations.TourDBHost = "localhost"
	}
	if configurations.TourDBPort == "" {
		configurations.TourDBPort = "5432"
	}
	if configurations.UserServiceDomain == "" {
		configurations.UserServiceDomain = "localhost"
	}
	if configurations.UserServicePort == "" {
		configurations.UserServicePort = "3000"
	}

}
