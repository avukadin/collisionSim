package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	sPerU = 1.0/60 // seconds per update
	pxPerM = 50 // pixels per meter

	ag = 9.8*pxPerM*sPerU*sPerU/5 // acceleration due to gravity
	coeff = 0.9
	r1 = 30
	r2 = 3
	nImages = 6
)

func (g *Game) loadImages(folder string) {
	for i := 1; i <= nImages; i++ {
		filePath := fmt.Sprintf("%s/ball%d.png", folder, i)
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}

		ebitenImg := ebiten.NewImageFromImage(img)
		g.imgs = append(g.imgs, ebitenImg)
	}
}

type Game struct {
	balls []Ball
	callDuration time.Duration
	imgs []*ebiten.Image
}

func main() {
	ebiten.SetWindowTitle("Particle Collision Simulation")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	g := Game{}
	g.loadImages("./images")
	// balls := GetFallingBalls(100, 10, 3)
	balls := GetFallingSingleBall(100, r1)
	// balls = append(balls, GetBallsOnGround(500, 10, r2)...)
	balls = append(balls, GetBallsRandom(1000, 1, r2)...)
	for i := range balls {
		balls[i].id = i
	}
	g.balls = balls
	err := ebiten.RunGame(&g)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f, Call Time %v", ebiten.ActualFPS(), g.callDuration))
	for i, ball := range g.balls {
		// Scale image to 10 by 10
		img := g.imgs[i % nImages]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2*ball.radius/float64(img.Bounds().Dx()), 2*ball.radius/float64(img.Bounds().Dy()))
		op.GeoM.Translate(ball.r.x-ball.radius, ball.r.y-ball.radius)
		screen.DrawImage(img, op)
	}
}

func isHit(b1, b2 *Ball) bool {
	dist := b1.r.distance(b2.r)
	return dist <= (b1.radius + b2.radius)
}


func (b *Ball) handleWallCollision() {
	if b.r.x-b.radius < 0 || b.r.x+b.radius > screenWidth {
		b.v.x = -coeff*b.v.x
	}
	if b.r.y-b.radius < 0 || b.r.y+b.radius > screenHeight {
		b.v.y = -coeff*b.v.y
	}
	//Move out of wall
	if b.r.x-b.radius < 0 {
		b.r.x = b.radius
	}
	if b.r.x+b.radius > screenWidth {
		b.r.x = screenWidth - b.radius
	}
	if b.r.y-b.radius < 0 {
		b.r.y = b.radius
	}
	if b.r.y+b.radius > screenHeight {
		b.r.y = screenHeight - b.radius
	}
}

func (g *Game) Update() error {
	// g.UpdateBrute()
	g.UpdateGrid()
	return nil
}

func (g *Game) UpdateBrute() error {	
	t0 := time.Now()
	for i := range g.balls {
		g.balls[i].addGravity()
	}	
	calls := 0
	for i := 0; i < len(g.balls); i++ {
		for j := i+1; j < len(g.balls); j++ {
			b1 := &g.balls[i]
			b2 := &g.balls[j]
			b1.handleWallCollision()
			calls++
			if isHit(b1, b2) {
				v1 := b1.getNewVelocity(b2)
				v2 := b2.getNewVelocity(b1)
				b1.v = v1
				b2.v = v2

				// move balls apart
				dist := b1.r.distance(b2.r)
				overlap := b1.radius + b2.radius - dist
				dir := b1.r.subtract(b2.r)
				dir = dir.normalize()

				b1.r = b1.r.add(dir.multiply(overlap/2))
				b2.r = b2.r.add(dir.multiply(-overlap/2))
			}
		}
	}

	for i := range g.balls {
		g.balls[i].r = g.balls[i].r.add(g.balls[i].v)
	}
	g.callDuration = time.Since(t0)

	return nil
}


func (g *Game) UpdateGrid() error {	
	to := time.Now()
	for i := range g.balls {
		g.balls[i].addGravity()
	}	
	grid := NewUniforGrid(g.balls, math.Max(r1, r2)*2)
	for i := 0; i < len(g.balls); i++ {
		b1 := &g.balls[i]
		neighbors := grid.GetNeighbors(b1)
		b1.handleWallCollision()
		for _, b2 := range neighbors {

			if b1.id == b2.id {
				continue
			}

			if isHit(b1, b2) {
				v1 := b1.getNewVelocity(b2)
				v2 := b2.getNewVelocity(b1)
				b1.v = v1
				b2.v = v2

				// move balls apart
				dist := b1.r.distance(b2.r)
				overlap := b1.radius + b2.radius - dist
				dir := b1.r.subtract(b2.r)
				dir = dir.normalize()

				b1.r = b1.r.add(dir.multiply(overlap/2))
				b2.r = b2.r.add(dir.multiply(-overlap/2))
			}
		}
	}

	for i := range g.balls {
		g.balls[i].r = g.balls[i].r.add(g.balls[i].v)
	}

	g.callDuration = time.Since(to)
			
	return nil
}

