package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type body struct {
	pos    pixel.Vec
	radius float64
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Gravity",
		Bounds: pixel.R(0, 0, 500, 500),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	body := body{
		pos:    pixel.V(10, 10),
		radius: 10,
	}
	imd := imdraw.New(nil)
	imd.Color = colornames.White
	imd.Push(body.pos)
	imd.Circle(body.radius, 0)

	for !win.Closed() {
		win.Clear(colornames.Black)
		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
