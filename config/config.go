package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type AppConfig struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	DSN                 string        `mapstructure:"dsn"`
	MaxOpenConns        int           `mapstructure:"max_open_conns"`
	MaxIdleConns        int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime     time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdletime     time.Duration `mapstructure:"conn_max_idle_time"`
	HealthCheckInterval time.Duration `mapstructure:"health_check_interval"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"` // info, debug, warn, error
	Production bool   `mapstructure:"production"`
}

func LoadConfig() (config AppConfig, err error) {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Warning: .env file not found, relying on environment variables.")
		} else {
			return AppConfig{}, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.read_timeout", "5s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("server.idle_timeout", "60s")

	viper.SetDefault("database.dsn", "postgres://postgres:root@localhost:5432/proj?sslmode=disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 25)
	viper.SetDefault("database.conn_max_lifetime", "5m")
	viper.SetDefault("database.conn_max_idle_time", "2m")
	viper.SetDefault("database.health_check_interval", "1m")

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.production", true)

	if err = viper.Unmarshal(&config); err != nil {
		return AppConfig{}, fmt.Errorf("unable to unmarshal config: %w", err)
	}
	return config, nil
}
