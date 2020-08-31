package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nickysemenza/gourd/parser"
	"github.com/nickysemenza/gourd/scraper"
	"github.com/spf13/cobra"
)

func main() {
	setupMisc()
	_ = rootCmd.Execute()
}

// nolint:gochecknoglobals
var rootCmd = &cobra.Command{
	Use:   "gourd",
	Short: "Go Universal Recipe Database",
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
				s := makeServer()

				r, err := scraper.FetchAndTransform(context.Background(), strings.Join(args, " "), s.GetResolver().Mutation().UpsertIngredient)
				if err != nil {
					return err
				}
				_, err = s.GetResolver().Mutation().UpdateRecipe(context.Background(), r)
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
