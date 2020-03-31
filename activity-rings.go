package activityrings

import (
	"image/color"
	"io"
	"math"

	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
)

// ActivityType is a type that represents one of the rings
type ActivityType int

const (
	// Stand is the activity type for the Apple Watch stand goal
	Stand ActivityType = iota
	// Exercise is the activity type for the exercise goal on the Apple Watch Activity app
	Exercise
	// Move is the activity type for the Move goal on the Apple Watch
	Move
)

const (
	// TODO: All of these values should be a function of the actual image size. right now they are appropriate for the size in the example cmd/activity/main.go
	lineWidth   = 86.0
	ringPadding = 6.0
	innerCircle = 115.0
)

// Inactive ring colors
var (
	InactiveRed   = color.RGBA{30, 1, 3, 255}
	InactiveGreen = color.RGBA{14, 32, 3, 255}
	InactiveBlue  = color.RGBA{6, 27, 33, 255}
)

// Active starting colors
var (
	// Credit to the MKRingProgressView Example project
	RedStart   = colorful.Color{R: 0.8823529412, G: 0, B: 0.07843137255}
	RedEnd     = colorful.Color{R: 1, G: 0.1960784314, B: 0.5294117647}
	GreenStart = colorful.Color{R: 0.2156862745, G: 0.862745098, B: 0}
	GreenEnd   = colorful.Color{R: 0.7176470588, G: 1, B: 0}
	BlueStart  = colorful.Color{R: 0, G: 0.7294117647, B: 0.8823529412}
	BlueEnd    = colorful.Color{R: 0, G: 0.9803921569, B: 0.8156862745}
)

// ActivityRing represents the data needed to draw a single activity ring including the type of activity
// it represents and the colors needed to do the drawing.
type ActivityRing struct {
	radius        float64
	inactiveColor color.RGBA
	startColor    colorful.Color
	endColor      colorful.Color

	ActivityType
}

// ActivityRingsImage is a representation of a collection of activity rings that can be drawn.
type ActivityRingsImage struct {
	ctx               *gg.Context
	backgroundColor   color.Color
	center            float64
	lineWidth         float64
	ringPadding       float64
	innerCircleRadius float64
	rings             map[ActivityType]*ActivityRing
}

// NewActivityRingsImage is a convenience method for creating a new ActivityRingsImage with the provided background color and full image size.
// The drawn image will be a square using the imgSize parameter.
func NewActivityRingsImage(imgSize int, bgColor color.Color) *ActivityRingsImage {

	radius := innerCircle + (lineWidth / 2.0)
	standRing := &ActivityRing{radius: radius, inactiveColor: InactiveBlue, startColor: BlueStart, endColor: BlueEnd, ActivityType: Stand}

	radius += lineWidth + ringPadding
	exerciseRing := &ActivityRing{radius: radius, inactiveColor: InactiveGreen, startColor: GreenStart, endColor: GreenEnd, ActivityType: Exercise}

	radius += lineWidth + ringPadding
	moveRing := &ActivityRing{radius: radius, inactiveColor: InactiveRed, startColor: RedStart, endColor: RedEnd, ActivityType: Move}

	context := gg.NewContext(imgSize, imgSize)
	context.SetLineCapRound()
	context.SetLineWidth(lineWidth)

	img := &ActivityRingsImage{
		ctx:               context,
		center:            float64(imgSize) / 2.0,
		backgroundColor:   bgColor,
		lineWidth:         lineWidth,
		ringPadding:       ringPadding,
		innerCircleRadius: innerCircle,
		rings:             map[ActivityType]*ActivityRing{Stand: standRing, Exercise: exerciseRing, Move: moveRing},
	}
	img.drawEmptyActivityRings()

	return img
}

func (img *ActivityRingsImage) drawEmptyActivityRings() {
	img.ctx.SetColor(img.backgroundColor)
	img.ctx.DrawRectangle(0, 0, float64(img.ctx.Width()), float64(img.ctx.Height()))
	img.ctx.Fill()

	for _, ring := range img.rings {

		img.ctx.SetStrokeStyle(gg.NewSolidPattern(ring.inactiveColor))
		img.ctx.DrawArc(img.center, img.center, ring.radius, 0, 2*math.Pi)
		img.ctx.Stroke()
	}
}

// DrawActivity is the main method for drawing a set of activity values. If an unrecognized activity type is provided, it is ignored. The activity
// values should represent a fraction of the activity goal completed. A value of 1.0 means the activity goal is met exactly, 2.0 means the activity is double the goal.
func (img *ActivityRingsImage) DrawActivity(values map[ActivityType]float64) {
	// TODO: This parameter could be expanded to allow an arbitrary number of rings? or possibly a mapping of ActivityTypes to values.

	for t, v := range values {
		if ring, ok := img.rings[t]; ok {
			img.drawActivityValue(ring, v)
		}
	}
}

func (img *ActivityRingsImage) drawActivityValue(ring *ActivityRing, value float64) {

	// TODO: This probably needs a special case for value == 0.0, should probably just draw a circle at the top with the active starting color fill

	startValue := value
	direction := topDownDirection
	startColor := ring.startColor
	needsShadow := false
	if value >= 1.0 {
		needsShadow = true
	}

	currentValue := 0.0
	for {
		if value <= 0.0 {
			break
		}
		startRadians := -0.5 * math.Pi
		if direction == bottomUpDirection {
			startRadians = 0.5 * math.Pi
		}
		if needsShadow && value <= 0.5 {
			angle := (startValue * (2.0 * math.Pi)) - (0.5 * math.Pi) + (0.01 * math.Pi)
			img.drawShadow(angle, ring.radius)
		}

		// Draw either the current value or a full half circle segment, which ever is smaller
		segmentValue := math.Min(value, 0.5)
		endRadians := ((segmentValue / 0.5) * math.Pi) + startRadians

		currentValue += segmentValue
		stopColor := NextColor(ring.startColor, ring.endColor, currentValue/startValue)

		grad := newRingGradient(img.center, direction, startColor, stopColor, ring.radius)
		img.ctx.SetStrokeStyle(grad)
		img.ctx.DrawArc(img.center, img.center, ring.radius, startRadians, endRadians)
		img.ctx.Stroke()

		value -= 0.5
		direction = direction.flip()
		startColor = stopColor
	}
}

func (img *ActivityRingsImage) drawShadow(angle, radius float64) {
	x, y := arcEnd(angle, radius, img.center)

	grad := gg.NewRadialGradient(x, y, (img.lineWidth/2.0)-11, x, y, (img.lineWidth / 2.0))
	grad.AddColorStop(0.0, color.RGBA{0, 0, 0, 255})
	grad.AddColorStop(1.0, color.RGBA{0, 0, 0, 10})

	img.ctx.SetFillStyle(grad)
	img.ctx.DrawCircle(x, y, (img.lineWidth / 2.0))
	img.ctx.Fill()
}

func arcEnd(angle, radius, center float64) (float64, float64) {
	x := center + radius*math.Cos(angle)
	y := center + radius*math.Sin(angle)
	return x, y
}

// SavePNG will write the activity rings image to the provided path, returning an error if the write operation fails.
func (img *ActivityRingsImage) SavePNG(path string) error {
	return img.ctx.SavePNG(path)
}

// EncodePNG will write the rings activity image to the provided Writer as a PNG.
func (img *ActivityRingsImage) EncodePNG(w io.Writer) error {
	return img.ctx.EncodePNG(w)
}
