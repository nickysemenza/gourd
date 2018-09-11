package main

import (
	"os"

	"github.com/nickysemenza/food/backend/app"
	"github.com/nickysemenza/food/backend/app/config"
	"github.com/urfave/cli"
)

func main() {

	globalConfig := config.GetConfig()
	mainApp := app.NewApp(globalConfig)

	cliApp := cli.NewApp()
	cliApp.Version = "1.0.0"
	cliApp.Name = "Recipe Hub Backend API"
	cliApp.Authors = []cli.Author{
		{Name: "Nicky", Email: "nicky@nickysemenza.com"},
	}
	cliApp.Action = func(c *cli.Context) error {
		mainApp.RunServer()
		return nil
	}
	cliApp.Commands = []cli.Command{
		{
			Name:    "export",
			Aliases: []string{"e"},
			Usage:   "Export the recipes to JSON",
			Action: func(c *cli.Context) error {
				pwd, _ := os.Getwd()
				pwd += "/recipes/"
				mainApp.Export(pwd)
				return nil
			},
		},
		{
			Name:    "import",
			Aliases: []string{"i", "ingest"},
			Usage:   "Import a folder recipes from JSON",
			Action: func(c *cli.Context) error {
				pwd, _ := os.Getwd()
				pwd += "/recipes/"
				mainApp.Import(pwd)
				return nil
			},
		},
	}
	cliApp.Run(os.Args)

}
