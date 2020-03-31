package activityrings

import (
	"math"

	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
)

type gradientDirection int

const (
	topDownDirection gradientDirection = iota
	bottomUpDirection
)

func (d gradientDirection) flip() gradientDirection {
	if d == topDownDirection {
		return bottomUpDirection
	}

	return topDownDirection
}

// NextColor will calculate the next intperolated color between start and end given the decimal value, t. For 0.0, start color is
// used, for a value of 1.0, end is used. This method should probably not be called directly, the activity rings image type will call this as needed.
func NextColor(start, end colorful.Color, t float64) colorful.Color {
	next := interpolateTo(start, end, t)
	return next
}

func newRingGradient(center float64, direction gradientDirection, startColor, stopColor colorful.Color, ringRadius float64) gg.Pattern {

	startY := 0.0
	endY := 0.0
	if direction == topDownDirection {
		startY = center - ringRadius + (lineWidth / 2.0)
		endY = center + ringRadius - (lineWidth / 2.0)
	} else {
		startY = center + ringRadius - (lineWidth / 2.0)
		endY = center - ringRadius + (lineWidth / 2.0)
	}
	grad := gg.NewLinearGradient(center, startY, center, endY)
	grad.AddColorStop(0.0, startColor)
	grad.AddColorStop(1.0, stopColor)

	return grad
}

func interpolateTo(start, end colorful.Color, t float64) colorful.Color {

	r1, g1, b1 := start.RGB255()
	r2, g2, b2 := end.RGB255()
	r := lerp(t, r1, r2)
	g := lerp(t, g1, g2)
	b := lerp(t, b1, b2)

	return colorful.Color{R: float64(r) / 255.0, G: float64(g) / 255.0, B: float64(b) / 255.0}
}

func lerp(t float64, a, b uint8) uint8 {
	return uint8(float64(a) + (math.Min(math.Max(t, 0), 1) * (float64(b) - float64(a))))
}
