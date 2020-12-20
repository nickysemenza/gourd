package unit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint:scopelint
func TestIsWeightVolume(t *testing.T) {
	tests := []struct {
		name     string
		unit     string
		isWeight bool
		isVolume bool
	}{
		{unit: "gram", isWeight: true},
		{unit: "grams", isWeight: true},
		{unit: "cup", isVolume: true},
		{unit: "cups", isVolume: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.isWeight, IsWeight(tt.unit))
			require.Equal(t, tt.isVolume, IsVolume(tt.unit))
		})
	}
}
