package container

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/container"
	dcontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog/log"
)

func NewContainer(ctx context.Context, opts ...ContainerOptions) *containerOptions {
	c := &containerOptions{
		ctx: ctx,
		logger: log.
			With().
			Str("actor", "container-docker").
			Logger(),
	}
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

func (c clientOptions) apply(opts *containerOptions) {
	opts.client = c
}

func (i imageOptions) apply(opts *containerOptions) {
	opts.image = i
}

func (p portOptions) apply(opts *containerOptions) {
	opts.ports = p
}

func WithImage(name, tag string) ContainerOptions {
	return imageOptions{
		name: name,
		tag:  tag,
	}
}

func WithClient(client *client.Client) ContainerOptions {
	return clientOptions{
		docker: client,
	}
}

func WithPort(ports []string) ContainerOptions {
	var pm []PortMapping

	for _, p := range ports {
		parts := strings.Split(p, ":")
		if len(parts) != 2 {
			panic("invalid port format, expected 'host:container'")
		}
		pm = append(pm, PortMapping{
			containerPort: parts[1],
			hostPort:      parts[0],
		})
	}

	return portOptions(pm)
}

func (c containerOptions) Run() string {
	config := &dcontainer.Config{
		Image: c.image.name,
	}

	portBindings := nat.PortMap{}
	for _, pm := range c.ports {
		portBindings[nat.Port(pm.containerPort+"/tcp")] = []nat.PortBinding{
			{
				HostIP:   "",
				HostPort: pm.hostPort,
			},
		}
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
	}

	resp, err := c.client.docker.ContainerCreate(c.ctx, config, hostConfig, nil, nil, "")
	err = c.client.docker.ContainerStart(c.ctx, resp.ID, dcontainer.StartOptions{})

	if err != nil {
		c.logger.
			Fatal().
			Str("status", "failed to run container").
			Str("reason", ""+err.Error()).
			Send()
	}
	return resp.ID
}
