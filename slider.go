package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (s *Slider) Update() {
    mx, my := ebiten.CursorPosition()
    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        if !s.dragging {
            if mx >= s.x-10 && mx <= s.x+s.width + 10 && my >= s.y-10 && my <= s.y+s.height+10 {
                s.dragging = true
            }
        }
    } else {
        s.dragging = false
    }

    if s.dragging {
        s.value = float64(mx-s.x) / float64(s.width)
        if s.value < 0 {
            s.value = 0
        }
        if s.value > 1 {
            s.value = 1
        }
    }
}

func (s *Slider) Draw(screen *ebiten.Image) {
	// Draw the slider background
	vector.DrawFilledRect(screen, float32(s.x), float32(s.y), float32(s.width), float32(s.height), hexToRGBA("A8DADC"), true)

	//Draw the slider handle as scircle
	handleX := s.x + int(s.value*float64(s.width)) - s.height/2
	vector.DrawFilledCircle(screen, float32(handleX), float32(s.y+s.height/2), float32(s.height*2), hexToRGBA("E63946"), true)

	// Draw text under the slider
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s: %.2f", s.label, s.value), s.x, s.y+20)
}
