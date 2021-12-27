package config

import "github.com/gnames/gnfmt"

type Config struct {
	Format gnfmt.Format
}

type Option func(*Config)

func OptFormat(f gnfmt.Format) Option {
	return func(cfg *Config) {
		cfg.Format = f
	}
}

func New(opts ...Option) Config {
	res := Config{
		Format: gnfmt.CSV,
	}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
