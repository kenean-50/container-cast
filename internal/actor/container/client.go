package container

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

type ClientOptions interface {
	apply(*clientOptions)
}

type clientOptions struct {
	ssh    sshOptions
	docker *client.Client
	logger zerolog.Logger
}

type sshOptions struct {
	client *ssh.Client
}

func NewClient(opts ...ClientOptions) *clientOptions {
	c := &clientOptions{
		logger: log.
			With().
			Str("actor", "container-client").
			Logger(),
	}
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

func (s sshOptions) apply(opts *clientOptions) {
	opts.ssh = s
}

func WithSsh(client *ssh.Client) ClientOptions {
	return sshOptions{
		client: client,
	}
}

func (s *clientOptions) Connect() (*client.Client, error) {
	var host string
	if s.ssh.client != nil {
		tunnel, err := NewSockTunnel(s.ssh.client)
		if err != nil {
			s.logger.
				Fatal().
				Str("status", "failed to connect").
				Str("reason", ""+err.Error()).
				Send()
		}
		// todo: close tunnel properly
		// defer tunnel.Close()
		host = fmt.Sprintf("tcp://%s", tunnel.listener.Addr().String())
	}

	if host == "" {
		s.logger.
			Fatal().
			Str("status", "failed to connect").
			Str("reason", "no host provided").
			Send()
	}

	return client.NewClientWithOpts(
		client.WithHost(host),
		client.WithAPIVersionNegotiation(),
	)
}
