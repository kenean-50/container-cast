package ssh

import (
	"os"

	"golang.org/x/crypto/ssh"
)

type authConfig struct {
	password   string
	privateKey string
}

func NewAuthPassword(password string) *authConfig {
	return &authConfig{
		password: password,
	}
}

func NewAuthPrivateKey(privateKey string) *authConfig {
	return &authConfig{
		privateKey: privateKey,
	}
}

func (r *authConfig) Password() AuthMethod {
	return AuthMethod{
		ssh.Password(r.password),
	}
}

func (r *authConfig) PrivateKey() (AuthMethod, error) {
	signer, err := sign(r.privateKey)

	if err != nil {
		return nil, err
	}

	return AuthMethod{
		ssh.PublicKeys(signer),
	}, nil
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
