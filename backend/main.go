package main

import (
	"fmt"
	"github.com/nickysemenza/food/backend/app"
	"github.com/urfave/cli"
	"os"
)

func main() {

	globalConfig := app.GetConfig()
	mainApp := &app.App{}
	appEnv := mainApp.Initialize(globalConfig)

	cliApp := cli.NewApp()
	cliApp.Name = "Food Backend"
	cliApp.Authors = []cli.Author{
		{Name: "Nicky", Email: "nicky@nickysemenza.com"},
	}
	cliApp.Action = func(c *cli.Context) error {
		fmt.Println("Hello! Running API server on port", globalConfig.Port)
		mainApp.Run(fmt.Sprintf(":%d", globalConfig.Port))
		return nil
	}
	cliApp.Commands = []cli.Command{
		{
			Name:    "export",
			Aliases: []string{"e"},
			Usage:   "export the recipes to JSON",
			Action: func(c *cli.Context) error {
				pwd, _ := os.Getwd()
				pwd += "/recipes/"
				app.Utils{appEnv}.Export(pwd)
				return nil
			},
		},
	}

	cliApp.Run(os.Args)

}
