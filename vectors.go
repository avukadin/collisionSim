package main

import (
	"math"
)

type Vector struct {
	x, y float64
}

func (v *Vector) subtract(v2 Vector) Vector {
	return Vector{v.x - v2.x, v.y - v2.y}
}

func (v *Vector) dot(v2 Vector) float64 {
	return v.x*v2.x + v.y*v2.y
}

func (v *Vector) multiply(scalar float64) Vector {
	return Vector{v.x * scalar, v.y * scalar}
}

func (v *Vector) squaredDist(v2 Vector) float64 {
	return (v.x-v2.x)*(v.x-v2.x) + (v.y-v2.y)*(v.y-v2.y)
}

func (v *Vector) distance(v2 Vector) float64 {
	return math.Sqrt(v.squaredDist(v2))
}

func (v * Vector) projection(v2 Vector) Vector {
	return v2.multiply(v.dot(v2)/v2.dot(v2))
}

func (v *Vector) add(v2 Vector) Vector {
	return Vector{v.x + v2.x, v.y + v2.y}
}

func (v *Vector) normalize() Vector {
	return v.multiply(1/v.distance(Vector{0, 0}))
}

