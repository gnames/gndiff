package config

type Config struct {
}

type Option func(*Config)

func New(opts ...Option) Config {
	res := Config{}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
