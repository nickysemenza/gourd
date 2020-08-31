//go:generate stringer -type=Unit -linecomment

package unit

import "strings"

type Unit int //

const (
	Gram  Unit = iota // gram
	Ounce             // oz
	Pound             // pound

	Liter      // liter
	Millileter // millileter
	Quart      // quart
	Cup        // cpu
	Teaspoon   // tsp
	Tablespoon // tbsp

	Item // item
)

// mapping for parsing
// nolint: gochecknoglobals
var (
	weights = map[string]Unit{
		"g":    Gram,
		"gr":   Gram,
		"gram": Gram,

		"oz":    Ounce,
		"ounce": Ounce,

		"lb":    Pound,
		"pound": Pound,
	}
	volumes = map[string]Unit{
		"c":   Cup,
		"cup": Cup,

		"tsp":      Teaspoon,
		"teaspoon": Teaspoon,

		"tbsp":       Tablespoon,
		"tablespoon": Tablespoon,

		"l":          Liter,
		"liter":      Liter,
		"millileter": Millileter,
		"ml":         Millileter,

		"q":     Quart,
		"quart": Quart,

		"oz":    Ounce,
		"ounce": Ounce,
	}
)

func normalize(unit string) string {
	return strings.TrimSuffix(strings.ToLower(unit), "s")
}

func IsWeight(unit string) bool {
	_, ok := weights[normalize(unit)]
	return ok
}
func IsGram(unit string) bool {
	val, ok := weights[normalize(unit)]
	if !ok {
		return false
	}
	return val == Gram
}
func IsVolume(unit string) bool {
	_, ok := volumes[normalize(unit)]
	return ok
}

func Parse(unit string) Unit {
	vol, isVol := volumes[normalize(unit)]
	if isVol {
		return vol
	}
	weight, isWeight := weights[normalize(unit)]
	if isWeight {
		return weight
	}
	return Item
}
