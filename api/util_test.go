package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecipeFromFileYAMLvsJSON(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	r1, err := RecipeFromFile(ctx, "../testdata/cookies_1.json")
	require.NoError(err)
	r2, err := RecipeFromFile(ctx, "../testdata/cookies_1.yaml")
	require.NoError(err)
	require.Equal(r1, r2)
	require.Equal("cookies 1", r1[0].Name)
}
