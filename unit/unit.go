package unit

import "strings"

type Unit int //

const (
	Gram Unit = iota
	Ounce

	Liter
	Quart
	Cup
	Teaspoon
	Tablespoon

	Item
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
	}
	volumes = map[string]Unit{
		"c":   Cup,
		"cup": Cup,

		"tsp":      Teaspoon,
		"teaspoon": Teaspoon,

		"tbsp":       Tablespoon,
		"tablespoon": Tablespoon,

		"l":     Liter,
		"liter": Liter,

		"q":     Quart,
		"quart": Quart,
	}
)

func IsWeight(unit string) bool {
	_, ok := weights[strings.TrimSuffix(unit, "s")]
	return ok
}
func IsVolume(unit string) bool {
	_, ok := volumes[strings.TrimSuffix(unit, "s")]
	return ok
}

func Parse(unit string) Unit {
	vol, isVol := volumes[strings.TrimSuffix(unit, "s")]
	if isVol {
		return vol
	}
	weight, isWeight := weights[strings.TrimSuffix(unit, "s")]
	if isWeight {
		return weight
	}
	return Item
}
