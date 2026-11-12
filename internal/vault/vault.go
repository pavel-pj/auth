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

func Load(ctx context.Context, cfg Config) (*Secret, error) {
	client, err := newClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("vault client: %w", err)
	}

	secret, err := client.KVv2(cfg.Mount).Get(ctx, cfg.SecretPath)
	if err != nil {
		return nil, fmt.Errorf("vault read %s/%s: %w", cfg.Mount, cfg.SecretPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("vault: secret %s/%s not found", cfg.Mount, cfg.SecretPath)
	}

	certPEM, _ := secret.Data["cert"].(string)
	keyPEM, _ := secret.Data["key"].(string)
	caPEM, _ := secret.Data["ca"].(string)
	if certPEM == "" || keyPEM == "" || caPEM == "" {
		return nil, fmt.Errorf("vault: fields cert/key/ca are required")
	}

	pair, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		return nil, fmt.Errorf("vault: parse key pair: %w", err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(caPEM)) {
		return nil, fmt.Errorf("vault: parse ca")
	}

	jwtSecret, err := client.KVv2(cfg.Mount).Get(ctx, cfg.JWTSecretPath)
	if err != nil {
		return nil, fmt.Errorf("vault read %s/%s: %w", cfg.Mount, cfg.JWTSecretPath, err)
	}
	if jwtSecret == nil || jwtSecret.Data == nil {
		return nil, fmt.Errorf("vault: secret %s/%s not found", cfg.Mount, cfg.JWTSecretPath)
	}

	accessSecret, ok := jwtSecret.Data[cfg.JWTAccessSecretField].(string)
	if !ok || accessSecret == "" {
		return nil, fmt.Errorf("vault: field %q is empty or not a string", cfg.JWTAccessSecretField)
	}

	refreshSecret, ok := jwtSecret.Data[cfg.JWTRefreshSecretField].(string)
	if !ok || refreshSecret == "" {
		return nil, fmt.Errorf("vault: field %q is empty or not a string", cfg.JWTRefreshSecretField)
	}

	return &Secret{
		TLS: TLS{
			ServerTLS: &tls.Config{
				Certificates: []tls.Certificate{pair},
				MinVersion:   tls.VersionTLS12,
			},
			CAPool:     pool,
			ServerName: cfg.ServerName,
		},
		JWT: JWT{
			AccessTokenSecret:  accessSecret,
			RefreshTokenSecret: refreshSecret,
		},
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
