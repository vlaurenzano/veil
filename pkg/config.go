package pkg

import (
	"os"
	"strconv"
	"log"
)

type Configuration struct {
	DB string
	ConnectionString string
	LimitDefault int
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
		limit, err := strconv.Atoi(envOrDefault("VEIL_LIMIT_DEFAULT", "30"))
		if err != nil {
			log.Fatal("cofiguration error: invalid limit value")
		}
		config.LimitDefault = limit
	}

	return config
}

