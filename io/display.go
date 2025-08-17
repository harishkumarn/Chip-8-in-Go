package io

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width int = 64
const height int = 32

var pixelColor color.Color = color.RGBA{255, 255, 255, 255}

type Display struct {
	fb       [height][width]bool
	TimeSync chan bool
}

func (d *Display) Update() error {
	d.TimeSync <- true
	return nil
}

func (d *Display) Draw(screen *ebiten.Image) {
	for y := range height {
		for x := range width {
			if d.fb[y][x] {
				vector.DrawFilledRect(
					screen,
					float32(x*10),
					float32(y*10),
					10,
					10,
					pixelColor,
					false,
				)
			}
		}
	}
}

func (d *Display) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func (d *Display) CLS() {
	for y := range height {
		for x := range width {
			d.fb[y][x] = false
		}
	}
}

func (d *Display) DrawAtPosition(x, y uint8, data []uint8) {
	for i, val := range data {
		for j := range 8 {
			if val&(1<<j) > 0 {
				fmt.Println("Setting pixel", val)
				d.fb[int(y)+i][int(x)+7-j] = true
			}
		}
	}
}

func (d *Display) Init() {
	ebiten.SetWindowSize(64*10, 32*10)
	ebiten.SetWindowTitle("Chip-8")
	d.TimeSync = make(chan bool)
	go func() {
		if err := ebiten.RunGame(d); err != nil {
			fmt.Println("Error initing the Display")
		}
	}()
}
