package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    *Ingredient
		wantErr bool
	}{
		{name: "empty", wantErr: true},
		{name: "malformed 1", s: "flour", wantErr: true},
		{name: "malformed 2", s: "cup", wantErr: true},
		{name: "malformed 3", s: "1", wantErr: true},
		{name: "malformed 3", s: " 1 1/e2 cup", wantErr: true},
		{name: "malformed 3", s: "\a", wantErr: true},
		{name: "basic vol", s: "1 cup flour", want: &Ingredient{Name: "flour", Volume: Measurement{Value: 1, Unit: "cup"}}},
		{name: "basic vol decimal", s: "1.2 cup flour", want: &Ingredient{Name: "flour", Volume: Measurement{Value: 1.2, Unit: "cup"}}},
		{name: "basic vol fraction", s: "1 ½ cup flour", want: &Ingredient{Name: "flour", Volume: Measurement{Value: 1.5, Unit: "cup"}}},
		{name: "basic mass", s: "100 grams flour", want: &Ingredient{Name: "flour", Weight: Measurement{Value: 100, Unit: "grams"}}},
		{name: "basic mass", s: "100 gr flour", want: &Ingredient{Name: "flour", Weight: Measurement{Value: 100, Unit: "gr"}}},
		{name: "basic mass and vol", s: "1/2 cup (60 grams) flour", want: &Ingredient{Name: "flour",
			Weight: Measurement{Value: 60, Unit: "grams"},
			Volume: Measurement{Value: 0.5, Unit: "cup"},
		}},
		{name: "mass and vol, modifier", s: "¾ cup (90 g) flour, sifted", want: &Ingredient{Name: "flour",
			Weight:   Measurement{Value: 90, Unit: "g"},
			Volume:   Measurement{Value: 0.75, Unit: "cup"},
			Modifier: "sifted",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(context.Background(), tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				// return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseAndStringify(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{in: "3 tablespoons/43 grams unsalted butter, melted", out: "43 grams (3 tablespoons) unsalted butter, melted"},
		{in: "1 ¼ teaspoons kosher salt", out: "1.25 teaspoons kosher salt"},
		{in: "8 ounces (225g) semisweet chocolate, roughly chopped with a knife into 1/2- to 1/4-inch chunks", out: "225 g (8 ounces) semisweet chocolate, chunks"}, // TODO: it's dropping some
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := Parse(context.Background(), tt.in)
			require.NoError(t, err)

			require.Equal(t, tt.out, got.ToString())
		})
	}
}
