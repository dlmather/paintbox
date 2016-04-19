package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

type CanvasStack struct {
	internal []*Canvas
}

func (cs *CanvasStack) Init() {
	cs.internal = make([]*Canvas, 0)
}

func (cs *CanvasStack) Push(can *Canvas) {
	cs.internal = append(cs.internal, can)
}

func (cs *CanvasStack) Pop() *Canvas {
	if len(cs.internal) == 0 {
		return nil
	}
	last := cs.internal[len(cs.internal)-1]
	if len(cs.internal) == 1 {
		return last
	}
	cs.internal = cs.internal[0 : len(cs.internal)-1]
	return last
}

func (can *Canvas) Init(w, h int) {
	can.Width = w
	can.Height = h
	can.Squares = make([][]int, can.Width)
	for i := 0; i < can.Width; i++ {
		can.Squares[i] = make([]int, can.Height)
	}
}

func NewCanvas(data [][]int) *Canvas {
	width := len(data)
	height := len(data[0])
	w, h := termbox.Size()
	if width != w || height != h {
		panic(fmt.Sprintf("Error : Trying to load into a window that is different than the original ORIG %d:%d, NEW %d:%d", width, height, w, h))
	}
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
	_, err = fHandle.WriteString(fmt.Sprintf("%v\n%v", can.Width, can.Height))
	for _, col := range can.Squares {
		_, err = fHandle.WriteString("\n")
		check(err)
		for _, val := range col {
			_, err = fHandle.WriteString(fmt.Sprintf("%d ", val))
			check(err)
		}
	}
	check(err)
}

func Load(fName string) *Canvas {
	fBytes, err := ioutil.ReadFile(fName)
	check(err)
	fData := string(fBytes)
	var width, height int
	lines := strings.Split(fData, "\n")

	fmt.Sscanf(lines[0], "%d", &width)
	fmt.Sscanf(lines[1], "%d", &height)
	lines = lines[2:]
	squares := make([][]int, width)
	fmt.Println(width, height, lines)
	//for index, line := range lines {
	//	vals := strings.Fields(line)
	//	squares[index] = make([]int, height)
	//	for internalIndex, val := range vals {
	//		fmt.Sscanf(val, "%d", &squares[index][internalIndex])
	//	}
	//}
	fmt.Println(squares)
	return NewCanvas(squares)
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

func (can *Canvas) FullBox(x0, y0, x1, y1, color int) {
	// Who is on the left?
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	// Who is above?
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	curX := x0
	for curX <= x1 {
		curY := y0
		for curY <= y1 {
			can.Squares[curX][curY] = color
			termbox.SetCell(curX, curY, ' ', termbox.Attribute(color), termbox.Attribute(color))
			curY++
		}
		curX++
	}
}

func (can *Canvas) CopyBox(x0, y0, x1, y1 int) *Canvas {
	// Who is on the left?
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	// Who is above?
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	curX := x0
	subCanvas := &Canvas{}
	subCanvas.Init(x1-x0+1, y1-y0+1)
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
		can.Squares[curX][y0] = color
		termbox.SetCell(curX, y0, ' ', termbox.Attribute(color), termbox.Attribute(color))
		can.Squares[curX][y1] = color
		termbox.SetCell(curX, y1, ' ', termbox.Attribute(color), termbox.Attribute(color))
		curX++
	}
	// LEFT and RIGHT with redundant ends
	for curY <= y1 {
		can.Squares[x0][curY] = color
		termbox.SetCell(x0, curY, ' ', termbox.Attribute(color), termbox.Attribute(color))
		can.Squares[x1][curY] = color
		termbox.SetCell(x1, curY, ' ', termbox.Attribute(color), termbox.Attribute(color))
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
