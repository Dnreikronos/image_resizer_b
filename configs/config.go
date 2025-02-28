package configs

import "github.com/spf13/viper"

var cfg *config

type config struct {
	API APIConfig
	DB  DBConfig
}

type APIConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	TimeZone string
}

func init() {
	viper.SetDefault("api.port", "9090")
	viper.SetDefault("database.host", "${DB_HOST}")
	viper.SetDefault("database.port", "${DB_PORT}")
	viper.SetDefault("database.user", "${DB_USER}")
	viper.SetDefault("database.pass", "${POSTGRES_PASSWORD}")
	viper.SetDefault("database.database", "${POSTGRES_DB}")
	viper.SetDefault("database.timezone", "${POSTGRES_TIME_ZONE}")
}
