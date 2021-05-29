package objects

import (
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
)

const G = 1

type Body struct {
	Pos    pixel.Vec
	Vel    pixel.Vec
	Radius float64
}

func InitializeBodies(N int, xMax, yMax, radius float64) (bodies []*Body) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < N; i++ {
		xPos := rand.Float64() * xMax
		yPos := rand.Float64() * yMax
		body := &Body{
			Pos:    pixel.V(xPos, yPos),
			Radius: radius,
		}
		bodies = append(bodies, body)
	}
	return
}

func UpdateVelocities(bodies []*Body, dt float64) {
	for idxA, bodyA := range bodies[:len(bodies)-1] {
		for _, bodyB := range bodies[idxA+1:] {
			diff := Difference(bodyB.Pos, bodyA.Pos)
			dist := Distance(diff)
			forceLen := G / math.Pow(dist, 2)
			// Force of B on A.
			forceBA := diff.Unit().Scaled(forceLen)

			velDiffA := forceBA.Scaled(dt)
			bodyA.Vel = bodyA.Vel.Add(velDiffA)

			velDiffB := velDiffA.Scaled(-1)
			bodyB.Vel = bodyB.Vel.Add(velDiffB)
		}
	}
}

func Difference(posA, posB pixel.Vec) pixel.Vec {
	return posA.Sub(posB)
}

func Distance(diff pixel.Vec) float64 {
	return diff.Len()
}
