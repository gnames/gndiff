package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DataPath string
}

type Option func(*Config)

func OptDataPath(s string) Option {
	return func(c *Config) {
		c.DataPath = s
	}
}

func New(opts ...Option) Config {
	cache, _ := os.UserCacheDir()
	res := Config{
		DataPath: filepath.Join(cache, "gndiff"),
	}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
