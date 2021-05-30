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
	var fireballs objects.Fireballs
	imd := imdraw.New(nil)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		dt /= 5

		bodies, fireballs = bodies.RemoveClose(fireballs)
		bodies.UpdateVelocities(dt)
		bodies.UpdatePositions(dt)
		fireballs = fireballs.Update()

		imd.Clear()
		for _, body := range bodies {
			imd.Color = body.Color
			imd.Push(body.Pos)
			imd.Circle(body.Radius, 0)
		}

		for _, fireball := range fireballs {
			imd.Color = fireball.Color
			imd.Push(fireball.Pos)
			imd.Circle(fireball.Radius, 0)
		}

		win.Clear(colornames.Black)
		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
