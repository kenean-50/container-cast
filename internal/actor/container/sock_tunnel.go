package container

import (
	"fmt"
	"io"
	"net"

	"golang.org/x/crypto/ssh"
)

type SockTunnel struct {
	sshClient *ssh.Client
	listener  net.Listener
}

func NewSockTunnel(sshClient *ssh.Client) (*SockTunnel, error) {

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		sshClient.Close()
		return nil, fmt.Errorf("unable to start local listener: %v", err)
	}

	tunnel := &SockTunnel{
		sshClient: sshClient,
		listener:  listener,
	}

	go tunnel.forward()

	return tunnel, nil
}

func (t *SockTunnel) forward() {
	for {
		local, err := t.listener.Accept()
		if err != nil {
			return
		}

		go func() {
			remote, err := t.sshClient.Dial("unix", "/var/run/docker.sock")
			if err != nil {
				local.Close()
				return
			}

			go func() { _, _ = io.Copy(local, remote) }()
			go func() { _, _ = io.Copy(remote, local) }()
		}()
	}
}

func (t *SockTunnel) Close() {
	t.listener.Close()
	t.sshClient.Close()
}
