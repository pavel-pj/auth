package config

type PGConfig interface {
	DSN() string
}

type GRPCConfig interface {
	Address() string
}

type HTTPConfig interface {
	Address() string
}

type SwaggerConfig interface {
	Address() string
}

type VaultConfig interface {
	Address() string
	Token() string
	RoleID() string
	SecretID() string
	KVMount() string
	SecretPath() string
	TLSServerName() string

	JWTSecretPath() string
	JWTAccessSecretField() string
	JWTRefreshSecretField() string
}
