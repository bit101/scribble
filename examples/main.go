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
	p := scribble.NewPen(0, 0)

	p.Line(context, 80, 100, 80, 700, 50000)
	p.Ellipse(context, 300, 400, 120, 300, 100000)
	p.Circle(context, 600, 200, 100, 50000)
	p.Rectangle(context, 500, 400, 200, 300, 50000)

	path := geom.NewPointList()
	path.AddXY(800, 100)
	path.AddXY(900, 300)
	path.AddXY(800, 500)
	path.AddXY(900, 700)
	p.Path(context, path, false, 50000)

	p.MoveTo(300, 400)
	for range 5000 {
		p.MoveTowards(300, 400)
		p.Update(context)
	}
	context.Stroke()
}
