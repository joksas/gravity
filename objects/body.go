package objects

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/faiface/pixel"
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
	colorChooser := CreateColorChooser()
	for i := 0; i < N; i++ {
		xPos := rand.Float64() * xMax
		yPos := rand.Float64() * yMax
		color := colorChooser.Pick().(color.Color)
		body := &Body{
			Pos:    pixel.V(xPos, yPos),
			Radius: radius,
			Mass:   1,
			Color:  color,
		}
		bodies = append(bodies, body)
	}
	return
}

func (bodies Bodies) Update(dt float64, fireballs Fireballs) (Bodies, Fireballs) {
	bodies, fireballs = bodies.RemoveClose(fireballs)
	bodies.UpdateVelocities(dt)
	bodies.UpdatePositions(dt)

	return bodies, fireballs
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
		var updatedBodies Bodies
		var updatedFireballs Fireballs
		for _, mergeGroup := range mergeGroups {
			var mergedBodies []*Body

			firstBody := bodies[mergeGroup[0]]
			mergedBodies = append(mergedBodies, firstBody)

			updatedPos := firstBody.Pos
			updatedVel := firstBody.Vel
			updatedRadius := firstBody.Radius
			updatedMass := firstBody.Mass
			updatedColor := firstBody.Color
			for _, nextBodyIdx := range mergeGroup[1:] {
				nextBody := bodies[nextBodyIdx]
				mergedBodies = append(mergedBodies, nextBody)

				updatedPos = PosAfterCollision(updatedMass, nextBody.Mass, updatedPos, nextBody.Pos)
				updatedVel = VelAfterCollision(updatedMass, nextBody.Mass, updatedVel, nextBody.Vel)
				updatedRadius = RadiusAfterCollision(updatedRadius, nextBody.Radius)
				updatedMass = MassAfterCollision(updatedMass, nextBody.Mass)
				updatedColor = ColorAfterCollision(updatedMass, nextBody.Mass, updatedColor, nextBody.Color)
			}
			updatedBody := &Body{
				Pos:    updatedPos,
				Vel:    updatedVel,
				Radius: updatedRadius,
				Mass:   updatedMass,
				Color:  updatedColor,
			}
			updatedBodies = append(updatedBodies, updatedBody)
			newFireballs := CreateFireballs(mergedBodies)
			updatedFireballs = append(updatedFireballs, newFireballs...)
		}
		return updatedBodies, append(fireballs, updatedFireballs...)
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

// Conservation of mass
func MassAfterCollision(massA, massB float64) float64 {
	massC := massA + massB
	return massC
}

// Not completely sure if this is the best approach, but I adapted
// [this](https://www.youtube.com/watch?v=LKnqECcg6Gw) to our scenario, where
// we also weigh the two colors by their masses.
func ColorAfterCollision(massA, massB float64, colorA, colorB color.Color) color.Color {
	pixelColorA := pixel.ToRGBA(colorA)
	pixelColorA2 := pixelColorA.Mul(pixelColorA)
	pixelColorB := pixel.ToRGBA(colorB)
	pixelColorB2 := pixelColorB.Mul(pixelColorB)
	term1 := pixelColorA2.Scaled(massA)
	term2 := pixelColorB2.Scaled(massB)
	combinedMass := massA + massB
	squaredContents := term1.Add(term2).Scaled(1 / combinedMass)
	r := math.Sqrt(squaredContents.R)
	g := math.Sqrt(squaredContents.G)
	b := math.Sqrt(squaredContents.B)
	color := pixel.RGBA{R: r, G: g, B: b, A: 1}
	return color
}

func Difference(posA, posB pixel.Vec) pixel.Vec {
	return posA.Sub(posB)
}

func Distance(diff pixel.Vec) float64 {
	return diff.Len()
}
