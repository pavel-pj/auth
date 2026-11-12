package env

import (
	"log"

	"github.com/joho/godotenv"
)

type VaultConfig struct {
	address    string
	token      string
	roleID     string
	secretID   string
	kvMount    string
	secretPath string
	serverName string

	jwtSecretPath         string
	jwtAccessSecretField  string
	jwtRefreshSecretField string
}

func NewVaultConfig() (VaultConfig, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := VaultConfig{}

	cfg.address = getEnv("VAULT_ADDR", "http://127.0.0.1:8200")
	cfg.token = getEnv("VAULT_TOKEN", "")
	cfg.roleID = getEnv("VAULT_ROLE_ID", "")
	cfg.secretID = getEnv("VAULT_SECRET_ID", "")
	cfg.kvMount = getEnv("VAULT_KV_MOUNT", "secret")
	cfg.secretPath = getEnv("VAULT_SECRET_PATH", "auth/tls")
	cfg.serverName = getEnv("VAULT_TLS_SERVER_NAME", "localhost")

	cfg.jwtSecretPath = getEnv("VAULT_JWT_SECRET_PATH", "auth/jwt")
	cfg.jwtAccessSecretField = getEnv("VAULT_JWT_ACCESS_SECRET_FIELD", "AccessTokenSecret")
	cfg.jwtRefreshSecretField = getEnv("VAULT_JWT_REFRESH_SECRET_FIELD", "RefreshTokenSecret")

	return cfg, nil
}

func (s VaultConfig) Address() string {
	return s.address
}

func (s VaultConfig) Token() string {
	return s.token
}

func (s VaultConfig) RoleID() string {
	return s.roleID
}

func (s VaultConfig) SecretID() string {
	return s.secretID
}

func (s VaultConfig) KVMount() string {
	return s.kvMount
}

func (s VaultConfig) SecretPath() string {
	return s.secretPath
}

func (s VaultConfig) TLSServerName() string {
	return s.serverName
}

func (s VaultConfig) JWTAccessSecretField() string {
	return s.jwtAccessSecretField
}

func (s VaultConfig) JWTRefreshSecretField() string {
	return s.jwtRefreshSecretField
}

func (s VaultConfig) JWTSecretPath() string {
	return s.jwtSecretPath
}
