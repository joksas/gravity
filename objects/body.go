package objects

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"gonum.org/v1/plot/palette/moreland"
)

const G = 30

type Body struct {
	Pos    pixel.Vec
	Vel    pixel.Vec
	Radius float64
	Mass   float64
	Color  color.Color
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

func (bodies Bodies) UpdateColors() {
	palette := moreland.BlackBody()
	palette.SetMin(0)
	palette.SetMax(1)

	for _, body := range bodies {
		speed := body.Vel.Len()
		// Something to determine which color to pick.
		fastness := speed / 20
		if fastness > 1 {
			fastness = 1
		}
		if fastness < 0.1 {
			fastness = 0.1
		}
		var err error
		body.Color, err = palette.At(fastness)
		if err != nil {
			panic(err)
		}
	}
}

func (bodies Bodies) RemoveClose() Bodies {
	var mergeGroups [][]int
	var mergedIdxs []int
	for idxA, bodyA := range bodies {
		// There might be some second-order merges, but, for
		// simplicity, we won't handle them. The merges will
		// simply occur in the next iteration.
		if Contains(mergedIdxs, idxA) {
			continue
		}
		mergeGroups = append(mergeGroups, []int{idxA})
		for idxB := idxA + 1; idxB < len(bodies); idxB++ {
			bodyB := bodies[idxB]
			if Contains(mergedIdxs, idxB) {
				continue
			}
			diff := Difference(bodyB.Pos, bodyA.Pos)
			dist := Distance(diff)
			if dist < bodyA.Radius || dist < bodyB.Radius {
				mergeGroups[len(mergeGroups)-1] = append(mergeGroups[len(mergeGroups)-1], idxB)
			}
		}
		mergedIdxs = append(mergedIdxs, mergeGroups[len(mergeGroups)-1]...)
	}

	if len(mergeGroups) < len(bodies) {
		var newBodies Bodies
		for _, mergeGroup := range mergeGroups {
			firstBody := bodies[mergeGroup[0]]
			mergedPos := firstBody.Pos
			mergedVel := firstBody.Vel
			mergedRadius := firstBody.Radius
			mergedMass := firstBody.Mass
			for _, newBodyIdx := range mergeGroup[1:] {
				newBody := bodies[newBodyIdx]
				mergedPos = PosAfterCollision(mergedMass, newBody.Mass, mergedPos, newBody.Pos)
				mergedVel = VelAfterCollision(mergedMass, newBody.Mass, mergedVel, newBody.Vel)
				mergedRadius = RadiusAfterCollision(mergedRadius, newBody.Radius)
				mergedMass += newBody.Mass
			}
			newBodies = append(newBodies, &Body{
				Pos:    mergedPos,
				Vel:    mergedVel,
				Radius: mergedRadius,
				Mass:   mergedMass,
			})
		}
		return newBodies
	} else {
		return bodies
	}
}

// Conservation of momentum
func VelAfterCollision(massA, massB float64, velA, velB pixel.Vec) pixel.Vec {
	momentumA := velA.Scaled(massA)
	momentumB := velB.Scaled(massB)
	massC := massA + massB
	velC := momentumA.Add(momentumB).Scaled(1 / massC)
	return velC
}

// Center of mass
func PosAfterCollision(massA, massB float64, posA, posB pixel.Vec) pixel.Vec {
	massDistanceA := posA.Scaled(massA)
	massDistanceB := posB.Scaled(massB)
	massC := massA + massB
	posC := massDistanceA.Add(massDistanceB).Scaled(1 / massC)
	return posC
}

// Conservation of radius (as we operate in 2D)
func RadiusAfterCollision(radiusA, radiusB float64) float64 {
	radiusC := math.Sqrt(math.Pow(radiusA, 2) + math.Pow(radiusB, 2))

	return radiusC
}

func Difference(posA, posB pixel.Vec) pixel.Vec {
	return posA.Sub(posB)
}

func Distance(diff pixel.Vec) float64 {
	return diff.Len()
}
