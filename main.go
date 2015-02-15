package main

import "github.com/nsf/termbox-go"
// import "math/rand"
import "time"

type Canvas struct {
	Width, Height int
}

type Cursor struct {
	xCoord, yCoord int
	color termbox.Attribute
}

func (c *Cursor) moveLeft() {
	if c.xCoord - 1 >= 0 {
		c.xCoord -= 1
	}
}

func (c *Cursor) moveRight() {
	w, _ := termbox.Size()
	if c.xCoord + 1 <= w {
		c.xCoord += 1
	}
}

func (c *Cursor) moveDown() {
	if c.yCoord - 1 >= 0 {
		c.yCoord -= 1
	}
}

func (c *Cursor) moveUp() {
	_, h := termbox.Size()
	if c.yCoord + 1 <= h {
		c.yCoord += 1
	}
}

func (c *Cursor) Position() (int, int) {
	return c.xCoord, c.yCoord
}

func (c *Cursor) placeColor() {
	x, y := c.Position()
	termbox.SetCell(x, y, ' ', c.color, c.color) 
}

func (c *Cursor) delete() {
	x, y := c.Position()
	termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault) 
}

func (c *Cursor) changeColor() {
	c.color = termbox.Attribute((c.color % 8) + 1) 
}

func draw(c *Cursor) {
	termbox.SetCursor(c.xCoord, c.yCoord)
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()
	
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	c := Cursor{xCoord: 0, yCoord: 0, color: termbox.ColorRed}
	cPtr := &c
	draw(cPtr)
loop:
	for {
		select {
		case ev := <-event_queue:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					break loop
				case termbox.KeyArrowDown:
					cPtr.moveUp()
				case termbox.KeyArrowUp:
					cPtr.moveDown()
				case termbox.KeyArrowRight:
					cPtr.moveRight()
				case termbox.KeyArrowLeft:
					cPtr.moveLeft()
				case termbox.KeySpace:
					cPtr.placeColor()
				case termbox.KeyTab:
					cPtr.changeColor()
				case termbox.KeyBackspace:
					cPtr.delete()
				default:
					draw(cPtr)
					time.Sleep(10 * time.Millisecond)
				}
			case termbox.EventError:
				panic(ev.Err)
			}
		}
		draw(cPtr)
	}
}