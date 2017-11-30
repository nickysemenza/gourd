package main

import (
	"github.com/nickysemenza/food/backend/app"
	"github.com/nickysemenza/food/backend/config"
)

func main() {
	globalConfig := config.GetConfig()

	mainApp := &app.App{}
	mainApp.Initialize(globalConfig)
	mainApp.Run(":4000")
}
