package ssh

import (
	"fmt"
	"net"
	"os"
	"time"

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
}

func NewClient(host, user string, port int, auth AuthMethod) (*clientConfig, error) {
	callback, err := defaultKnownHosts()

	if err != nil {
		return nil, err
	}

	return &clientConfig{
		auth:     auth,
		host:     host,
		port:     port,
		user:     user,
		callback: callback,
	}, nil
}

func (r *clientConfig) Connect() (*ssh.Client, error) {
	return ssh.Dial("tcp", net.JoinHostPort(r.host, fmt.Sprint(r.port)), &ssh.ClientConfig{
		User:            r.user,
		Auth:            r.auth,
		Timeout:         r.timeout,
		HostKeyCallback: r.callback,
		BannerCallback:  r.bannerCallback,
	})
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
