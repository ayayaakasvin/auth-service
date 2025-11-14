package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	configPathEnvKey 	= "CONFIG_PATH"
	postgresURLEnvKey 	= "POSTGRES_URL"
	valkeyURLEnvKey 	= "VALKEY_URL"
	jwtSecretKey	 	= "JWT_SECRET"
)

// Config represents the configuration structure
type Config struct {
    // Add your configuration fields here
    // Example:
    // Port int `yaml:"port"`
	HTTPServer						`yaml:"http-server"`

	PostgreSQLConfig

	ValkeyConfig

	JWTSecret						`yaml:"jwt"`
}

type HTTPServer struct {
	Address     	string        	`yaml:"address"					env-default:"localhost:8080"`
	Timeout     	time.Duration 	`yaml:"timeout" 															env-required:"true"`
	IdleTimeout 	time.Duration 	`yaml:"idle_timeout" 														env-required:"true"`
}

type PostgreSQLConfig struct {
	URL 			string
}

type ValkeyConfig struct {
	URL				string
}

type JWTSecret struct {
	Secret 			string			
	AccessTokenTTL	time.Duration	`yaml:"access_ttl"				env-default:"10m"							env-required:"true"`
	RefreshTokenTTL	time.Duration	`yaml:"refresh_ttl"				env-default:"168h"							env-required:"true"`
}

// MustLoadConfig loads the configuration from the specified path
func MustLoadConfig () *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, falling back to local variables")
	}

	configPath := os.Getenv(configPathEnvKey)
	if configPath == "" {
		log.Fatalf("%s is not set up", configPathEnvKey)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist: %s", configPath, err.Error())
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config file: %s", err.Error())
	}

	postgresURL := os.Getenv(postgresURLEnvKey)
	valkeyURL := os.Getenv(valkeyURLEnvKey)
	jwtSecret := os.Getenv(jwtSecretKey)

	cfg.PostgreSQLConfig.URL = postgresURL
	cfg.ValkeyConfig.URL = valkeyURL
	cfg.JWTSecret.Secret = jwtSecret

	return &cfg
}