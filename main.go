package main

import (
	"os"

	"github.com/SpeedVan/gitlab-2-jira-hook/controller"
	"github.com/SpeedVan/go-common/app/web"
	"github.com/SpeedVan/go-common/config/env"
	"github.com/SpeedVan/go-common/log"
)

func main() {
	if cfg, err := env.LoadAllWithoutPrefix("G2J_"); err == nil {
		logger := log.NewCommon(log.Debug)
		app := web.New(cfg, logger)
		app.HandleController(&controller.HookController{})
		if err := app.Run(log.Debug); err != nil {
			os.Exit(1)
		}
	}
}
