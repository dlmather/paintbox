package actions

import (
	"errors"

	"github.com/dlmather/paintbox/painter"
	"github.com/nsf/termbox-go"
)

type Action func(*painter.Painter) error

var ErrQuit = errors.New("you've quit")

func Quit(_ *painter.Painter) error {
	return ErrQuit
}

var ActionBinding = map[termbox.Key]Action{
	termbox.KeyEsc: Quit,
	// Movement
	termbox.KeyArrowDown:  painter.MoveDown,
	termbox.KeyArrowUp:    painter.MoveUp,
	termbox.KeyArrowRight: painter.MoveRight,
	termbox.KeyArrowLeft:  painter.MoveLeft,
	// Basics
	termbox.KeyTab:        painter.CycleColor,
	termbox.KeySpace:      painter.Point,
	termbox.KeyBackspace:  painter.Erase,
	termbox.KeyBackspace2: painter.Erase,
	// Shapes
	termbox.KeyCtrlX: painter.FullBox,
	termbox.KeyCtrlB: painter.Box,
	termbox.KeyCtrlL: painter.Line,
	termbox.KeyCtrlF: painter.FloodFill,
}
