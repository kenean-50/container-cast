package deploy

import (
	"context"
	"fmt"
	"sync"

	"github.com/docker/docker/client"
	"github.com/kenean-50/container-cast/internal/actor/config"
	"github.com/kenean-50/container-cast/internal/actor/container"
	"github.com/kenean-50/container-cast/internal/actor/ssh"
	"github.com/rs/zerolog/log"
)

func (d *DeployConfig) Apply() {
	var wg sync.WaitGroup

	for _, server := range d.config.Values.Servers {
		wg.Add(1)

		go func(server config.Server) {
			defer wg.Done()

			auth := ssh.NewAuth(ssh.WithPrivateKey(server.PrivateKeyPath)).AuthMethod()
			sClient := ssh.NewClient(
				server.Host,
				server.User,
				server.SSHPort,
				auth,
			).Connect()
			defer sClient.Close()

			dClient, err := container.NewClient(
				container.WithSsh(sClient),
			).Connect()
			if err != nil {
				log.Error().Msg(err.Error())
				return
			}
			defer dClient.Close()

			var containerWg sync.WaitGroup

			for name, service := range d.config.Values.Services {
				containerWg.Add(1)

				go func(name string, service config.Service) {
					defer containerWg.Done()

					runContainer(dClient, service.Image, name, service.Ports)
				}(name, service)
			}

			containerWg.Wait()
		}(server)
	}

	wg.Wait()
}

func runContainer(dClient *client.Client, imageName, imageTag string, ports []string) {
	con := container.NewContainer(
		context.Background(),
		container.WithDockerClient(dClient),
		container.WithImage(imageName, imageTag),
		container.WithPort(ports),
	)

	con.PullImage(imageName)
	out := con.Run()
	fmt.Println(out)
}
