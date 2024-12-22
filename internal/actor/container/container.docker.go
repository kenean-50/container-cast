package container

import (
	"context"

	dcontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

func (c containerOptions) Run() string {
	config := &dcontainer.Config{
		Image: c.image.name,
	}

	resp, err := c.client.docker.ContainerCreate(c.ctx, config, nil, nil, nil, "")
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
