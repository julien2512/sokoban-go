package main

import (
	"time"

	"github.com/TheInvader360/sokoban-go/controller"
	"github.com/TheInvader360/sokoban-go/model"
	"github.com/TheInvader360/sokoban-go/view"

	pixelgl "github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

const (
	width       = 512
	height      = 256
	scaleFactor = 3
)

func run() {
	cfg := opengl.WindowConfig{
		Title:  "Sokoban",
		Bounds: pixelgl.R(0, 0, width*scaleFactor, height*scaleFactor),
		VSync:  true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	m := model.NewModel()
	v := view.NewView(m, win, scaleFactor)
	c := controller.NewController(m)
	lastKey := pixelgl.UnknownButton
	c.StartNewGame()

	// Main game loop
	for !win.Closed() {
		if win.Pressed(pixelgl.KeyEscape) {
			return
		}

		// Fire an event once per key press (no repeats if the key is held down)
		// Note: JustPressed() is a cleaner way to achieve this, but Pressed() more closely matches the Jack OS API
		if win.Pressed(pixelgl.KeyUp) {
			if lastKey != pixelgl.KeyUp {
				c.HandleInput(pixelgl.KeyUp)
			}
			lastKey = pixelgl.KeyUp
		} else if win.Pressed(pixelgl.KeyDown) {
			if lastKey != pixelgl.KeyDown {
				c.HandleInput(pixelgl.KeyDown)
			}
			lastKey = pixelgl.KeyDown
		} else if win.Pressed(pixelgl.KeyLeft) {
			if lastKey != pixelgl.KeyLeft {
				c.HandleInput(pixelgl.KeyLeft)
			}
			lastKey = pixelgl.KeyLeft
		} else if win.Pressed(pixelgl.KeyRight) {
			if lastKey != pixelgl.KeyRight {
				c.HandleInput(pixelgl.KeyRight)
			}
			lastKey = pixelgl.KeyRight
		} else if win.Typed() == "z" {
			if lastKey != pixelgl.KeyZ {
				c.HandleInput(pixelgl.KeyZ)
			}
			lastKey = pixelgl.KeyZ
		} else if win.Typed() == "f" {
			if lastKey != pixelgl.KeyF {
				c.HandleInput(pixelgl.KeyF)
			}
			lastKey = pixelgl.KeyF
		} else if win.Typed() == "r" {
			if lastKey != pixelgl.KeyR {
				c.HandleInput(pixelgl.KeyR)
			}
			lastKey = pixelgl.KeyR
		} else if win.Typed() == "a" {
			if lastKey != pixelgl.KeyA {
				c.HandleInput(pixelgl.KeyA)
			}
			lastKey = pixelgl.KeyR
		} else if win.Pressed(pixelgl.KeySpace) {
			if lastKey != pixelgl.KeySpace {
				c.HandleInput(pixelgl.KeySpace)
			}
			lastKey = pixelgl.KeySpace
		} else {
			lastKey = pixelgl.UnknownButton
		}

		m.Update()

		v.Draw(c.ShowFreeSpace)

		time.Sleep(50 * time.Millisecond)
	}
}

func main() {
	opengl.Run(run)
}
