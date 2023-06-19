package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecipeFromFileYAMLvsJSON(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	a := API{}
	r1, err := a.RecipeFromFile(ctx, "../../tooling/testdata/cookies_1.json")
	require.NoError(err)
	require.Len(r1, 1)
	r2, err := a.RecipeFromFile(ctx, "../../tooling/testdata/cookies_1.yaml")
	require.NoError(err)
	require.Len(r2, 1)
	r1[0].Sources = nil
	r2[0].Sources = nil
	require.Equal(r1, r2)
	require.Equal("cookies 1", r1[0].Name)
}
