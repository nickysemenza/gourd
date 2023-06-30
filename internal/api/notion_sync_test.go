//go:build integration
// +build integration

package api

import (
	"context"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/nickysemenza/gourd/internal/clients/notion"
	"github.com/nickysemenza/gourd/internal/clients/notion/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestSync(t *testing.T) {
	require := require.New(t)
	apiManager := makeAPI(t)
	ctx := context.Background()
	err := apiManager.Sync(ctx, 14)
	require.NoError(err)

	items, _, err := apiManager.searchRecipes(ctx, DefaultPagination, "subpagetitle")
	require.NoError(err)
	require.Len(items, 1)
	rd := items[0].Detail.Id
	res, err := apiManager.recipeById(ctx, rd)
	require.NoError(err)

	require.Len(res.Detail.Sections, 1)
	require.Equal("bread", res.Detail.Sections[0].Ingredients[0].Ingredient.Ingredient.Name)
	require.Equal("eat", res.Detail.Sections[0].Instructions[0].Instruction)

	meals, err := apiManager.listMeals(ctx)
	require.NoError(err)
	require.Len(meals, 1)
	spew.Dump(meals)

	// l, err := apiManager.Latex(ctx, rd)
	// require.NoError(err)
	// require.Greater(len(l), 10000)

}

func TestSyncThingsChanged(t *testing.T) {
	nClient := mocks.NewClient(t)

	require := require.New(t)
	apiManager := makeAPI(t, WithNotionClient(nClient))
	ctx := context.Background()
	ctx = boil.WithDebug(ctx, true)

	base := time.Now().Truncate(time.Hour * 24).Add(-time.Hour * 24 * 5)
	r := notion.Recipe{
		Title: "example 1",
		UID:   "MEAL123",
		Time:  notion.GetTimeByName([]string{"dinner"}, base),
		Raw:   "name: example 1\n---\n2 cups flour",
		Tags:  []string{"a"},
	}

	r.Time = notion.GetTimeByName([]string{"dinner"}, base.Add(time.Hour*24*2))
	r.Raw = "name: example 1\n---\n3 cups flour"
	r.Tags = []string{"b"}
	nClient.On("GetAll", mock.Anything, time.Hour*24*14, "").
		Return([]notion.Recipe{r}, nil).Once()

	require.NoError(apiManager.Sync(ctx, 14))

	nClient.On("GetAll", mock.Anything, time.Hour*24*14, "").
		Return([]notion.Recipe{r}, nil).Once()
	require.NoError(apiManager.Sync(ctx, 14))

}
