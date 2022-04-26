package monies

import "math"

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
