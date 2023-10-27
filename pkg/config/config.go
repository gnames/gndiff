package config

import "github.com/gnames/gnfmt"

// Config provides configuration parameters.
type Config struct {

	// Format sets the output format for CLI and Web interfaces.
	// There are 4 formats available: 'CSV', 'TSV', 'CompactJSON' and
	// 'PrettyJSON'.
	Format gnfmt.Format

	// WithUninomialFuzzyMatch is true when it is allowed to use fuzzy match for
	// uninomial names.
	WithUninomialFuzzyMatch bool
}

type Option func(*Config)

func OptFormat(f gnfmt.Format) Option {
	return func(cfg *Config) {
		cfg.Format = f
	}
}

func OptWithUninomialFuzzyMatch(b bool) Option {
	return func(cfg *Config) {
		cfg.WithUninomialFuzzyMatch = b
	}
}

func New(opts ...Option) Config {
	res := Config{
		Format:                  gnfmt.CSV,
		WithUninomialFuzzyMatch: false,
	}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
