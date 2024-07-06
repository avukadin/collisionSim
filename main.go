package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1920*0.8
	screenHeight = 1080*0.7

	pxPerM = 50 // pixels per meter
	ag = 9.8/(60*60)*pxPerM

	nImages = 1
)

func (g *Game) loadImages(folder string) {
	for i := 1; i <= nImages; i++ {
		filePath := fmt.Sprintf("%s/test%d.png", folder, i)
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
	imgs []*ebiten.Image

	padding float64

	gravSlider Slider
	restSlider Slider

	callDuration time.Duration
	maxRadius float64
	frame int
}

type Slider struct {
	x, y, width, height	int
	value			float64
	dragging		bool
	label			string
}

func main() {
	ebiten.SetWindowTitle("Particle Collision Simulation")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	g := Game{
		padding: 60,
		gravSlider: Slider{x: 100, y: screenHeight-50, width: 200, height: 3, value: 1, label: "Gravity"},
		restSlider: Slider{x: 400, y: screenHeight-50, width: 200, height: 3, value: 0.95, label: "Restitution"},
	}
	g.loadImages("./images")

	balls := make([]Ball, 0)
	// balls = append(balls, GetFallingSingleBall(100, 100, g.padding)...)
	// balls = append(balls, Get2BallsToward(2, 100, g.padding)...)
	balls = append(balls, GetBallsRandom(2000,1, 3, g.padding)...)
	// balls = append(balls, GetBallsOnGround(20000,1, 2, g.padding)...)

	for i := range balls {
		balls[i].id = i
	}

	g.balls = balls
	for i := range g.balls {
		g.maxRadius = maxVal(g.maxRadius, g.balls[i].radius)
	}

	err := ebiten.RunGame(&g)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f, Call Time %v", ebiten.ActualFPS()/2, g.callDuration))
	bgColor := color.RGBA{6, 10, 22, 255}
	screen.Fill(bgColor)
	for i, ball := range g.balls {
		// Scale image to 10 by 10
		img := g.imgs[i % nImages]
		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterLinear
		op.GeoM.Scale(1.25*2*ball.radius/float64(img.Bounds().Dx()), 1.25*2*ball.radius/float64(img.Bounds().Dy()))
		op.GeoM.Translate(ball.r.x-ball.radius, ball.r.y-ball.radius)
		screen.DrawImage(img, op)
	}
	g.gravSlider.Draw(screen)
	g.restSlider.Draw(screen)
}

func isHit(b1, b2 *Ball) bool {
	dist := b1.r.distance(b2.r)
	return dist <= (b1.radius + b2.radius)
}

func (g *Game) Update() error {
	g.UpdateBrute()
	// g.UpdateGrid()

	g.gravSlider.Update()
	g.restSlider.Update()

	return nil
}

func (g *Game) UpdateBrute() error {	
	g.frame ++
	t0 := time.Now()
	for i := 0; i < len(g.balls); i++ {
		for j := i+1; j < len(g.balls); j++ {
			b1 := &g.balls[i]
			b2 := &g.balls[j]
			if isHit(b1,b2){
				v1 := b1.ballHitVelocity(
					b2, 
					g.restSlider.value,
				)
				v2 := b2.ballHitVelocity(b1, 
					g.restSlider.value)
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
		g.balls[i].addGravity(g.gravSlider.value)
		g.balls[i].handleWallCollision(g.padding, g.restSlider.value)
		g.balls[i].r = g.balls[i].r.add(g.balls[i].v)
	}
	g.callDuration = time.Since(t0)

	return nil
}


func (g *Game) UpdateGrid() error {	
	to := time.Now()
	grid := NewUniforGrid(g.balls, g.maxRadius*2)
	for i := 0; i < len(g.balls); i++ {
		b1 := &g.balls[i]
		neighbors := grid.GetNeighbors(b1)
		b1.handleWallCollision(g.padding, g.restSlider.value)
		for _, b2 := range neighbors {

			if b1.id == b2.id {
				continue
			}

			if isHit(b1, b2) {
				v1 := b1.ballHitVelocity(b2, g.restSlider.value)
				v2 := b2.ballHitVelocity(b1, g.restSlider.value)
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
		g.balls[i].addGravity(g.gravSlider.value)
		g.balls[i].handleWallCollision(g.padding, g.restSlider.value)
		g.balls[i].r = g.balls[i].r.add(g.balls[i].v)
	}

	g.callDuration = time.Since(to)
			
	return nil
}








































