package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	Host     string //  include :port if local
	DBPath   string // file where sqlite db will be
	Env      string // dev | staging | prod
	RedisURL string // redis url e.g. localhost:6379
	Port     int    // for http server PORT config
}

func (c *Config) validate() error {
	if c.Host == "" {
		return errors.New("HOST is required")
	}
	if c.DBPath == "" {
		return errors.New("DB_PATH is required")
	}
	if c.Port == 0 {
		return errors.New("PORT is required")
	}
	return nil
}

func Load() *Config {
	c := &Config{}

	c.Host = os.Getenv("HOST")
	c.DBPath = os.Getenv("DB_PATH")
	c.Env = os.Getenv("APP_ENV")
	c.RedisURL = os.Getenv("REDIS_URL")
	sp := os.Getenv("PORT")

	port, err := strconv.Atoi(sp)
	if err != nil {
		panic(fmt.Sprintf("invalid PORT value %q: %v", sp, err))
	}

	c.Port = port

	if err = c.validate(); err != nil {
		panic(err.Error())
	}

	slog.Info("config loaded", "host", c.Host, "db_path", c.DBPath, "env", c.Env, "redis_url", c.RedisURL, "port", c.Port)

	return c
}
