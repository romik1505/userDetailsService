package config

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type ConfigFile struct {
	AppLevel    string
	DBConfig    DBConfig
	KafkaConfig KafkaConfig
	Cache       CacheConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Driver   string
	User     string
	Password string
	Name     string
}

func (d DBConfig) ConnectionString() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		d.Driver,
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
	)
}

type CacheConfig struct {
	Host     string
	Port     string
	Password string
}

type KafkaConfig struct {
	Brokers string
	GroupID string
}

var (
	Config ConfigFile
	once   sync.Once
)

func newConfig() ConfigFile {

	c := ConfigFile{
		AppLevel: os.Getenv("APP_LEVEL"),
		DBConfig: DBConfig{
			Host:     mustGetEnv("DB_HOST"),
			Port:     mustGetEnv("DB_PORT"),
			Driver:   mustGetEnv("DB_DRIVER"),
			User:     mustGetEnv("DB_USER"),
			Password: mustGetEnv("DB_PASSWORD"),
			Name:     mustGetEnv("DB_NAME"),
		},
		KafkaConfig: KafkaConfig{
			Brokers: mustGetEnv("KAFKA_BROKERS"),
			GroupID: mustGetEnv("KAFKA_CONSUMER_GROUP"),
		},
		Cache: CacheConfig{
			Host:     mustGetEnv("REDIS_HOST"),
			Port:     mustGetEnv("REDIS_PORT"),
			Password: mustGetEnv("REDIS_PASSWORD"),
		},
	}

	return c
}

func mustGetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		log.Printf("config [%s]=%s\n", key, value)
		return value
	}
	log.Fatalf("config key %s not set", key)

	return ""
}

func init() {
	once.Do(func() {
		Config = newConfig()
	})
}
