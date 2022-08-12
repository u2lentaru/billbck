package models

import (
	"testing"
)

func BenchmarkEmptySlice(b *testing.B) {
	bench_arr := []CableResistance{}

	for i := 0; i < b.N; i++ {
		bench_arr = append(bench_arr, CableResistance{Id: 7, CableResistanceName: "95 мм", Resistance: 0.195, MaterialType: false})
	}

}

func BenchmarkMakeSlice(b *testing.B) {
	bench_arr := make([]CableResistance, 0, b.N)

	for i := 0; i < b.N; i++ {
		bench_arr = append(bench_arr, CableResistance{Id: 7, CableResistanceName: "95 мм", Resistance: 0.195, MaterialType: false})
	}

}
