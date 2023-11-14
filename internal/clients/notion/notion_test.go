package notion

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSmoke(t *testing.T) {
	c := NewFakeNotion(t)
	res, err := c.GetAll(context.Background(), 14, "")
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, "page1title", res[0].Title)
	require.Len(t, res[0].Photos, 1)
}
