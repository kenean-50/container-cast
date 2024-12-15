package ssh

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type clientConfig struct {
	*ssh.Client
	host           string
	user           string
	port           int
	auth           AuthMethod
	timeout        time.Duration
	callback       ssh.HostKeyCallback
	bannerCallback ssh.BannerCallback
	logger         zerolog.Logger
}

func NewClient(host, user string, port int, auth AuthMethod) *clientConfig {
	callback, err := defaultKnownHosts()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to get private key auth")
	}

	return &clientConfig{
		auth:     auth,
		host:     host,
		port:     port,
		user:     user,
		callback: callback,
		logger: log.
			With().
			Str("actor", "ssh").
			Logger(),
	}
}

func (r *clientConfig) Connect() *ssh.Client {

	conn, err := ssh.Dial("tcp", net.JoinHostPort(r.host, fmt.Sprint(r.port)), &ssh.ClientConfig{
		User:            r.user,
		Auth:            r.auth,
		Timeout:         r.timeout,
		HostKeyCallback: r.callback,
		BannerCallback:  r.bannerCallback,
	})

	if err != nil {
		r.logger.
			Fatal().
			Str("status", "failed to connect").
			Str("reason", ""+err.Error()).
			Send()
	}

	return conn
}

// todo: add checks for new hosts and save them
func defaultKnownHosts() (ssh.HostKeyCallback, error) {

	path, err := defaultKnownHostsPath()
	if err != nil {
		return nil, err
	}

	return knownHosts(path)
}

func knownHosts(file string) (ssh.HostKeyCallback, error) {
	return knownhosts.New(file)
}

func defaultKnownHostsPath() (string, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.ssh/known_hosts", home), err
}
