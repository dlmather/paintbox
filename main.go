package main

import (
	"flag"

	"github.com/dlmather/paintbox/actions"
	"github.com/dlmather/paintbox/painter"
	"github.com/nsf/termbox-go"
)

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
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	myPainter := painter.New()
	myPainter.Draw()
	for {
		select {
		case ev := <-eventQueue:
			switch ev.Type {
			case termbox.EventKey:
				// TODO: Rewritable bindings
				if action, ok := actions.ActionBinding[ev.Key]; ok {
					err := action(myPainter)
					if err != nil {
						panic(err)
					}
				}
				myPainter.Draw()
			case termbox.EventError:
				panic(ev.Err)
			}
		}
	}
}
