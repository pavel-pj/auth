package vault

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
)

type Config struct {
	Address    string
	Token      string
	RoleID     string
	SecretID   string
	Mount      string
	SecretPath string
	ServerName string

	JWTSecretPath         string
	JWTAccessSecretField  string
	JWTRefreshSecretField string

	PKIMount      string
	PKIRole       string
	PKICommonName string
	PKITTL        string
	RootCAFile    string
}

type TLS struct {
	ServerTLS  *tls.Config
	CAPool     *x509.CertPool
	ServerName string
}

type JWT struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
}

type Secret struct {
	TLS TLS
	JWT JWT
}

func LoadJWT(ctx context.Context, cfg Config) (JWT, error) {
	client, err := newClient(ctx, cfg)
	if err != nil {
		return JWT{}, fmt.Errorf("vault jwt: client: %w", err)
	}

	secret, err := client.KVv2(cfg.Mount).Get(ctx, cfg.JWTSecretPath)
	if err != nil {
		return JWT{}, fmt.Errorf("vault jwt: read %s/%s: %w", cfg.Mount, cfg.JWTSecretPath, err)
	}
	if secret == nil || secret.Data == nil {
		return JWT{}, fmt.Errorf("vault jwt: secret %s/%s not found", cfg.Mount, cfg.JWTSecretPath)
	}

	access, ok := secret.Data[cfg.JWTAccessSecretField].(string)
	if !ok || access == "" {
		return JWT{}, fmt.Errorf("vault jwt: field %q is empty", cfg.JWTAccessSecretField)
	}

	refresh, ok := secret.Data[cfg.JWTRefreshSecretField].(string)
	if !ok || refresh == "" {
		return JWT{}, fmt.Errorf("vault jwt: field %q is empty", cfg.JWTRefreshSecretField)
	}

	return JWT{
		AccessTokenSecret:  access,
		RefreshTokenSecret: refresh,
	}, nil
}

func newClient(ctx context.Context, cfg Config) (*api.Client, error) {
	vc := api.DefaultConfig()
	vc.Address = cfg.Address

	client, err := api.NewClient(vc)
	if err != nil {
		return nil, err
	}

	if cfg.RoleID != "" && cfg.SecretID != "" {
		auth, err := approle.NewAppRoleAuth(
			cfg.RoleID,
			&approle.SecretID{FromString: cfg.SecretID},
		)
		if err != nil {
			return nil, err
		}

		login, err := auth.Login(ctx, client)
		if err != nil {
			return nil, err
		}
		if login == nil || login.Auth == nil || login.Auth.ClientToken == "" {
			return nil, fmt.Errorf("vault: empty AppRole login response")
		}

		client.SetToken(login.Auth.ClientToken)
		return client, nil
	}

	if cfg.Token != "" {
		client.SetToken(cfg.Token)
		return client, nil
	}

	return nil, fmt.Errorf("vault: no credentials")
}
