package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	MySQL  MySQLConfig
	Redis  RedisConfig
}

type ServerConfig struct {
	Port int
}

type MySQLConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Charset  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	Database int
}

var AppConfig *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	log.Println("Config loaded successfully")
}

func (m *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		m.Username, m.Password, m.Host, m.Port, m.Database, m.Charset)
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
