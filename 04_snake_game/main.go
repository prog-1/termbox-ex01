package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	snakeBody    = 'S'
	snakeFgColor = termbox.ColorGreen
	// Use the default background color for the snake.
	snakeBgColor = termbox.ColorDefault
	appleIcon    = '0'
	appleFgColor = termbox.ColorRed
	appleBgColor = termbox.ColorDefault
	borderColor  = termbox.ColorBlue
)

// writeText writes a string to the buffer.
func writeText(x, y int, s string, fg, bg termbox.Attribute) {
	for i, ch := range s {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

// coord is a coordinate on a plane.
type coord struct {
	x, y int
}

// snake is a struct with fields representing a snake.
type snake struct {
	// Position of a snake.
	pos coord
}

type apple struct {
	// Position of an apple.
	po coord
}

// game represents a state of the game.
type game struct {
	sn snake
	ap apple
	v  coord
	// Game field dimensions.
	fieldWidth, fieldHeight int
}

// newSnake returns a new struct instance representing a snake.
// The snake is placed in a random position in the game field.
// The movement direction is right.
func newSnake(maxX, maxY int) snake {
	// rand.Intn generates a pseudo-random number:
	// https://pkg.go.dev/math/rand#Intn
	return snake{coord{rand.Intn(maxX), rand.Intn(maxY)}}
}

func newApple(maxX, maxY int) apple {
	// rand.Intn generates a pseudo-random number:
	// https://pkg.go.dev/math/rand#Intn
	return apple{coord{rand.Intn(maxX), rand.Intn(maxY)}}
}

// newGame returns a new game state.
func newGame() game {
	// Sets game field dimensions to the size of the terminal.
	w, h := termbox.Size()
	return game{
		fieldWidth:  w,
		fieldHeight: h,
		sn:          newSnake(w, h),
		ap:          newApple(w, h),
		v:           coord{1, 0},
	}
}

// drawSnakePosition draws the current snake position (as a debugging
// information) in the buffer.
func drawSnakePosition(g game) {
	str := fmt.Sprintf("(%d, %d)", g.sn.pos.x, g.sn.pos.y)
	writeText(g.fieldWidth-len(str), 0, str, snakeFgColor, snakeBgColor)
}

func drawApplePosition(g game) {
	str := fmt.Sprintf("(%d, %d)", g.ap.po.x, g.ap.po.y)
	writeText(g.fieldWidth-len(str), 0, str, appleFgColor, appleBgColor)
}

// drawSnake draws the snake in the buffer.
func drawSnake(sn snake) {
	termbox.SetCell(sn.pos.x, sn.pos.y, snakeBody, snakeFgColor, snakeBgColor)
}

func drawApple(ap apple) {
	termbox.SetCell(ap.po.x, ap.po.y, appleIcon, appleFgColor, appleBgColor)
}

func drawBorder(width, height int) {
	for i := 0; i <= width; i++ {
		termbox.SetBg(i, 0, borderColor)
	}
	for i := 0; i <= height; i++ {
		termbox.SetBg(width-1, i, borderColor)
	}
	for i := 0; i <= width; i++ {
		termbox.SetBg(i, height-1, borderColor)
	}
	for i := 0; i <= height; i++ {
		termbox.SetBg(0, i, borderColor)
	}
}

// Redraws the terminal.
func draw(g game) {
	// Clear the old "frame".
	termbox.Clear(snakeFgColor, snakeBgColor)
	drawSnakePosition(g)
	drawSnake(g.sn)
	// termbox.Clear(appleFgColor, appleBgColor)
	drawApplePosition(g)
	drawApple(g.ap)
	drawCount(g)
	drawBorder(g.fieldWidth, g.fieldHeight)
	// Update the "frame".
	termbox.Flush()
}

// mod is a custom implementation of the '%' (modulo) operator that always
// returns positive numbers.
func mod(n, m int) int {
	return (n%m + m) % m
}

// Makes a move for a snake. Returns a snake with an updated position.
func moveSnake(s snake, v coord, fw, fh int) snake {
	s.pos.x = mod(s.pos.x+v.x, fw)
	s.pos.y = mod(s.pos.y+v.y, fh)
	return s
}

func appleCount(g game) int {
	w, h := termbox.Size()
	count := 0
	for g.sn.pos.x == g.ap.po.x && g.sn.pos.y == g.ap.po.y {
		g.ap = newApple(w, h)
		count++

	}
	return count
}

func drawCount(g game) {
	c := appleCount(g)
	writeText(g.fieldWidth/2, 0, fmt.Sprint("Points:", c), termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault)
}

func step(g game) game {
	w, h := termbox.Size()
	if g.sn.pos.x == g.ap.po.x && g.sn.pos.y == g.ap.po.y {
		g.ap = newApple(w, h)
	}
	g.sn = moveSnake(g.sn, g.v, g.fieldWidth, g.fieldHeight)
	draw(g)
	return g
}

func moveLeft(g game) game  { g.v = coord{-1, 0}; return g }
func moveRight(g game) game { g.v = coord{1, 0}; return g }
func moveUp(g game) game    { g.v = coord{0, -1}; return g }
func moveDown(g game) game  { g.v = coord{0, 1}; return g }

// Tasks:
func main() {
	// Initialize termbox.
	err := termbox.Init()
	if err != nil {
		log.Fatalf("failed to init termbox: %v", err)
	}
	defer termbox.Close()

	// Other initialization.
	rand.Seed(time.Now().UnixNano())
	g := newGame()
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(150 * time.Millisecond)
	defer ticker.Stop()

	// This is the main event loop.
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyArrowDown:
					g = moveDown(g)
				case termbox.KeyArrowUp:
					g = moveUp(g)
				case termbox.KeyArrowLeft:
					g = moveLeft(g)
				case termbox.KeyArrowRight:
					g = moveRight(g)
				case termbox.KeyEsc:
					return
				}
			}
		case <-ticker.C:
			g = step(g)
		}
	}
}
