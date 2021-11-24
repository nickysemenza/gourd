package notion

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestSmoke(t *testing.T) {

	c := NewFakeNotion(t)
	res, err := c.Dump(context.Background())
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, res[0].Title, "test1")
	spew.Dump(res)
	require.Len(t, res[0].Photos, 1)
}
