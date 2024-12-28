package container

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/ssh"
)

type Container interface {
	Run() string
	Build()
	Start()
	Stop()
	Tail()
	DeleteContainer()
}

type ContainerOptions interface {
	apply(*containerOptions)
}

type containerOptions struct {
	ctx    context.Context
	logger zerolog.Logger
	client clientOptions
	image  imageOptions
	ports  portOptions
	// networks []networkOption
	// registry registryOption
	// configs  []configOption
	// secrets  sshOptions
}

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

type imageOptions struct {
	name string
	tag  string
}

type PortMapping struct {
	containerPort string
	hostPort      string
}

type portOptions []PortMapping

// type clientOptions struct {
// 	client *client.Client
// }

// type networkOption struct {
// 	network string
// }

// type registryOption struct {
// 	host     string
// 	passowrd string
// }

// type secretOption struct {
// 	key   string
// 	value string
// }

// type configOption struct {
// 	key   string
// 	value string
// }

// type RunOptions interface {
// 	apply(*runOptions)
// }

// type runOptions struct {
// 	networks []networkOption
// 	ports    []portOptions
// 	registry registryOption
// 	configs  []configOption
// 	secrets  sshOptions
// }
