package main

import "github.com/nsf/termbox-go"

// import "io/ioutil"
// import "math/rand"
import "time"
import "fmt"
import "os"

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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (can *Canvas) Save() {
	fHandle, err := os.Create(fmt.Sprintf("./paintbox-%v.pnt", time.Now()))
	check(err)
	defer fHandle.Close()
	_, err = fHandle.WriteString(fmt.Sprintf("%v", *can))
	check(err)
}

type Cursor struct {
	xCoord, yCoord int
	color          termbox.Attribute
	colorInt int
}

func (c *Cursor) moveLeft() {
	if c.xCoord-1 >= 0 {
		c.xCoord -= 1
	}
}

func (c *Cursor) moveRight() {
	w, _ := termbox.Size()
	if c.xCoord+1 <= w {
		c.xCoord += 1
	}
}

func (c *Cursor) moveDown() {
	if c.yCoord-1 >= 0 {
		c.yCoord -= 1
	}
}

func (c *Cursor) moveUp() {
	_, h := termbox.Size()
	if c.yCoord+1 <= h {
		c.yCoord += 1
	}
}

func (c *Cursor) Position() (int, int) {
	return c.xCoord, c.yCoord
}

func (c *Cursor) placeColor(can *Canvas) {
	x, y := c.Position()
	can.Squares[x][y] = c.colorInt
	termbox.SetCell(x, y, ' ', c.color, c.color)
}

func (c *Cursor) delete(can *Canvas) {
	x, y := c.Position()
	can.Squares[x][y] = 0
	termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
}

func (c *Cursor) changeColor() {
	c.color = termbox.Attribute((c.color % 8) + 1)
	c.colorInt = (c.colorInt % 8) + 1
}

func draw(c *Cursor, canPtr *Canvas) {
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
	// disTimer := time.NewTicker(2 * time.Second)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	c := Cursor{xCoord: 0, yCoord: 0, color: termbox.ColorRed, colorInt: 1}
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
					canPtr.Save()
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
				case termbox.KeySpace:
					cPtr.placeColor(canPtr)
				case termbox.KeyBackspace, termbox.KeyBackspace2:
					cPtr.delete(canPtr)
				// default:
				// 	draw(cPtr)
				// 	time.Sleep(10 * time.Millisecond)
				}
				draw(cPtr, canPtr)
			case termbox.EventError:
				panic(ev.Err)
			}
		// case <-disTimer.C:
		// 	termbox.HideCursor()
		// 	termbox.Flush()
		}
	}
}
