package main

import (
	"bytes"
	"flag"
	"image/color"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rking788/activityrings"
)

const (
	totalImageSize = 782
	lineWidth      = 86.0
)

func main() {

	outPath := flag.String("out-path", "rings.png", "The full path at which the activity rings image will be written including the filename")
	standValue := flag.Float64("stand", 0.0, "The decimal value to use for the stand ring")
	moveValue := flag.Float64("move", 0.0, "The decimal value to use for the move ring")
	exerciseValue := flag.Float64("exercise", 0.0, "The decimal value to use for the exercise ring")
	useHTTP := flag.Bool("http", false, "Passing this flag will start an http server for serving up activity ring images")
	flag.Parse()

	if *useHTTP {
		router := gin.Default()
		router.GET("/rings", ringsHandler)
		router.Run(":8082")
	}

	img := activityrings.NewActivityRingsImage(totalImageSize, color.Black)
	img.DrawActivity(map[activityrings.ActivityType]float64{activityrings.Stand: *standValue, activityrings.Exercise: *exerciseValue, activityrings.Move: *moveValue})
	img.SavePNG(*outPath)
}

func ringsHandler(ctx *gin.Context) {

	buf := &bytes.Buffer{}

	standValue, err := strconv.ParseFloat(ctx.Query("stand"), 64)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Bad request, could not parse stand value as a float")
		return
	}

	exerciseValue, err := strconv.ParseFloat(ctx.Query("exercise"), 64)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Bad request, could not parse exercise value as a float")
		return
	}

	moveValue, err := strconv.ParseFloat(ctx.Query("move"), 64)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Bad request, could not parse move value as a float")
		return
	}

	img := activityrings.NewActivityRingsImage(totalImageSize, color.Transparent)
	img.DrawActivity(map[activityrings.ActivityType]float64{activityrings.Stand: standValue, activityrings.Exercise: exerciseValue, activityrings.Move: moveValue})
	img.EncodePNG(buf)

	ctx.DataFromReader(http.StatusOK, int64(buf.Len()), "image/png", buf, nil)
}
