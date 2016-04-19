package main

import (
	"fmt"
)

import (
	"github.com/nsf/termbox-go"
)

// Stores state for on-screen cursor and its operations
type Cursor struct {
	xCoord, yCoord         int
	lineXCoord, lineYCoord int
	boxXCoord, boxYCoord   int // Re-use line?
	color                  termbox.Attribute
	colorInt               int
}

// TODO: Move based on canvas params, not termbox
func (cur *Cursor) MoveLeft() {
	if cur.xCoord > 0 {
		cur.xCoord -= 1
	}
}

func (cur *Cursor) MoveRight() {
	w, _ := termbox.Size()
	if cur.xCoord < w-1 {
		cur.xCoord += 1
	}
}

func (cur *Cursor) MoveDown() {
	if cur.yCoord > 0 {
		cur.yCoord -= 1
	}
}

func (cur *Cursor) MoveUp() {
	_, h := termbox.Size()
	if cur.yCoord < h-1 {
		cur.yCoord += 1
	}
}

func (cur *Cursor) Position() (int, int) {
	return cur.xCoord, cur.yCoord
}

// Places a single dot
func (cur *Cursor) PlaceColor(can *Canvas) {
	x, y := cur.Position()
	can.Squares[x][y] = cur.colorInt
	termbox.SetCell(x, y, ' ', cur.color, cur.color)
}

// Fills an area
func (cur *Cursor) FloodFill(can *Canvas) {
	x, y := cur.Position()
	targetColor := can.Squares[x][y]
	replaceColor := cur.colorInt
	can.FloodFill(x, y, targetColor, replaceColor)
}

// Draws a box starting at initial corner
// Ending at final corner
func (cur *Cursor) Box(can *Canvas) {
	x, y := cur.Position()
	lineX := cur.lineXCoord
	lineY := cur.lineYCoord
	if lineX == -1 && lineY == -1 {
		termbox.SetCell(x, y, 'B', termbox.Attribute((cur.color%8)+1), cur.color)
		cur.lineXCoord = x
		cur.lineYCoord = y
	} else {
		can.Box(lineX, lineY, x, y, cur.colorInt)
		cur.lineXCoord = -1
		cur.lineYCoord = -1
	}
}

// Draws a filled box starting at initial corner
// Ending at final corner
func (cur *Cursor) FullBox(can *Canvas) {
	x, y := cur.Position()
	lineX := cur.lineXCoord
	lineY := cur.lineYCoord
	if lineX == -1 && lineY == -1 {
		termbox.SetCell(x, y, 'B', termbox.Attribute((cur.color%8)+1), cur.color)
		cur.lineXCoord = x
		cur.lineYCoord = y
	} else {
		can.FullBox(lineX, lineY, x, y, cur.colorInt)
		cur.lineXCoord = -1
		cur.lineYCoord = -1
	}
}

// Draws a line between the point selected when first run
// and point selected when second run
func (cur *Cursor) Line(can *Canvas) {
	x, y := cur.Position()
	lineX := cur.lineXCoord
	lineY := cur.lineYCoord
	if lineX == -1 && lineY == -1 {
		termbox.SetCell(x, y, 'x', termbox.Attribute((cur.color%8)+1), cur.color)
		cur.lineXCoord = x
		cur.lineYCoord = y
	} else {
		can.BresenhamLine(lineX, lineY, x, y, cur.colorInt)
		cur.lineXCoord = -1
		cur.lineYCoord = -1
	}
}

// Useful for debugging
func (cur *Cursor) Pos() {
	x, y := cur.Position()
	fmt.Printf("\t%v:%v", x, y)
}

func (cur *Cursor) Delete(can *Canvas) {
	x, y := cur.Position()
	can.Squares[x][y] = 0
	termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
}

func (cur *Cursor) ChangeColor() {
	cur.color = termbox.Attribute((cur.color % 8) + 1)
	cur.colorInt = (cur.colorInt % 8) + 1
}
