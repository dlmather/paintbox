package painter

import (
	"github.com/dlmather/paintbox/display"
	"github.com/nsf/termbox-go"
)

const (
	prevUnset  = -1
	colorLimit = 8
)

// Stores state for on-screen cursor and its operations
type Cursor struct {
	Color, PrevColor int
	X, Y             int
	PrevX, PrevY     int
}

func NewCursor() *Cursor {
	return &Cursor{
		Color: 1,
		PrevX: prevUnset,
		PrevY: prevUnset,
	}
}

func (c *Cursor) PrevUnset() bool {
	return c.PrevX == prevUnset || c.PrevY == prevUnset
}

type Painter struct {
	Canvas        *Canvas
	CanvasHistory *CanvasStack
	Cursor        *Cursor
	Display       *display.Display
}

func New() *Painter {
	w, h := termbox.Size()
	p := &Painter{
		Canvas:        NewCanvas(w, h-3),
		CanvasHistory: NewCanvasStack(),
		Cursor:        NewCursor(),
		Display:       display.New(w, h),
	}
	p.Display.Draw()
	return p
}

func (p *Painter) Draw() {
	termbox.SetCursor(p.Cursor.X, p.Cursor.Y)
	p.Canvas.apply()
	termbox.Flush()
}

func MoveLeft(p *Painter) error {
	if p.Cursor.X == 0 {
		// Resize!
	}
	p.Cursor.X -= 1
	return nil
}

func MoveRight(p *Painter) error {
	if p.Cursor.X == p.Canvas.Width {
		// Resize!
	}
	p.Cursor.X += 1
	return nil
}

func MoveUp(p *Painter) error {
	if p.Cursor.Y == 0 {
		// Resize!
	}
	p.Cursor.Y -= 1
	return nil
}

func MoveDown(p *Painter) error {
	if p.Cursor.Y == p.Canvas.Height {
		// Resize!
	}
	p.Cursor.Y += 1
	return nil
}

func CycleColor(p *Painter) error {
	p.Cursor.Color, p.Cursor.PrevColor = (p.Cursor.Color%colorLimit)+1, p.Cursor.Color
	p.Display.CurrentColor = p.Cursor.Color
	p.Display.Draw()
	return nil
}

func Point(p *Painter) error {
	x, y, color := p.Cursor.X, p.Cursor.Y, p.Cursor.Color
	p.Canvas.Point(x, y, color)
	return nil
}

func Erase(p *Painter) error {
	x, y := p.Cursor.X, p.Cursor.Y
	p.Canvas.Erase(x, y)
	return nil
}

func Line(p *Painter) error {
	x, y, color := p.Cursor.X, p.Cursor.Y, p.Cursor.Color
	prevX, prevY := p.Cursor.PrevX, p.Cursor.PrevY
	if prevX == prevUnset || prevY == prevUnset {
		p.Canvas.Point(x, y, color)
		p.Cursor.PrevX, p.Cursor.PrevY = x, y
	} else {
		p.Canvas.BresenhamLine(x, y, prevX, prevY, color)
		p.Cursor.PrevX, p.Cursor.PrevY = prevUnset, prevUnset
	}
	return nil
}

func FloodFill(p *Painter) error {
	x, y, color := p.Cursor.X, p.Cursor.Y, p.Cursor.Color
	replaceColor := p.Canvas.Squares[x][y]
	p.Canvas.FloodFill(x, y, replaceColor, color)
	return nil
}

func Box(p *Painter) error {
	x, y, color := p.Cursor.X, p.Cursor.Y, p.Cursor.Color
	prevX, prevY := p.Cursor.PrevX, p.Cursor.PrevY
	if prevX == prevUnset || prevY == prevUnset {
		p.Canvas.Point(x, y, color)
		p.Cursor.PrevX, p.Cursor.PrevY = x, y
	} else {
		p.Canvas.Box(x, y, prevX, prevY, color)
		p.Cursor.PrevX, p.Cursor.PrevY = prevUnset, prevUnset
	}
	return nil
}

func FullBox(p *Painter) error {
	x, y, color := p.Cursor.X, p.Cursor.Y, p.Cursor.Color
	prevX, prevY := p.Cursor.PrevX, p.Cursor.PrevY
	if prevX == prevUnset || prevY == prevUnset {
		p.Canvas.Point(x, y, color)
		p.Cursor.PrevX, p.Cursor.PrevY = x, y
	} else {
		p.Canvas.FullBox(x, y, prevX, prevY, color)
		p.Cursor.PrevX, p.Cursor.PrevY = prevUnset, prevUnset
	}
	return nil
}
