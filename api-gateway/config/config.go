package config

import "os"

type Config struct {
	Address            string
	AuthServiceAddress string
}

func GetConfig() Config {
	return Config{
		AuthServiceAddress: os.Getenv("AUTH_SERVICE_ADDRESS"),
		Address:            os.Getenv("GATEWAY_ADDRESS"),
	}
}
