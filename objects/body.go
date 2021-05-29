package objects

import (
	"math/rand"
	"time"

	"github.com/faiface/pixel"
)

type Body struct {
	Pos    pixel.Vec
	Radius float64
}

func InitializeBodies(N int, xMax, yMax, radius float64) (bodies []Body) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < N; i++ {
		xPos := rand.Float64() * xMax
		yPos := rand.Float64() * yMax
		body := Body{
			Pos:    pixel.V(xPos, yPos),
			Radius: radius,
		}
		bodies = append(bodies, body)
	}
	return
}
