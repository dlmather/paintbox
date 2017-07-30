package display

import (
	"errors"

	"github.com/nsf/termbox-go"
)

const (
	colorWord = "Current Color:"
)

type Display struct {
	Width, Height       int
	X, Y                int
	HorizontalSeperator rune
	VerticalSeperator   rune
	CurrentColor        int
}

func New(canW, canH int) *Display {
	return &Display{
		Width:               canW - 2,
		Height:              2,
		X:                   0,
		Y:                   canH - 3,
		HorizontalSeperator: '-',
		VerticalSeperator:   '|',
		CurrentColor:        1,
	}
}

func drawWord(word string, x, y, limit int) error {
	if len(word) > limit {
		return errors.New("length of word exceeds limit")
	}
	for i, r := range word {
		termbox.SetCell(x+i, y, r, termbox.Attribute(0), termbox.Attribute(0))
	}
	return nil
}

func (d *Display) Draw() {
	// Draw borders
	for curX := d.X; curX <= d.Width; curX++ {
		termbox.SetCell(curX, d.Y, d.HorizontalSeperator, termbox.Attribute(0), termbox.Attribute(0))
		termbox.SetCell(curX, d.Y+d.Height, d.HorizontalSeperator, termbox.Attribute(0), termbox.Attribute(0))
	}

	_ = drawWord(colorWord, d.X, d.Y+1, d.Width)
	termbox.SetCell(d.X+len(colorWord)+1, d.Y+1, ' ', termbox.Attribute(d.CurrentColor), termbox.Attribute(d.CurrentColor))
	termbox.Flush()
}
