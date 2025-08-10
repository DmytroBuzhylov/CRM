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
	JWT      JWTConfig
	Minio    MinioConfig
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

type JWTConfig struct {
	JWTAccessSecret    string        `mapstructure:"jwt_access_secret"`
	JWTRefreshSecret   string        `mapstructure:"jwt_refresh_secret"`
	JWTAccessLifetime  time.Duration `mapstructure:"jwt_access_token_lifetime"`
	JWTRefreshLifetime time.Duration `mapstructure:"jwt_refresh_token_lifetime"`
}

type MinioConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyId     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	SSL             bool   `mapstructure:"ssl"`
}

func LoadConfig() (config AppConfig, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
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

	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 25)
	viper.SetDefault("database.conn_max_lifetime", "5m")
	viper.SetDefault("database.conn_max_idle_time", "2m")
	viper.SetDefault("database.health_check_interval", "1m")

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.production", true)

	viper.SetDefault("jwt.jwt_access_token_lifetime", "15m")
	viper.SetDefault("jwt.jwt_refresh_token_lifetime", "168h")

	if err = viper.Unmarshal(&config); err != nil {
		return AppConfig{}, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return config, nil
}
