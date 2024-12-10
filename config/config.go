package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppName       string
	AppEnv        string
	ServerPort    string
	DatabaseURI   string
	JWTSecret     string
	EmailHost     string
	EmailPort     int
	EmailUser     string
	EmailPassword string
}


func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")           // Name of the config file (without extension)
	viper.SetConfigType("env")              // If the config file is in env format
	viper.AddConfigPath(".")                // Look for the config file in the current directory
	viper.AddConfigPath("./config")         // Look for the config file in the config directory
	viper.AutomaticEnv()                    // Override config values with environment variables if they exist
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, reading from environment variables: %v", err)
	}

	config := &Config{
		AppName:       viper.GetString("APP_NAME"),
		AppEnv:        viper.GetString("APP_ENV"),
		ServerPort:    viper.GetString("SERVER_PORT"),
		DatabaseURI:   viper.GetString("DATABASE_URI"),
		JWTSecret:     viper.GetString("JWT_SECRET"),
		EmailHost:     viper.GetString("EMAIL_HOST"),
		EmailPort:     viper.GetInt("EMAIL_PORT"),
		EmailUser:     viper.GetString("EMAIL_USER"),
		EmailPassword: viper.GetString("EMAIL_PASSWORD"),
	}

	return config, nil
}

func InitConfig() *Config {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	return config
}

func (c *Config) GetServerPort() string {
	port := os.Getenv("PORT")
	if port != "" {
		return port
	}
	return c.ServerPort
}
