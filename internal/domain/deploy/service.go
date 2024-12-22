package deploy

import (
	"context"
	"fmt"
	"sync"

	"github.com/kenean-50/vm-container-manager/internal/actor/manifest"

	"github.com/docker/docker/client"
	"github.com/kenean-50/vm-container-manager/internal/actor/container"
	"github.com/kenean-50/vm-container-manager/internal/actor/ssh"
	"github.com/rs/zerolog/log"
)

func (d *DeployConfig) Apply() {
	var wg sync.WaitGroup

	for _, server := range d.manifest.Values.Servers {
		wg.Add(1)

		go func(server manifest.Server) {
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

			for name, service := range d.manifest.Values.Services {
				containerWg.Add(1)

				go func(name string, service manifest.Service) {
					defer containerWg.Done()
					runContainer(dClient, service.Image, name)
				}(name, service)
			}

			containerWg.Wait()
		}(server)
	}

	wg.Wait()
}

func runContainer(dClient *client.Client, imageName, imageTag string) {
	con := container.NewContainer(
		context.Background(),
		container.WithClient(dClient),
		container.WithImage(imageName, imageTag),
	)
	out := con.Run()
	fmt.Println(out)
}
