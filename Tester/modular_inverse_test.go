// modular_inverse_test.go
package main

import (
	"testing"
)

func BenchmarkModInv_1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ModInv_1(3)
	}
}
func BenchmarkModInv_2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ModInv_2(3)
	}
}
func BenchmarkModInv_3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ModInv_3(3)
	}
}
func BenchmarkModInv_4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ModInv_4(3)
	}
}
