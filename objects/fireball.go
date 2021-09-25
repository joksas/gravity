package objects

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"gonum.org/v1/plot/palette/moreland"
)

const FireballLifetime = 100

type Fireball struct {
	Pos            pixel.Vec
	Radius         float64
	Color          color.Color
	IterationsLeft int
}

type Fireballs []Fireball

func (fireballs Fireballs) Update() Fireballs {
	palette := moreland.BlackBody()
	palette.SetMin(0)
	palette.SetMax(1)

	var newFireballs Fireballs
	for idx := 0; idx < len(fireballs); idx++ {
		fireball := &fireballs[idx]
		if fireball.IterationsLeft > 0 {
			fireball.Radius *= 1 + 1/float64(FireballLifetime)

			// Intensity will decrease with inverse square.
			intensity := math.Pow(float64(fireball.IterationsLeft)/FireballLifetime, 2)
			tempColor, err := palette.At(intensity)
			if err != nil {
				panic(err)
			}
			color := pixel.ToRGBA(tempColor)
			color = color.Mul(pixel.Alpha(intensity))
			fireball.Color = color

			fireball.IterationsLeft -= 1

			newFireballs = append(newFireballs, *fireball)
		}
	}

	if len(newFireballs) < len(fireballs) {
		return newFireballs
	} else {
		return fireballs
	}
}

func CreateFireballs(mergedBodies []Body) Fireballs {
	heaviestBodyIdx := 0
	for idx, body := range mergedBodies {
		if body.Mass > mergedBodies[heaviestBodyIdx].Mass {
			heaviestBodyIdx = idx
		}
	}
	heaviestBody := mergedBodies[heaviestBodyIdx]

	var newFireballs Fireballs
	for idx, body := range mergedBodies {
		if idx == heaviestBodyIdx {
			continue
		}
		pos := FireballPos(body.Pos, heaviestBody.Pos, body.Radius, heaviestBody.Radius)
		radius := FireballRadius(body.Mass, heaviestBody.Mass, body.Vel, heaviestBody.Vel)
		newFireball := Fireball{
			Pos:            pos,
			Radius:         radius,
			IterationsLeft: FireballLifetime,
		}
		newFireballs = append(newFireballs, newFireball)
	}

	return newFireballs
}

func FireballPos(posA, posB pixel.Vec, radiusA, radiusB float64) pixel.Vec {
	term1 := posA.Scaled(radiusB)
	term2 := posA.Scaled(radiusA)
	combinedRadius := radiusA + radiusB
	posC := term1.Add(term2).Scaled(1 / combinedRadius)
	return posC
}

func FireballRadius(massA, massB float64, velA, velB pixel.Vec) float64 {
	massC := massA + massB
	velC := VelAfterCollision(massA, massB, velA, velB)
	energyBefore := KineticEnergy(massA, velA.Len()) + KineticEnergy(massB, velB.Len())
	energyAfter := KineticEnergy(massC, velC.Len())
	energyDiff := energyBefore - energyAfter
	// Will make area proportional to square root of energy.
	area := math.Sqrt(energyDiff)
	// And radius is, of course, proportional to square root of
	// area.
	radius := math.Sqrt(area)
	// The scaling parameters are tuned to taste.
	radius *= 0.25
	return radius
}

func KineticEnergy(mass, speed float64) float64 {
	return 0.5 * mass * math.Pow(speed, 2)
}
