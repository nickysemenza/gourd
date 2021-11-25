package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/davecgh/go-spew/spew"
	"github.com/nickysemenza/gourd/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	if err := setupMisc(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Info("cleaning up tracer")
		err := otel.GetTracerProvider().(*tracesdk.TracerProvider).ForceFlush(context.Background())
		if err != nil {
			log.Error(err)
		}
	}()
	err := rootCmd.Execute()
	if err != nil {
		log.Error(err)
	}

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
			Use:   "tmp",
			Short: "misc",
			RunE: func(cmd *cobra.Command, args []string) error {
				s, err := makeServer()
				if err != nil {
					return err
				}
				res, err := s.APIManager.Notion.GetAll(context.Background())
				if err != nil {
					return err
				}
				spew.Dump(res)
				return nil
			},
		},
		&cobra.Command{
			Use:   "sync",
			Short: "Run the server",
			RunE: func(cmd *cobra.Command, args []string) error {
				s, err := makeServer()
				if err != nil {
					return err
				}
				ctx := context.Background()

				err = s.APIManager.DoSync(ctx)
				if err != nil {
					return err
				}

				return nil
			},
		},

		&cobra.Command{
			Use:   "scrape [url]",
			Short: "scrape a recipe",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				s, err := makeServer()
				if err != nil {
					return err
				}
				ctx := context.Background()

				r, err := s.APIManager.FetchAndTransform(ctx, strings.Join(args, " "), s.APIManager.IngredientIdByName)
				if err != nil {
					return err
				}
				_, err = s.APIManager.CreateRecipe(ctx, r)
				return err
			},
		},
		&cobra.Command{
			Use:   "load-ingredients [file]",
			Short: "Load ingredient -> fdc mappings from afile",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				s, err := makeServer()
				if err != nil {
					return err
				}
				ctx := context.Background()

				mappings, err := api.IngredientMappingFromFile(ctx, strings.Join(args, ""))
				if err != nil {
					return err
				}

				err = s.APIManager.LoadIngredientMappings(ctx, mappings)
				return err
			},
		},
		&cobra.Command{
			Use:   "import [file]",
			Short: "import a recipe",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				s, err := makeServer()
				if err != nil {
					return err
				}
				ctx := context.Background()

				recipes, err := s.APIManager.RecipeFromFile(ctx, strings.Join(args, " "))
				if err != nil {
					return err
				}

				for x := range recipes {
					out, err := s.APIManager.CreateRecipe(ctx, &api.RecipeWrapperInput{Detail: recipes[x]})
					if err != nil {
						return err
					}
					res, err := glamour.Render(fmt.Sprintf(`
# Import Complete
Imported `+" **%s** as `%s`", out.Detail.Name, out.Detail.Id), "dark")
					if err != nil {
						return err
					}
					fmt.Print(res)
				}
				return nil
			},
		},
	)
}
