package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

import (
	"github.com/nsf/termbox-go"
)


// Set the cursor to its defined location + flush
func draw(cur *Cursor, canPtr *Canvas) {
	termbox.SetCursor(cur.xCoord, cur.yCoord)
	termbox.Flush()
}

// Attempt to load canvas file
func load(path string) (*Canvas, error) {
	fBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fData := string(fBytes)
	lines := strings.Split(fData, "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("Bad paintbox file")
	}
	return nil, fmt.Errorf("Can't load yet.")
}

// For now, config contains path for load
type Config struct {
	LoadPath string
}

var config Config

func init() {
	flag.StringVar(&config.LoadPath, "load", "", "file path to load a previous work from")
	flag.Parse()
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
}

func main() {
	defer termbox.Close()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	c := Cursor{xCoord: 0, yCoord: 0, color: termbox.ColorRed, colorInt: 2, lineXCoord: -1, lineYCoord: -1}
	canvas := Canvas{}
	canvas.Init()
	canPtr := &canvas
	cPtr := &c
	draw(cPtr, canPtr)
loop:
	for {
		select {
		case ev := <-event_queue:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					// canPtr.Save()
					break loop
				case termbox.KeyArrowDown:
					cPtr.moveUp()
				case termbox.KeyArrowUp:
					cPtr.moveDown()
				case termbox.KeyArrowRight:
					cPtr.moveRight()
				case termbox.KeyArrowLeft:
					cPtr.moveLeft()
				case termbox.KeyTab:
					cPtr.changeColor()
				case termbox.KeyCtrlL:
					cPtr.Line(canPtr)
				case termbox.KeyCtrlF:
					cPtr.FloodFill(canPtr)
				case termbox.KeyCtrlP:
					cPtr.Pos()
				case termbox.KeySpace:
					cPtr.placeColor(canPtr)
				case termbox.KeyBackspace, termbox.KeyBackspace2:
					cPtr.delete(canPtr)
				}
				draw(cPtr, canPtr)
			case termbox.EventError:
				panic(ev.Err)
			}
		}
	}
}
