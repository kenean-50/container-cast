package ssh

import (
	"golang.org/x/crypto/ssh"
)

type Ssh interface {
	Auth() (AuthMethod, error)
	Client() (Client error)
	Session() (Session, error)
}

type AuthMethod []ssh.AuthMethod

type Client interface {
	Connect() (*ssh.Client, error)
}

type Session interface {
	Run() error
	Upload() error
	Download() error
	Interactive() error
}
