package main

import (
	"math"
	"math/rand"
)

func getBalls(n int, mass float64, radius float64) []Ball{
	balls := make([]Ball, n)
	for i := 0; i < n; i++ {
		ball := Ball{
			r: Vector{0, 0},
			v: Vector{0, 0},
			radius: radius,
			mass: mass,
		}
		balls[i] = ball
	}
	return balls
}

func Get1BallForward(mass float64, radius float64, padding float64) []Ball{
	balls := getBalls(1, mass, radius)
	balls[0].r = Vector{radius + padding + 1, screenHeight/2}

	var speed float64 = 5
	balls[0].v = Vector{speed, 0}

	return balls
}
func Get2BallsToward(mass float64, radius float64, padding float64) []Ball{
	balls := getBalls(2, mass, radius)
	balls[0].r = Vector{radius + padding + 1, screenHeight/2}
	balls[1].r = Vector{screenWidth - radius - padding - 1, screenHeight/2}

	var speed float64 = 5
	balls[0].v = Vector{speed, 0}
	balls[1].v = Vector{-speed, 0}

	return balls
}

func Get2BallsSame(mass float64, radius float64, padding float64) []Ball{
	balls := getBalls(2, mass, radius)
	balls[0].r = Vector{radius + padding + 1, screenHeight/2}
	balls[1].r = Vector{screenWidth/4 + radius + padding, screenHeight/2}

	var speed float64 = 5
	balls[0].v = Vector{speed, 0}
	balls[1].v = Vector{speed/2, 0}

	return balls
}

func GetFallingSingleBall(mass float64, radius float64, padding float64) []Ball{
	balls := getBalls(1, mass, radius)
	balls[0].r = Vector{screenWidth/2, radius + padding}
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

func GetBallsOnGround(n int, mass float64, radius float64, padding float64) []Ball{
	balls := getBalls(n, mass, radius)
	var maxBallsPerLevel float64 = math.Floor((screenWidth-2*padding)/(2*radius))
	var levels float64 = math.Ceil(float64(n)/maxBallsPerLevel)
	for level := 0; level < int(levels); level++ {
		var index float64 = 0;
		for i := level*int(maxBallsPerLevel); i < (level+1)*int(maxBallsPerLevel); i++ {
			if i >= n {
				break
			}
			balls[i].r = Vector{padding + radius + index*2*radius + float64(level%2)*radius, screenHeight-padding-radius*2*float64(level)}
			balls[i].v = Vector{0, 0}
			index++
		}
	}

	return balls
}

func GetBallsRandom(n int, mass float64, radius float64, padding float64) []Ball{
	balls := getBalls(n, mass, radius)
	maxX := screenWidth-padding-radius-1
	maxY := screenHeight-padding-radius-1
	minX := padding+radius+1
	minY := padding+radius+1
	for i := 0; i < n; i++ {
		balls[i].r = Vector{(rand.Float64() * (maxX - minX)) + minX, (rand.Float64() * (maxY - minY)) + minY}
		balls[i].v = Vector{rand.Float64()*2 - 1, rand.Float64()*2 - 1}
	}

	return balls
}

