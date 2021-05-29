package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/joksas/gravity/objects"
	"golang.org/x/image/colornames"
)

type Body struct {
	Pos    pixel.Vec
	Radius float64
}

func run() {
	screenWidth := 500.0
	screenHeight := 500.0
	cfg := pixelgl.WindowConfig{
		Title:  "Gravity",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	bodies := objects.InitializeBodies(100, screenWidth, screenHeight, 5)
	imd := imdraw.New(nil)
	imd.Color = colornames.White
	for _, body := range bodies {
		imd.Push(body.Pos)
		imd.Circle(body.Radius, 0)
	}

	for !win.Closed() {
		win.Clear(colornames.Black)
		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
