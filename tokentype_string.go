// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[L_BRACE-0]
	_ = x[R_BRACE-1]
	_ = x[L_SQUARE-2]
	_ = x[R_SQUARE-3]
	_ = x[COMMA-4]
	_ = x[COLON-5]
	_ = x[NULL-6]
	_ = x[FALSE-7]
	_ = x[TRUE-8]
	_ = x[NUMBER-9]
	_ = x[STRING-10]
	_ = x[EOF-11]
}

const _TokenType_name = "L_BRACER_BRACEL_SQUARER_SQUARECOMMACOLONNULLFALSETRUENUMBERSTRINGEOF"

var _TokenType_index = [...]uint8{0, 7, 14, 22, 30, 35, 40, 44, 49, 53, 59, 65, 68}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
