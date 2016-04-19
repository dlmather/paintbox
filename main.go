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
	if config.LoadPath != "" {
		canvas = *Load(config.LoadPath)
	} else {
		canvas.Init()
	}
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
				case termbox.KeyCtrlQ:
					canPtr.Save()
					break loop
				case termbox.KeyEsc:
					break loop
				case termbox.KeyArrowDown:
					cPtr.MoveUp()
				case termbox.KeyArrowUp:
					cPtr.MoveDown()
				case termbox.KeyArrowRight:
					cPtr.MoveRight()
				case termbox.KeyArrowLeft:
					cPtr.MoveLeft()
				case termbox.KeyTab:
					cPtr.ChangeColor()
				case termbox.KeyCtrlB:
					cPtr.Box(canPtr)
				case termbox.KeyCtrlL:
					cPtr.Line(canPtr)
				case termbox.KeyCtrlF:
					cPtr.FloodFill(canPtr)
				case termbox.KeyCtrlP:
					cPtr.Pos()
				case termbox.KeySpace:
					cPtr.PlaceColor(canPtr)
				case termbox.KeyBackspace, termbox.KeyBackspace2:
					cPtr.Delete(canPtr)
				}
				draw(cPtr, canPtr)
			case termbox.EventError:
				panic(ev.Err)
			}
		}
	}
}
