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
	Mass   float64
}

type Bodies []*Body

func InitializeBodies(N int, xMax, yMax, radius float64) (bodies Bodies) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < N; i++ {
		xPos := rand.Float64() * xMax
		yPos := rand.Float64() * yMax
		body := &Body{
			Pos:    pixel.V(xPos, yPos),
			Radius: radius,
			Mass:   1,
		}
		bodies = append(bodies, body)
	}
	return
}

func (bodies Bodies) UpdateVelocities(dt float64) {
	for idxA, bodyA := range bodies[:len(bodies)-1] {
		for _, bodyB := range bodies[idxA+1:] {
			diff := Difference(bodyB.Pos, bodyA.Pos)
			dist := Distance(diff)
			// Force magnitude with masses set to 1 (for more
			// efficient calculation).
			unitMassForceLen := G * bodyA.Mass * bodyB.Mass / math.Pow(dist, 2)
			// Force of B on A (with unit masses).
			unitMassForceBA := diff.Unit().Scaled(unitMassForceLen)

			// Change in velocity (with unit masses).
			unitMassVelDiffA := unitMassForceBA.Scaled(dt)
			velDiffA := unitMassVelDiffA.Scaled(bodyB.Mass)
			bodyA.Vel = bodyA.Vel.Add(velDiffA)

			unitMassVelDiffB := unitMassVelDiffA.Scaled(-1)
			velDiffB := unitMassVelDiffB.Scaled(bodyA.Mass)
			bodyB.Vel = bodyB.Vel.Add(velDiffB)
		}
	}
}

func (bodies Bodies) UpdatePositions(dt float64) {
	for _, body := range bodies {
		displacement := body.Vel.Scaled(dt)
		body.Pos = body.Pos.Add(displacement)
	}
}

func Difference(posA, posB pixel.Vec) pixel.Vec {
	return posA.Sub(posB)
}

func Distance(diff pixel.Vec) float64 {
	return diff.Len()
}
