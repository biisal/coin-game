package main

import (
	"testing"

	"github.com/biisal/fun-game/utils"
)

// Benchmark for DrawWithMath
//
//	func BenchmarkDrawWithMath(b *testing.B) {
//		for b.Loop() {
//			utils.DrawWithMath()
//		}
//	}
func BenchmarkMove(b *testing.B) {
	for b.Loop() {
		utils.Move(0, 0)
	}
}

// Benchmark for Draw (builder version)
// func BenchmarkDraw(b *testing.B) {
// 	for b.Loop() {
// 		utils.Draw()
// 	}
// }
