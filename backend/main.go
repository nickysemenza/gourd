package main

import (
	"github.com/nickysemenza/food/backend/app"
	"github.com/nickysemenza/food/backend/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":3000")
}
