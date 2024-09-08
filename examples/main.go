// Package main renders an image or video
package main

import (
	"math"

	"github.com/bit101/bitlib/blmath"
	"github.com/bit101/bitlib/geom"
	cairo "github.com/bit101/blcairo"
	"github.com/bit101/blcairo/render"
	"github.com/bit101/scribble"
)

//revive:disable:unused-parameter
const (
	tau = blmath.Tau
	pi  = math.Pi
)

func main() {
	fileName := "scribble"
	render.CreateAndViewImage(1000, 800, "out/"+fileName+".png", scene1, 0.0)
}

func scene1(context *cairo.Context, width, height, percent float64) {
	context.BlackOnWhite()
	pen := scribble.NewPen(0, 0)

	pen.Line(80, 100, 80, 700, 50000)
	pen.Ellipse(300, 400, 120, 300, 100000)
	pen.Dot(300, 400, 5000)
	pen.Circle(600, 200, 100, 50000)
	pen.Rectangle(500, 400, 200, 300, 70000)
	pen.Arc(875, 600, 100, 1, -1, false, 30000)

	path := geom.NewPointList()
	path.AddXY(800, 100)
	path.AddXY(950, 200)
	path.AddXY(800, 300)
	path.AddXY(950, 400)
	pen.Path(path, false, 50000)

	for _, list := range pen.GetPoints() {
		context.StrokePath(list, false)
	}
}
