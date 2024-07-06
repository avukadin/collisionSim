
package main

import (
	"fmt"
	"image/color"
)

func hexToRGBA(hex string) color.RGBA {
	var r, g, b uint8
	if _, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b); err != nil {
		return color.RGBA{}
	}
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func minVal[T int|float32|float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func maxVal[T int|float32|float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}
