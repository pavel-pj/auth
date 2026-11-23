package vault

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadPKI(ctx context.Context, cfg Config) (TLS, error) {
	if cfg.PKIRole == "" {
		return TLS{}, fmt.Errorf("vault pki: role is empty")
	}

	if cfg.PKICommonName == "" {
		return TLS{}, fmt.Errorf("vault pki: common name is empty")
	}

	client, err := newClient(ctx, cfg)
	if err != nil {
		return TLS{}, fmt.Errorf("vault pki: client: %w", err)
	}

	mount := cfg.PKIMount
	if mount == "" {
		mount = "pki"
	}

	data := map[string]interface{}{
		"common_name": cfg.PKICommonName,
	}

	if cfg.PKITTL != "" {
		data["ttl"] = cfg.PKITTL
	}

	path := fmt.Sprintf("%s/issue/%s", mount, cfg.PKIRole)

	secret, err := client.Logical().Write(path, data)
	if err != nil {
		return TLS{}, fmt.Errorf("vault pki: issue %s: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return TLS{}, fmt.Errorf("vault pki: empty issue response")
	}

	certPEM, err := secretString(secret.Data, "certificate")
	if err != nil {
		return TLS{}, err
	}

	keyPEM, err := secretString(secret.Data, "private_key")
	if err != nil {
		return TLS{}, err
	}

	pair, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		return TLS{}, fmt.Errorf("vault pki: parse key pair: %w", err)
	}

	chain, err := secretChain(secret.Data)
	if err != nil {
		return TLS{}, err
	}

	if err := appendCertificateChain(&pair, chain); err != nil {
		return TLS{}, err
	}

	rootPEM, err := os.ReadFile(cfg.RootCAFile)
	if err != nil {
		return TLS{}, fmt.Errorf("vault pki: read root ca %q: %w", cfg.RootCAFile, err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(rootPEM) {
		return TLS{}, fmt.Errorf("vault pki: parse root ca %q", cfg.RootCAFile)
	}

	return TLS{
		ServerTLS: &tls.Config{
			Certificates: []tls.Certificate{pair},
			MinVersion:   tls.VersionTLS12,
		},
		CAPool:     pool,
		ServerName: cfg.ServerName,
	}, nil
}

func secretString(data map[string]interface{}, key string) (string, error) {
	value, ok := data[key]
	if !ok {
		return "", fmt.Errorf("vault pki: field %q is missing", key)
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("vault pki: field %q is not a string", key)
	}

	if str == "" {
		return "", fmt.Errorf("vault pki: field %q is empty", key)
	}

	return str, nil
}

func secretChain(data map[string]interface{}) ([]string, error) {
	switch chain := data["ca_chain"].(type) {
	case []interface{}:
		result := make([]string, 0, len(chain))

		for i, item := range chain {
			str, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("vault pki: ca_chain[%d] is not a string", i)
			}
			result = append(result, str)
		}

		if len(result) > 0 {
			return result, nil
		}

	case []string:
		if len(chain) > 0 {
			return chain, nil
		}

	case nil:
		// fallback to issuing_ca
	default:
		return nil, fmt.Errorf("vault pki: ca_chain has unsupported type %T", chain)
	}

	issuingCA, _ := data["issuing_ca"].(string)
	if issuingCA == "" {
		return nil, fmt.Errorf("vault pki: ca_chain and issuing_ca are empty")
	}

	return []string{issuingCA}, nil
}

func appendCertificateChain(pair *tls.Certificate, chain []string) error {
	for i, certPEM := range chain {
		block, _ := pem.Decode([]byte(certPEM))
		if block == nil {
			return fmt.Errorf("vault pki: ca_chain[%d]: failed to decode PEM", i)
		}

		if block.Type != "CERTIFICATE" {
			return fmt.Errorf("vault pki: ca_chain[%d]: unexpected PEM type %q", i, block.Type)
		}

		pair.Certificate = append(pair.Certificate, block.Bytes)
	}

	return nil
}
