package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/joksas/gravity/objects"
	"golang.org/x/image/colornames"
)

func run() {
	screenWidth := 1000.0
	screenHeight := 700.0
	cfg := pixelgl.WindowConfig{
		Title:  "Gravity",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	bodies := objects.InitializeBodies(300, screenWidth, screenHeight, 1)
	imd := imdraw.New(nil)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		bodies = bodies.RemoveClose()
		bodies.UpdateVelocities(dt)
		bodies.UpdatePositions(dt)
		bodies.UpdateColors()

		imd.Clear()
		for _, body := range bodies {
			imd.Color = body.Color
			imd.Push(body.Pos)
			imd.Circle(body.Radius, 0)
		}

		win.Clear(colornames.Black)
		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
