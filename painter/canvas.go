package painter

import (
	"errors"

	"github.com/nsf/termbox-go"
)

const (
	maxDimension      = 1024
	dimensionIncrease = 64
)

type Update struct {
	Color int
	X, Y  int
}

// A representation of the actual state of the canvas
type Canvas struct {
	Width, Height int
	Squares       [][]int
	Updates       []Update
}

func NewCanvas(w, h int) *Canvas {
	can := Canvas{
		Width:   w,
		Height:  h,
		Updates: make([]Update, 0),
	}
	can.Squares = make([][]int, can.Width)
	for i := 0; i < can.Width; i++ {
		can.Squares[i] = make([]int, can.Height)
	}
	return &can
}

func (can *Canvas) push(x, y, color int) {
	update := Update{
		Color: color,
		X:     x,
		Y:     y,
	}
	can.Updates = append(can.Updates, update)
}

func (can *Canvas) apply() {
	for _, update := range can.Updates {
		x, y, color := update.X, update.Y, update.Color
		can.Squares[x][y] = color
		termbox.SetCell(x, y, ' ', termbox.Attribute(color), termbox.Attribute(color))
	}
	can.Updates = make([]Update, 0)
}

func (can *Canvas) Resize(w, h int) (*Canvas, error) {
	newWidth := can.Width + w
	newHeight := can.Height + h
	if newHeight > maxDimension || newHeight > maxDimension {
		return can, errors.New("max canvas size exceeded")
	}
	newCanvas := NewCanvas(newWidth, newHeight)
	for j, row := range can.Squares {
		for i, color := range row {
			newCanvas.push(i, j, color)
		}
	}
	newCanvas.apply()
	return newCanvas, nil
}

func (can *Canvas) Point(x, y, color int) {
	can.push(x, y, color)
}

func (can *Canvas) Erase(x, y int) {
	can.push(x, y, 0)
}

// http://en.wikipedia.org/wiki/Flood_fill
func (can *Canvas) FloodFill(x, y, targetColor, replaceColor int) {
	if targetColor == replaceColor {
		return
	}
	if can.Squares[x][y] != targetColor {
		return
	}
	can.push(x, y, replaceColor)
	if x > 0 {
		can.FloodFill(x-1, y, targetColor, replaceColor)
	}
	if x < can.Width-1 {
		can.FloodFill(x+1, y, targetColor, replaceColor)
	}
	if y > 0 {
		can.FloodFill(x, y-1, targetColor, replaceColor)
	}
	if y < can.Height-1 {
		can.FloodFill(x, y+1, targetColor, replaceColor)
	}
	return
}

func (can *Canvas) FullBox(x0, y0, x1, y1, color int) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	curX := x0
	for curX <= x1 {
		curY := y0
		for curY <= y1 {
			can.push(curX, curY, color)
			curY++
		}
		curX++
	}
}

func (can *Canvas) CopyBox(x0, y0, x1, y1 int) *Canvas {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	curX := x0
	subCanvas := NewCanvas(x1-x0+1, y1-y0+1)
	for curX <= x1 {
		curY := y0
		for curY <= y1 {
			subCanvas.Squares[curX][curY] = can.Squares[curX][curY]
			curY++
		}
		curX++
	}
	return subCanvas
}

func (can *Canvas) PasteBox(x0, y0 int, subCanvas *Canvas) {
	curX := x0
	subCurX := 0
	for subCurX < subCanvas.Width {
		subCurY := 0
		curY := y0
		for subCurY < subCanvas.Height {
			// BOUNDS CHECK
			// if curX > can.Width-1 || curY > can.Height-1 {
			// 	// SKIP
			// 	curY++
			// 	continue
			// }
			color := subCanvas.Squares[subCurX][subCurY]
			termbox.SetCell(curX, curY, ' ', termbox.Attribute(color), termbox.Attribute(color))
			can.Squares[curX][curY] = color
			subCurY++
			curY++
		}
		subCurX++
		curX++
	}
}

func (can *Canvas) Box(x0, y0, x1, y1, color int) {
	// Who is on the left?
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	// Who is above?
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	curX := x0
	curY := y0
	// 4 sides
	// TOP and BOTTOM
	for curX <= x1 {
		can.push(curX, y0, color)
		can.push(curX, y1, color)
		curX++
	}
	// LEFT and RIGHT with redundant ends
	for curY <= y1 {
		can.push(x0, curY, color)
		can.push(x1, curY, color)
		curY++
	}
}

// http://en.wikipedia.org/wiki/Bresenham's_line_algorithm
func (can *Canvas) BresenhamLine(x0, y0, x1, y1, color int) {
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy
	for {
		can.push(x0, y0, color)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}
