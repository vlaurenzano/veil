package pkg

import (
	"os"
	"strconv"
	"log"
	"strings"
)

type Configuration struct {
	DB               string //what db we are using
	ConnectionString string //our storage connection string
	LimitDefault     int    //our default upper limit\

	//our permissions
	GetPermissions    map[string]string
	PutPermissions    map[string]string
	PostPermissions   map[string]string
	DeletePermissions map[string]string
}

func envOrDefault(env string, def string) string {
	v, exists := os.LookupEnv(env)
	if exists {
		return v
	}
	return def
}

func parsePermissionConf(pStr string) map[string]string {
	conf := make(map[string] string)
	for _, p := range strings.Split(pStr, ";"){
		rv := strings.Split(p, ":")
		conf[rv[0]] = rv[1]
	}
	return conf
}

var config *Configuration

func Config() *Configuration {
	if config == nil {
		config = &Configuration{}
		config.DB = envOrDefault("VEIL_DB", "MYSQL")
		config.ConnectionString = envOrDefault("VEIL_DB_CONN", "root:root@tcp(127.0.0.1:3306)/veil")

		limit, err := strconv.Atoi(envOrDefault("VEIL_LIMIT_DEFAULT", "30"))
		if err != nil {
			log.Fatal("cofiguration error: invalid limit value")
		}
		config.LimitDefault = limit

		config.GetPermissions = parsePermissionConf(envOrDefault("VEL_GET_PERMISSIONS", "global:allow"))
		config.PutPermissions = parsePermissionConf(envOrDefault("VEL_PUT_PERMISSIONS", "global:deny"))
		config.PostPermissions = parsePermissionConf(envOrDefault("VEL_POST_PERMISSIONS", "global:deny"))
		config.DeletePermissions = parsePermissionConf(envOrDefault("VEL_DELETE_PERMISSIONS", "global:deny"))
	}

	return config
}
