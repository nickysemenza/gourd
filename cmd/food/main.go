package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nickysemenza/food/parser"
	"github.com/nickysemenza/food/scraper"
	"github.com/spf13/cobra"
)

func main() {
	setupMisc()
	_ = rootCmd.Execute()
}

// nolint:gochecknoglobals
var rootCmd = &cobra.Command{
	Use:   "food",
	Short: "Food app",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "version",
			Short: "Print the version number of Hugo",
			Long:  `All software has versions. This is Hugo's`,
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
			},
		},
		&cobra.Command{
			Use:   "server",
			Short: "Run the server",
			Run: func(cmd *cobra.Command, args []string) {
				runServer()
			},
		},
		&cobra.Command{
			Use:   "ingredient-parse [ingredient]",
			Short: "parse an ingredient",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := context.Background()
				ingredient, err := parser.Parse(ctx, strings.Join(args, " "))
				if err != nil {
					return err
				}
				fmt.Println(ingredient.ToString())
				return nil
			},
		},
		&cobra.Command{
			Use:   "scrape [url]",
			Short: "scrape a recipe",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				_, err := scraper.GetIngredients(context.Background(), strings.Join(args, " "))
				return err
			},
		},
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
