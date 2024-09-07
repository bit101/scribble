// Package scribble holds the scribble library.
package scribble

import (
	"math"

	"github.com/bit101/bitlib/blmath"
	"github.com/bit101/bitlib/geom"
	"github.com/bit101/bitlib/random"
)

// Pen represents a pen drawing a scribbled line.
type Pen struct {
	x, y, vx, vy, vr, damp, step, curl, pull, reverse float64
	points                                            []geom.PointList
}

// NewPen creates a new pen object.
func NewPen(x, y float64) *Pen {
	pen := &Pen{
		x:       x,
		y:       y,
		vx:      0,
		vy:      0,
		vr:      0,
		damp:    0.7,
		step:    1,
		reverse: 0.5,
	}
	pen.SetPull(30)
	pen.SetCurl(30)
	pen.points = []geom.PointList{}
	return pen
}

// Position returns the pen's current x, y position.
func (p *Pen) GetPosition() (float64, float64) {
	return p.x, p.y
}

// MoveTo moves the pen to the x, y position.
func (p *Pen) MoveTo(x, y float64) {
	p.x = x
	p.y = y
	p.points = append(p.points, geom.NewPointList())
	p.addCurrentPoint()
}

// SetCurl sets how much curl the scribbles will have.
// Suggested range is 0-100 and default is 30.
func (p *Pen) SetCurl(c float64) {
	p.curl = blmath.Map(c, 0, 100, 0, math.Pi)
}

// SetDamp sets how much the drawing velocity will be dampened.
// Suggested range is 0-1. 0 is full dampening, 1 is no dampening. The default is 0.7.
func (p *Pen) SetDamp(d float64) {
	p.damp = d
}

// SetPull sets how quickly a pen is drawn to a target point.
// Suggested range is 0-100 and default is 30.
func (p *Pen) SetPull(g float64) {
	p.pull = blmath.Map(g, 0, 100, 0, 0.5)
}

// SetReverse sets the negative range of the random curl.
// The default of 0.5 gives you random values from -curl*0.5 to curl.
// 1 gives you -curl to curl.
// 0 gives you 0 to curl.
// Higher values create a more meandering scribble.
// Lower values create more tight loops.
// Suggested range is 0-1 and default is 0.5.
func (p *Pen) SetReverse(r float64) {
	p.reverse = r
}

// SetStep determines how much each pen will travel on each iteration.
// Higher values will create a larger, chunkier, less smooth curve.
// Suggested range is 0.5 to 5. Default is 1.
func (p *Pen) SetStep(s float64) {
	p.step = s
}

// Update updates and draws the path of this pen.
func (p *Pen) Update() {
	p.x += p.vx
	p.y += p.vy
	p.addCurrentPoint()

	p.vr += random.FloatRange(-p.curl*p.reverse, p.curl)
	p.vx += math.Cos(p.vr) * p.step
	p.vy += math.Sin(p.vr) * p.step
	p.vx *= p.damp
	p.vy *= p.damp
}

func (p *Pen) addCurrentPoint() {
	p.points[len(p.points)-1].AddXY(p.x, p.y)
}

// ClearAll deletes all points on all paths.
func (p *Pen) ClearAll() {
	p.points = []geom.PointList{}
}

// PullTowards pulls a pen towards this location.
// Multiple PullTowards can be called and the point will tend towards the average of these.
func (p *Pen) PullTowards(x, y float64) {
	dx := x - p.x
	dy := y - p.y
	angle := math.Atan2(dy, dx)
	p.vx += math.Cos(angle) * p.pull
	p.vy += math.Sin(angle) * p.pull
}

// GetPoints returns the list of list of points.
func (p *Pen) GetPoints() []geom.PointList {
	return p.points
}

// Line draws a single scribbled line from one point to another,
// with a specified number of iterations.
func (p *Pen) Line(x0, y0, x1, y1 float64, count int) {
	countf := float64(count)
	p.MoveTo(x0, y0)
	for i := range count {
		f := float64(i)
		x := blmath.Map(f, 0, countf, x0, x1)
		y := blmath.Map(f, 0, countf, y0, y1)
		p.PullTowards(x, y)
		p.Update()
	}
}

// Circle draws a single scribbled circle,
// with the specified center, radius and number of iterations.
func (p *Pen) Circle(xc, yc, radius float64, count int) {
	countf := float64(count)
	p.MoveTo(xc+radius, yc)
	for i := range count {
		f := float64(i)
		a := f / countf * blmath.Tau
		x := xc + math.Cos(a)*radius
		y := yc + math.Sin(a)*radius
		p.PullTowards(x, y)
		p.Update()
	}
}

// Arc draws a single scribbled arc,
// with the specified center, radius, start and end angle, cw/ccw, and number of iterations.
func (p *Pen) Arc(xc, yc, radius, start, end float64, ccw bool, count int) {
	countf := float64(count)
	p.MoveTo(xc+math.Cos(start)*radius, yc+math.Sin(start)*radius)
	start = blmath.WrapTau(start)
	end = blmath.WrapTau(end)
	for i := range count {
		f := float64(i)
		a := f / countf * (end - start)
		if ccw {
			a = start - a
		} else {
			a = start + a
		}
		x := xc + math.Cos(a)*radius
		y := yc + math.Sin(a)*radius
		p.PullTowards(x, y)
		p.Update()
	}
}

// Ellipse draws a single scribbled ellipse,
// with the specified center, radii and number of iterations.
func (p *Pen) Ellipse(xc, yc, xr, yr float64, count int) {
	countf := float64(count)
	p.MoveTo(xc+xr, yc)
	for i := range count {
		f := float64(i)
		a := f / countf * blmath.Tau
		x := xc + math.Cos(a)*xr
		y := yc + math.Sin(a)*yr
		p.PullTowards(x, y)
		p.Update()
	}
}

// Rectangle draws a single scribbled rectangle
// with the specified x, y, width, height and number of iterations.
func (p *Pen) Rectangle(x0, y0, w, h float64, count int) {
	countf := float64(count)
	p.MoveTo(x0, y0)
	xCount := countf / 2 * w / (w + h)
	yCount := countf/2 - xCount
	p.Line(x0, y0, x0+w, y0, int(xCount))
	p.Line(x0+w, y0, x0+w, y0+h, int(yCount))
	p.Line(x0, y0+h, x0+w, y0+h, int(xCount))
	p.Line(x0, y0, x0, y0+h, int(yCount))
}

// Path draws a scribbled path between a list of points.
func (p *Pen) Path(points geom.PointList, closed bool, count int) {
	countf := float64(count)
	length := points.Length()
	if closed {
		length += points.First().Distance(points.Last())
	}
	p.MoveTo(points.First().X, points.First().Y)
	for i := range len(points) - 1 {
		p0 := points[i]
		p1 := points[i+1]
		dist := p0.Distance(p1)
		p.Line(p0.X, p0.Y, p1.X, p1.Y, int(countf/length*dist))
	}
	if closed {
		dist := points.First().Distance(points.Last())
		p.Line(points.First().X, points.First().Y, points.Last().X, points.Last().Y, int(countf/length*dist))
	}
}

// Dot scribbles a dot at the specified location.
func (p *Pen) Dot(x, y float64, count int) {
	p.MoveTo(x, y)
	for range count {
		p.PullTowards(x, y)
		p.Update()
	}
}

// TODO:
// beziers
