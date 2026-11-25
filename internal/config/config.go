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

	PKIMount() string
	PKIRole() string
	PKICommonName() string
	PKITTL() string
	RootCAFile() string
}

type LoggerConfig interface {
	Level() string
	Filename() string
	MaxSize() int
	MaxBackups() int
	MaxAge() int
}

type PrometeusConfig interface {
	Address() string
	Namespace() string
	AppName() string
}
