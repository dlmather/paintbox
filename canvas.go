package main

import (
	"fmt"
	"os"
	"time"
)

import (
	"github.com/nsf/termbox-go"
)

// A representation of the actual state of the canvas
type Canvas struct {
	Width, Height int
	Squares       [][]int
}

func (can *Canvas) Init() {
	w, h := termbox.Size()
	can.Width = w
	can.Height = h
	can.Squares = make([][]int, can.Width)
	for i := 0; i < can.Width; i++ {
		can.Squares[i] = make([]int, can.Height)
	}
}

func NewCanvas(data [][]int) *Canvas {
	can := Canvas{Width: len(data), Height: len(data[0]), Squares: data}
	return &can
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (can *Canvas) Save() {
	fHandle, err := os.Create(fmt.Sprintf("./paintbox-%v.pnt", time.Now()))
	check(err)
	defer fHandle.Close()
	_, err = fHandle.WriteString(fmt.Sprintf("%v:%v\n%v", can.Width, can.Height, can.Squares))
	check(err)
}

func (can *Canvas) draw() {
	for x, column := range can.Squares {
		for y, color := range column {
			termbox.SetCell(x, y, ' ', termbox.Attribute(color), termbox.Attribute(color))
		}
	}
}

// http://en.wikipedia.org/wiki/Flood_fill
func (can *Canvas) FloodFill(x, y, targetColor, replaceColor int) {
	if targetColor == replaceColor {
		return
	}
	if can.Squares[x][y] != targetColor {
		return
	}
	can.Squares[x][y] = replaceColor
	termbox.SetCell(x, y, ' ', termbox.Attribute(replaceColor), termbox.Attribute(replaceColor))
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
		can.Squares[x0][y0] = color
		termbox.SetCell(x0, y0, ' ', termbox.Attribute(color), termbox.Attribute(color))
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
