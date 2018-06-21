package pkg

import "os"

type Configuration struct {
	DB string
	ConnectionString string
}

func envOrDefault(env string, def string) string{
	v, exists := os.LookupEnv(env)
	if exists {
		return v
	}
	return def
}

var config *Configuration

func Config() *Configuration{
	if config == nil {
		config = &Configuration{}
		config.DB = envOrDefault("VEIL_DB", "MYSQL")
		config.ConnectionString = envOrDefault("VEIL_DB_CONN", "root:root@tcp(127.0.0.1:3306)/veil")
	}
	return config
}

