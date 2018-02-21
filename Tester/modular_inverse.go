// modular_inverse.go
package main

/*
	Newton's method summarized by Marc B. Reynolds from various authors:
	http://marc-b-reynolds.github.io/math/2017/09/18/ModInverse.html
*/

func ModInv_1(a uint8) uint8 {
	x := a       //  3 bits
	x *= 2 - a*x //  6
	x *= 2 - a*x // 12
	return x
}

func ModInv_2(a uint8) uint8 {
	x := (a * a) + a - 1 //  4 bits (For serial comment below: a*a & a-1 are independent)
	x *= 2 - a*x         //  8
	return x
}

func ModInv_3(a uint8) uint8 {
	x := 3*a ^ 2 //  5 bits
	x *= 2 - a*x // 10
	return x
}

func ModInv_4(a uint8) uint8 {
	u := 2 - a
	i := a - 1
	i *= i
	u *= i + 1
	i *= i
	u *= i + 1
	return u
}
