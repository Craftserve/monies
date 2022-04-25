package money

import "math"

// @NOTE: Not sure if this file is really that needed

func add(a, b int64) int64 {
	return a + b
}

func subtract(a, b int64) int64 {
	return a - b
}

func multiply(a, m int64) int64 {
	return a * m
}

func divide(a int64, d int64) int64 {
	return a / d
}

func modulus(a int64, d int64) int64 {
	return a % d
}

func allocate(a int64, r, s int) int64 {
	return a * int64(r) / int64(s)
}

func absolute(a int64) int64 {
	if a < 0 {
		return -a
	}

	return a
}

func negative(a int64) int64 {
	if a > 0 {
		return -a
	}

	return a
}

func round(a int64, e int) int64 {
	if a == 0 {
		return 0
	}

	absam := absolute(a)
	exp := int64(math.Pow(10, float64(e)))
	m := absam % exp

	if m > (exp / 2) {
		absam += exp
	}

	absam = (absam / exp) * exp

	if a < 0 {
		a = -absam
	} else {
		a = absam
	}

	return a
}
