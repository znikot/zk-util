package misc

// simulate ternary operation
//
//	val := If(true, 0, 1) // val = 0
func If[T any](condition bool, trueVal T, falseVal T) T {
	if condition {
		return trueVal
	}

	return falseVal
}

// simulate ternary operation with function
//
//	val := IfFunc(true, func() int { return 0 }, func() int { return 1 }) // val = 0
func IfFunc[T any](condition bool, trueVal func() T, falseVal func() T) T {
	if condition {
		return trueVal()
	}

	return falseVal()
}
