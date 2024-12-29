package main

import (
	"os"
	"runtime/pprof"

	"github.com/kenean-50/container-cast/internal/actor/cli"
	"github.com/kenean-50/container-cast/internal/actor/config"
	"github.com/kenean-50/container-cast/internal/domain/deploy"
	sConfig "github.com/kenean-50/container-cast/internal/util/config"
	"github.com/kenean-50/container-cast/internal/util/logger"
	"github.com/rs/zerolog/log"
)

func main() {

	f, err := os.Create("cpu.pprof")
	if err != nil {
		panic(err)
	}
	// set profile sample
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Init System Config
	//
	c := sConfig.NewConfig()
	err = c.LoadConfigFile(".", "env", ".env")
	if err != nil {
		log.
			Error().
			Msg(err.Error())
	}

	// Init Logger
	//
	logger.InitLogger(c.GetString("LOG_LEVEL"), c.GetString("ENVIRONMENT"))

	// load config
	//
	config := config.NewConfig()
	cfg, err := config.LoadFromYAML(c.GetString("MANIFEST_PATH"), c.GetString("MANIFEST_FILE"))

	if err != nil {
		log.
			Error().
			Msg(err.Error())
	}

	// Initialize modules
	//
	deployModule := deploy.NewDeployModule(cfg)

	// Initialize Cobra Cli
	//
	cm := cli.NewCobraCli(
		c.GetString("APP_NAME"),
		deployModule,
	)

	if err := cm.Execute(); err != nil {
		log.
			Error().
			Msg(err.Error())
	}
}
