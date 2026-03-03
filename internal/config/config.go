package config

import "os"

type Config struct {
	Host   string //  include :port if local
	DBPath string // file where sqlite db will be
	Env    string // dev | staging | prod
}

func Load() *Config {
	c := &Config{}

	c.Host = os.Getenv("HOST")
	c.DBPath = os.Getenv("DB_PATH")
	c.Env = os.Getenv("APP_ENV")

	return c
}
