package config

type PGConfig interface {
	DSN() string
}

type GRPCConfig interface {
	Address() string
}
