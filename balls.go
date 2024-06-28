package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
)

var colors = []color.Color{
	hexToRGBA("0096c7"),
	hexToRGBA("f72585"),
	hexToRGBA("7209b7"),
	hexToRGBA("ffd166"),
	hexToRGBA("f4a261"),
	hexToRGBA("00cecb"),
}

func hexToRGBA(hex string) color.RGBA {
	var r, g, b uint8
	if _, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b); err != nil {
		return color.RGBA{}
	}
	return color.RGBA{R: r, G: g, B: b, A: 155}
}

type Ball struct {
	r Vector // position
	v Vector // velocity

	radius float64
	mass float64
	color color.Color
	id int
}

func (b *Ball) addGravity() {
	b.v.y += ag
	b.r.y += float64(b.v.y)
}

func (b *Ball) getNewVelocity(b2 *Ball) Vector {
	mRatio := (coeff+1)*b2.mass/(b.mass+b2.mass)
	vDiff := b.v.subtract(b2.v)
	rDiff := b.r.subtract(b2.r)
	proj := vDiff.projection(rDiff)
	return b.v.subtract(proj.multiply(mRatio))
}

func getBalls(n int, mass float64, radius float64) []Ball{
	balls := make([]Ball, n)
	for i := 0; i < n; i++ {
		ball := Ball{
			r: Vector{0, 0},
			v: Vector{0, 0},
			radius: radius,
			mass: mass,
			color: colors[i%len(colors)],
		}
		balls[i] = ball
	}
	return balls
}

func GetFallingSingleBall(mass float64, radius float64) []Ball{
	n:=1
	balls := getBalls(n, mass, radius)
	var ballsPerLevel float64 = screenWidth/(2*radius)
	var levels int = int(math.Ceil(float64(n)/float64(ballsPerLevel)))
	var loc float64 = radius
	for level := 0; level < levels; level++ {
		for i := level*int(ballsPerLevel); i < (level+1)*int(ballsPerLevel); i++ {
			if i >= n {
				break
			}
			loc = screenWidth/2
			wiggle := rand.Float64() * 2 * radius - radius
			balls[i].r = Vector{loc + wiggle, radius*float64(level)}
			balls[i].v = Vector{0, 0}
		}
	}

	return balls

}

func GetFallingBalls(n int, mass float64, radius float64) []Ball{
	balls := getBalls(n, mass, radius)
	var ballsPerLevel float64 = screenWidth/(2*radius)
	var levels int = int(math.Ceil(float64(n)/float64(ballsPerLevel)))
	var loc float64 = radius
	interval := (screenWidth-2*radius)/float64(n-1)
	for level := 0; level < levels; level++ {
		for i := level*int(ballsPerLevel); i < (level+1)*int(ballsPerLevel); i++ {
			if i >= n {
				break
			}
			wiggle := rand.Float64() * 2 * radius - radius
			balls[i].r = Vector{loc + wiggle, radius*float64(level)}
			balls[i].v = Vector{0, 0}
			loc += interval+wiggle
		}
	}

	return balls

}


func GetBallsOnGround(n int, mass float64, radius float64) []Ball{
	balls := getBalls(n, mass, radius)
	var ballsPerLevel float64 = screenWidth/(2*radius)
	var levels int = int(math.Ceil(float64(n)/float64(ballsPerLevel)))
	var loc float64 = radius
	interval := (screenWidth-2*radius)/float64(n-1)
	for level := 0; level < levels; level++ {
		for i := level*int(ballsPerLevel); i < (level+1)*int(ballsPerLevel); i++ {
			if i >= n {
				break
			}
			wiggle := rand.Float64() * 2 * radius - radius
			balls[i].r = Vector{loc + wiggle, screenHeight - radius*float64(level)}
			balls[i].v = Vector{0, 0}
			loc += interval+wiggle
		}
	}

	return balls
}

func GetBallsRandom(n int, mass float64, radius float64) []Ball{
	balls := getBalls(n, mass, radius)
	for i := 0; i < n; i++ {
		balls[i].r = Vector{rand.Float64()*screenWidth, rand.Float64()*screenHeight}
		if rand.Float64() < 0.5 {
			balls[i].v = Vector{1, 0}
		} else {
			balls[i].v = Vector{-1, 0}
		}
	}


	return balls
}
