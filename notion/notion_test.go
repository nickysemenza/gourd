package notion

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSmoke(t *testing.T) {

	c := NewFakeNotion(t)
	res, err := c.GetAll(context.Background())
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, res[0].Title, "page1title")
	require.Len(t, res[0].Photos, 1)
}
