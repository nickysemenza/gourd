package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/nickysemenza/gourd/common"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/manager"
	"github.com/stretchr/testify/require"
)

func TestRecipeFromFile(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	r, err := RecipeFromFile(ctx, "../testdata/cookies_1.json")
	require.NoError(err)
	require.Equal("cookies 1", r.Detail.Name)
	baseName := fmt.Sprintf("%s-%s", r.Detail.Name, common.UUID())
	r.Detail.Name = baseName

	tdb := db.NewDB(t)
	m := manager.New(tdb, nil, nil)
	apiManager := NewAPI(m)

	r2, err := apiManager.CreateRecipe(ctx, r)
	require.NoError(err)

	require.Equal(baseName, r2.Detail.Name)
	r.Detail.Id = "" // reset so we create a dup instead of update, ptr

	r3, err := apiManager.CreateRecipe(ctx, r)
	require.NoError(err)

	require.Equal((*r2.Detail.Version)+1, *r3.Detail.Version)
}
func TestRecipeReferencingRecipe(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	r, err := RecipeFromFile(ctx, "../testdata/dep_1.yaml")
	require.NoError(err)
	tdb := db.NewDB(t)
	m := manager.New(tdb, nil, nil)
	apiManager := NewAPI(m)
	_, err = apiManager.CreateRecipe(ctx, r)
	require.NoError(err)
}

func TestRecipeFromFileYAMLvsJSON(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	r1, err := RecipeFromFile(ctx, "../testdata/cookies_1.json")
	require.NoError(err)
	r2, err := RecipeFromFile(ctx, "../testdata/cookies_1.yaml")
	require.NoError(err)
	require.Equal(r1, r2)
}
