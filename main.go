package main

import (
	"fmt"
	"os"

	"github.com/ajstarks/svgo"
)

type RectStyle struct {
	Stroke      string
	StrokeWidth int
	Fill        string
}

func (s *RectStyle) String() string {
	return fmt.Sprintf("stroke-opacity:0.4;stroke:%s;stroke-width:%d;fill:%s;", s.Stroke, s.StrokeWidth, s.Fill)
}
func (s *RectStyle) SetStroke(stroke string) *RectStyle {
	s.Stroke = stroke
	return s
}

func (s *RectStyle) SetStrokeWidth(strokeWidth int) *RectStyle {
	s.StrokeWidth = strokeWidth
	return s
}

func (s *RectStyle) SetFill(fill string) *RectStyle {
	s.Fill = fill
	return s
}

var (
	DefaultRectStyle = RectStyle{
		Stroke:      "black",
		StrokeWidth: 1,
		Fill:        "white",
	}
)

func DefaultStyleRect() *RectStyle {
	tmp := DefaultRectStyle
	return &tmp
}

type Pos struct {
	X, Y int
}

type Grid struct {
	Pos    Pos
	W, H   int
	WN, HN int
	Margin int
	Style  *RectStyle
}

func (g *Grid) PosAt(x, y int) (int, int) {
	return g.Pos.X + g.W*x,
		g.Pos.Y + g.H*y
}

func (g *Grid) ChildGrid(x, y, wn, hn int) *Grid {
	sx, sy := g.PosAt(x, y)
	childGrid := &Grid{
		Pos: Pos{sx, sy},
		W:   g.W, H: g.H,
		WN: wn, HN: hn,
		Margin: g.Margin + 1,
		Style:  DefaultStyleRect().SetStroke("green").SetFill("transparent").SetStrokeWidth(2),
	}
	return childGrid
}

func (g *Grid) ChangeGridRate(rw, rh int) {
	g.W *= rw
	g.H *= rh
	g.WN /= rw
	g.HN /= rh
}

func (g *Grid) Draw(canvas *svg.SVG) {
	for j := 0; j < g.HN; j++ {
		for i := 0; i < g.WN; i++ {
			x, y := g.PosAt(i, j)
			lastLineMargin := 0
			if j == g.HN-1 {
				// NOTE: last line
				lastLineMargin = g.Margin * 2
			}
			canvas.Rect(x+g.Margin, y+g.Margin, g.W-g.Margin, g.H-g.Margin-lastLineMargin, g.Style.String())
		}
	}
	if g.WN == 0 {
		x, y := g.PosAt(0, 0)
		canvas.Rect(x, y, 1, g.H*g.HN, g.Style.String())
	}
	if g.HN == 0 {
		x, y := g.PosAt(0, 0)
		canvas.Rect(x, y, g.W*g.WN, 1, g.Style.String())
	}
}

func main() {
	width := 1024
	height := 800
	canvas := svg.New(os.Stdout)
	canvas.Start(width, height)

	// canvas.Circle(width/2, height/2, 100)
	// canvas.Text(width/2, height/2, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:white")

	grid := Grid{
		Pos: Pos{10, 10},
		W:   16, H: 16,
		WN: 48, HN: 24,
		Style: DefaultStyleRect().SetStroke("gray"),
	}
	grid.Draw(canvas)

	cnt := 0
	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			tx, ty := grid.PosAt(1+16*i, 1+j)
			canvas.Text(tx+8, ty+14, fmt.Sprint(cnt), "text-anchor:middle;font-size:16px;fill:black")
			cnt++
		}
	}

	for i := 0; i < 1; i++ {
		// NOTE: up
		{
			childGrid := grid.ChildGrid(1+16*i, 0, 16, 1)
			childGrid.ChangeGridRate(16, 1)
			childGrid.Style.SetStrokeWidth(4)
			childGrid.Style.SetStroke("blue")
			childGrid.Draw(canvas)
		}
		// NOTE: down
		{
			childGrid := grid.ChildGrid(1+16*i, 2, 16, 1)
			childGrid.ChangeGridRate(16, 1)
			childGrid.Style.SetStrokeWidth(4)
			childGrid.Style.SetStroke("blue")
			childGrid.Draw(canvas)
		}
		// NOTE: right
		{
			childGrid := grid.ChildGrid(2+16*i, 1, 16, 1)
			childGrid.ChangeGridRate(16, 1)
			childGrid.Style.SetStrokeWidth(4)
			childGrid.Style.SetStroke("purple")
			childGrid.Draw(canvas)
		}
		// NOTE: left
		{
			childGrid := grid.ChildGrid(0+16*i, 1, 16, 1)
			childGrid.ChangeGridRate(16, 1)
			childGrid.Style.SetStrokeWidth(4)
			childGrid.Style.SetStroke("purple")
			childGrid.Draw(canvas)
		}
	}
	// NOTE: alignment line
	for i := 0; i < 3; i++ {
		childGrid := grid.ChildGrid(1+16*i, 0, 0, 16)
		childGrid.Style.SetStroke("red")
		childGrid.Style.SetStrokeWidth(8)
		childGrid.Draw(canvas)
	}
	canvas.End()
}
