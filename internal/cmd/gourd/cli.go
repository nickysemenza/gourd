package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/nickysemenza/gourd/internal/api"
	"github.com/nickysemenza/gourd/internal/clients/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

func newClient() (*client.Client, error) {
	return client.New("http://localhost:4242/api/")
}

// nolint:gochecknoglobals
var rootCmd = &cobra.Command{
	Use:   "gourd",
	Short: "Go Universal Recipe Database",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		mode := "client"
		if cmd.Name() == "server" {
			mode = "server"
		}

		if err := setupMisc(mode); err != nil {
			return err
		}
		ctx, span := otel.Tracer(mode).Start(cmd.Context(), cmd.Name())
		defer span.End()
		cmd.SetContext(ctx)

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func makeRenderer() *glamour.TermRenderer {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	return r
}

// nolint:gochecknoinits
func init() {

	rootCmd.AddCommand(
		importCommand,
		&cobra.Command{
			Use: "version",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("v0")
			},
		},
		&cobra.Command{
			Use:   "server",
			Short: "Run the server",
			RunE: func(cmd *cobra.Command, args []string) error {
				return runServer(cmd.Context())
			},
		},
		&cobra.Command{
			Use:   "tmp",
			Short: "misc",
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				s, err := makeServer(ctx)
				if err != nil {
					return err
				}
				id := "rd_4b85d29a"
				_, err = s.APIManager.Latex(ctx, id)
				if err != nil {
					return err
				}
				return nil
			},
		},
		&cobra.Command{
			Use:   "usda",
			Short: "usda",
			RunE: func(cmd *cobra.Command, args []string) error {
				api.Test()
				return nil
				// ctx := cmd.Context()
				// s, err := makeServer(ctx)
				// if err != nil {
				// 	return err
				// }
				// id := "rd_4b85d29a"
				// _, err = s.APIManager.Latex(ctx, id)
				// if err != nil {
				// 	return err
				// }
				// return nil
			},
		},
		&cobra.Command{
			Use:   "sync",
			Short: "Run the server",
			RunE: func(cmd *cobra.Command, args []string) error {
				c, err := newClient()
				if err != nil {
					return err
				}
				res, err := c.DoSyncWithResponse(cmd.Context(), &api.DoSyncParams{LookbackDays: 7})
				if err != nil {
					return fmt.Errorf("sync: %w", err)
				}
				fmt.Printf("%s", res.Body)

				return nil
			},
		},

		&cobra.Command{
			Use:   "scrape [url]",
			Short: "scrape a recipe",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				url := strings.Join(args, " ")

				ctx, span := otel.Tracer("client").Start(cmd.Context(), "scrape")
				defer span.End()

				c, err := newClient()
				if err != nil {
					return err
				}
				res, err := c.ScrapeRecipeWithResponse(ctx, api.ScrapeRecipeJSONRequestBody{Url: url})
				if err != nil {
					return err
				}

				log.Println(res.JSON201.Detail.Id)

				out, err := makeRenderer().Render(fmt.Sprintf("# scrape Complete\nScraped `%s` as `%s`", url, res.JSON201.Detail.Id))
				if err != nil {
					return err
				}
				fmt.Print(out)

				return err
			},
		},
		&cobra.Command{
			Use:   "load-ingredients [file]",
			Short: "Load ingredient -> fdc mappings from afile",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {

				ctx, span := otel.Tracer("client").Start(cmd.Context(), "load-mappings")
				defer span.End()

				mappings, err := api.IngredientMappingFromFile(ctx, strings.Join(args, ""))
				if err != nil {
					return err
				}

				c, err := newClient()
				if err != nil {
					return err
				}
				_, err = c.LoadIngredientMappings(ctx, api.LoadIngredientMappingsJSONRequestBody{IngredientMappings: mappings})
				if err != nil {
					return err
				}
				log.Infof("loaded %d mappings", len(mappings))
				return nil
			},
		},
		&cobra.Command{
			Use: "list-recipes",
			RunE: func(cmd *cobra.Command, args []string) error {

				ctx, span := otel.Tracer("client").Start(cmd.Context(), "listrecipes")
				defer span.End()

				c, err := newClient()
				if err != nil {
					return err
				}

				var l api.LimitParam = 10
				resp, err := c.ListRecipesWithResponse(ctx, &api.ListRecipesParams{Limit: &l})
				if err != nil {
					return err
				}
				var sb strings.Builder

				for _, r := range *resp.JSON200.Recipes {
					sb.WriteString(fmt.Sprintf("# %s \n `v%d` (%s)\n", r.Detail.Name, r.Detail.Meta.Version, r.Detail.Id))
				}
				res, err := makeRenderer().Render(sb.String())
				if err != nil {
					return err
				}
				fmt.Print(res)
				return nil
			},
		},
	)
}

var importCommand = &cobra.Command{
	Use:   "import [file]",
	Short: "import a recipe",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		s, err := makeServer(ctx)
		if err != nil {
			return err
		}

		recipes, err := s.APIManager.RecipeFromFile(ctx, strings.Join(args, " "))
		if err != nil {
			return err
		}

		for x := range recipes {
			out, err := s.APIManager.CreateRecipe(ctx, &api.RecipeWrapperInput{Detail: recipes[x]})
			if err != nil {
				return err
			}
			res, err := makeRenderer().Render(fmt.Sprintf("# Import Complete\nImported **%s** as `%s`", out.Detail.Name, out.Detail.Id))
			if err != nil {
				return err
			}
			fmt.Print(res)
		}
		return nil
	},
}
