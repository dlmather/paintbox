package main

import (
	"io/ioutil"
	"time"
	"fmt"
	"os"
	"flag"
	"strings"
	"github.com/nsf/termbox-go"
)

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

func (can *Canvas) FloodFill(x, y, targetColor, replaceColor int) {
	// fmt.Printf("FLOODING %v,%v\n", targetColor, replaceColor)
	// return
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


type Cursor struct {
	xCoord, yCoord int
	color          termbox.Attribute
	colorInt       int
}

func (cur *Cursor) moveLeft() {
	if cur.xCoord > 0 {
		cur.xCoord -= 1
	}
}

func (cur *Cursor) moveRight() {
	w, _ := termbox.Size()
	if cur.xCoord < w-1 {
		cur.xCoord += 1
	}
}

func (cur *Cursor) moveDown() {
	if cur.yCoord > 0 {
		cur.yCoord -= 1
	}
}

func (cur *Cursor) moveUp() {
	_, h := termbox.Size()
	if cur.yCoord < h-1 {
		cur.yCoord += 1
	}
}

func (cur *Cursor) Position() (int, int) {
	return cur.xCoord, cur.yCoord
}

func (cur *Cursor) placeColor(can *Canvas) {
	x, y := cur.Position()
	can.Squares[x][y] = cur.colorInt
	termbox.SetCell(x, y, ' ', cur.color, cur.color)
}

func (cur *Cursor) FloodFill(can *Canvas) {
	x, y := cur.Position()
	targetColor := can.Squares[x][y]
	replaceColor := cur.colorInt
	can.FloodFill(x, y, targetColor, replaceColor)
}

func (cur *Cursor) delete(can *Canvas) {
	x, y := cur.Position()
	can.Squares[x][y] = 0
	termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
}

func (cur *Cursor) changeColor() {
	cur.color = termbox.Attribute((cur.color % 8) + 1)
	cur.colorInt = (cur.colorInt % 8) + 1
}

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
	// disTimer := time.NewTicker(2 * time.Second)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	c := Cursor{xCoord: 0, yCoord: 0, color: termbox.ColorRed, colorInt: 2}
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
				case termbox.KeyCtrlF:
					cPtr.FloodFill(canPtr)
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
