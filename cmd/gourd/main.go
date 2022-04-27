package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/nickysemenza/gourd/api"
	"github.com/nickysemenza/gourd/client"
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
		log.Debug("cleaning up tracer")
		err := otel.GetTracerProvider().(*tracesdk.TracerProvider).ForceFlush(context.Background())
		if err != nil {
			log.Error(err)
		}
	}()
	ctx, span := otel.Tracer("client").Start(context.Background(), "main")
	defer span.End()

	err := rootCmd.ExecuteContext(ctx)
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
			RunE: func(cmd *cobra.Command, args []string) error {
				return runServer(cmd.Context())
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
				// res, err := s.APIManager.Notion.GetAll(context.Background())
				id := "rd_4b85d29a"
				_, err = s.APIManager.Latex(context.Background(), id)
				if err != nil {
					return err
				}
				// spew.Dump(res)
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
				lookbackDays := 30

				err = s.APIManager.DoSync(cmd.Context(), lookbackDays)
				if err != nil {
					return fmt.Errorf("sync: %w", err)
				}

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

				c, err := client.New("http://localhost:4242/api/")
				if err != nil {
					return err
				}
				res, err := c.ScrapeRecipeWithResponse(ctx, api.ScrapeRecipeJSONRequestBody{Url: url})

				log.Println(res.JSON201.Detail.Id)

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

				c, err := client.New("http://localhost:4242/api/")
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

				c, err := client.New("http://localhost:4242/api/")
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
					sb.WriteString(fmt.Sprintf("# %s \n `v%d` (%s)\n", r.Detail.Name, r.Detail.Version, r.Detail.Id))
				}
				res, err := glamour.Render(sb.String(), "dark")
				if err != nil {
					return err
				}
				fmt.Print(res)
				return nil
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

				ctx := cmd.Context()
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
