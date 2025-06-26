package config

import "github.com/spf13/viper"

type AppConfig struct {
	App struct {
		Name  string
		Debug bool
		Port  uint16
	}

	Database struct {
		DSN string
	}

	Redis struct {
		URL string
	}
}

var cfg *AppConfig

func Config() *AppConfig {
	if cfg == nil {
		loadConfig()
	}

	return cfg
}

func loadConfig() {
	readEnvVariables()

	cfg = &AppConfig{}

	// App.
	cfg.App.Name = viper.GetString("APP_NAME")
	cfg.App.Debug = viper.GetBool("APP_DEBUG")
	cfg.App.Port = viper.GetUint16("APP_PORT")

	// Database.
	cfg.Database.DSN = viper.GetString("DATABASE_DSN")

	// Redis.
	cfg.Redis.URL = viper.GetString("REDIS_URL")
}

func readEnvVariables() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
}
