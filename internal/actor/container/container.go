package container

import (
	"context"

	"github.com/rs/zerolog"
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
	image  imageOptions
	client clientOptions
	// ssh    sshOptions
	// ports  []portOptions
	// networks []networkOption
	// registry registryOption
	// configs  []configOption
	// secrets  sshOptions
}

type imageOptions struct {
	name string
	tag  string
}

// type clientOptions struct {
// 	client *client.Client
// }

// type sshOptions struct {
// 	client *ssh.Client
// }

type portOptions struct {
	containerPort string
	hostPort      string
}

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
