package main

import (
	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/rking788/activityrings"
)

func main() {

	// Draw a sample linear set of color squares to illustrate how the ring gradient will look from top to bottom
	drawSwatches()
}

func drawSwatches() {
	ctx := gg.NewContext(60, 240)

	bases := [3]colorful.Color{activityrings.BlueStart, activityrings.GreenStart, activityrings.RedStart}
	ends := [3]colorful.Color{activityrings.BlueEnd, activityrings.GreenEnd, activityrings.RedEnd}

	for column := 0; column < 3; column++ {
		activeColor := bases[column]
		for row := 0; row < 12; row++ {
			drawSwatch(ctx, row, column, activeColor)

			activeColor = activityrings.NextColor(bases[column], ends[column], float64(row)/12.0)
		}
	}

	ctx.SavePNG("sample.png")
}

func drawSwatch(ctx *gg.Context, row, col int, c colorful.Color) {

	ctx.SetFillStyle(gg.NewSolidPattern(c))
	ctx.DrawRectangle(float64(col*20), float64(row*20), 20, 20)
	ctx.Fill()
}
