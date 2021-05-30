package objects

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/mroth/weightedrand"
)

func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func CreateColorChooser() *weightedrand.Chooser {
	// [Colors of solar system planets](https://astronomy.stackexchange.com/a/14040)
	colors := []color.Color{
		color.RGBA{26, 26, 26, 255},
		color.RGBA{230, 230, 230, 255},
		color.RGBA{47, 106, 105, 255},
		color.RGBA{153, 61, 0, 255},
		color.RGBA{176, 127, 23, 255},
		color.RGBA{176, 143, 54, 255},
		color.RGBA{85, 128, 170, 255},
		color.RGBA{54, 104, 150, 255},
	}
	rand.Seed(time.Now().UTC().UnixNano())
	var colorChoices []weightedrand.Choice
	for _, color := range colors {
		// We use exponential weights so that certain colors
		// dominate, and we don't get a boring gray mass in the
		// end. This does still happen sometimes though...
		weight := uint(math.Pow(2, float64(rand.Intn(20))))
		colorChoices = append(colorChoices, weightedrand.Choice{
			Item:   color,
			Weight: weight,
		})
	}
	colorChooser, _ := weightedrand.NewChooser(colorChoices...)

	return colorChooser
}
