package objects

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

const G = 100

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
			Color:  colornames.Gray,
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

func (bodies Bodies) RemoveClose(fireballs Fireballs) (Bodies, Fireballs) {
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
		var newFireballs Fireballs
		for _, mergeGroup := range mergeGroups {
			var mergedBodies []*Body

			firstBody := bodies[mergeGroup[0]]
			mergedBodies = append(mergedBodies, firstBody)

			newPos := firstBody.Pos
			newVel := firstBody.Vel
			newRadius := firstBody.Radius
			newMass := firstBody.Mass
			for _, nextBodyIdx := range mergeGroup[1:] {
				nextBody := bodies[nextBodyIdx]
				mergedBodies = append(mergedBodies, nextBody)

				newPos = PosAfterCollision(newMass, nextBody.Mass, newPos, nextBody.Pos)
				newVel = VelAfterCollision(newMass, nextBody.Mass, newVel, nextBody.Vel)
				newRadius = RadiusAfterCollision(newRadius, nextBody.Radius)
				newMass += nextBody.Mass
			}
			newBody := &Body{
				Pos:    newPos,
				Vel:    newVel,
				Radius: newRadius,
				Mass:   newMass,
				Color:  colornames.Gray,
			}
			newBodies = append(newBodies, newBody)
			fireballs := CreateFireballs(mergedBodies)
			newFireballs = append(newFireballs, fireballs...)
		}
		return newBodies, append(fireballs, newFireballs...)
	} else {
		return bodies, fireballs
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

// Conservation of area (as we operate in 2D)
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
