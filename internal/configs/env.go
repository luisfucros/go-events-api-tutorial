package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	DB     DatabaseConfig
	JWT    JWTConfig
	Redis  RedisConfig
}

type ServerConfig struct {
	PublicHost string
	Port       int64
}

type DatabaseConfig struct {
	User     string
	Password string
	Address  string
	Name     string
}

type JWTConfig struct {
	Secret              string
	ExpirationInSeconds int64
}

type RedisConfig struct {
	Addr    string
	Password string
	DB      int
	Enabled bool
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Server: ServerConfig{
			PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
			Port:       getEnvAsInt("PORT", 8080),
		},
		DB: DatabaseConfig{
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "events_password"),
			Address: fmt.Sprintf(
				"%s:%s",
				getEnv("DB_HOST", "127.0.0.1"),
				getEnv("DB_PORT", "3306"),
			),
			Name: getEnv("DB_NAME", "events"),
		},
		JWT: JWTConfig{
			Secret:              getEnv("JWT_SECRET", "super-secret"),
			ExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PW", ""),
			DB:       int(getEnvAsInt("REDIS_DB", 0)),
			Enabled:  getEnvAsBool("REDIS_ENABLED", false),
		},
	}
}

// Gets the env by key or fallbacks
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return boolVal
}