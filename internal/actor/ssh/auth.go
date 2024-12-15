package ssh

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

type authConfig struct {
	password   string
	privateKey string
	logger     zerolog.Logger
}

func NewAuthPassword(password string) *authConfig {
	return &authConfig{
		password: password,
		logger: log.
			With().
			Str("actor", "ssh").
			Logger(),
	}
}

func NewAuthPrivateKey(privateKey string) *authConfig {
	return &authConfig{
		privateKey: privateKey,
		logger: log.
			With().
			Str("actor", "ssh").
			Logger(),
	}
}

func (r *authConfig) Password() AuthMethod {
	return AuthMethod{
		ssh.Password(r.password),
	}
}

func (r *authConfig) PrivateKey() AuthMethod {
	signer, err := sign(r.privateKey)

	if err != nil {
		r.logger.
			Fatal().
			Str("status", "failed to get private key auth").
			Str("reason", ""+err.Error()).
			Send()
	}

	return AuthMethod{
		ssh.PublicKeys(signer),
	}
}

func sign(PrivateKey string) (ssh.Signer, error) {
	var (
		err    error
		signer ssh.Signer
	)

	privateKey, err := getPrivateKey(PrivateKey)
	if err != nil {
		return nil, err
	} else {
		signer, err = ssh.ParsePrivateKey(privateKey)
	}
	return signer, err
}

func getPrivateKey(privateKey string) ([]byte, error) {
	key, err := os.ReadFile(privateKey)
	// todo: get key from raw string
	if err != nil {
		return nil, err
	} else {
		return key, nil
	}
}
