// Code generated by "stringer -type=TokenKind"; DO NOT EDIT.

package parser

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[MeasureWeight-0]
	_ = x[MeasureVolume-1]
	_ = x[UnkString-2]
	_ = x[KindIngredientName-3]
	_ = x[Modifier-4]
	_ = x[UnkFloat-5]
	_ = x[WeightFloat-6]
	_ = x[VolumeFloat-7]
	_ = x[None-8]
}

const _TokenKind_name = "MeasureWeightMeasureVolumeUnkStringKindIngredientNameModifierUnkFloatWeightFloatVolumeFloatNone"

var _TokenKind_index = [...]uint8{0, 13, 26, 35, 53, 61, 69, 80, 91, 95}

func (i TokenKind) String() string {
	if i < 0 || i >= TokenKind(len(_TokenKind_index)-1) {
		return "TokenKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenKind_name[_TokenKind_index[i]:_TokenKind_index[i+1]]
}
