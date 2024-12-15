package ssh

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

type AuthOption interface {
	apply(*authOptions)
}

type authOptions struct {
	password   passwordOption
	privateKey privateKeyOption
	logger     zerolog.Logger
}

type passwordOption struct {
	password string
}

type privateKeyOption struct {
	privateKey string
}

func (p privateKeyOption) apply(opts *authOptions) {
	opts.privateKey = p
}

func (p passwordOption) apply(opts *authOptions) {
	opts.password = p
}

func NewAuth(opts ...AuthOption) *authOptions {
	a := &authOptions{
		logger: log.
			With().
			Str("actor", "ssh").
			Logger(),
	}
	for _, opt := range opts {
		opt.apply(a)
	}
	return a
}

func WithPassword(password string) AuthOption {
	return passwordOption{
		password: password,
	}
}

func WithPrivateKey(privateKey string) AuthOption {
	return privateKeyOption{
		privateKey: privateKey,
	}
}

func (r *authOptions) AuthMethod() AuthMethod {
	var methods []ssh.AuthMethod
	if r.password.password != "" {
		methods = append(methods, ssh.Password(r.password.password))
	}

	if r.privateKey.privateKey != "" {
		signer, err := sign(r.privateKey.privateKey)
		if err != nil {
			r.logger.
				Fatal().
				Str("status", "failed to get private key auth").
				Str("reason", ""+err.Error()).
				Send()
		}
		methods = append(methods, ssh.PublicKeys(signer))
	}

	return methods
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
