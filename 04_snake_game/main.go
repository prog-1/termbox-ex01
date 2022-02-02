package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
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
	pos []coord
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
	count                   int
}

// newSnake returns a new struct instance representing a snake.
// The snake is placed in a random position in the game field.
// The movement direction is right.

// func mod(n, m int) int {
// 	return (n%m + m) % m
// }

func newSnake(maxX, maxY int) snake {
	// rand.Intn generates a pseudo-random number:
	// https://pkg.go.dev/math/rand#Intn
	return snake{[]coord{{5, 5}, {4, 5}, {3, 5}, {2, 5}}}
}

func newApple(maxX, maxY int) apple {
	maxX, maxY = termbox.Size()
	return apple{coord{rand.Intn(maxX - 1), rand.Intn(maxY - 1)}}
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
	str := fmt.Sprintf("(%d, %d)", g.sn.pos[0].x, g.sn.pos[0].y)
	writeText(g.fieldWidth-len(str), 0, str, snakeFgColor, snakeBgColor)
}

func drawApplePosition(g game) {
	str := fmt.Sprintf("(%d, %d)", g.ap.po.x, g.ap.po.y)
	writeText(g.fieldWidth, 0, str, appleFgColor, appleBgColor)
}

// drawSnake draws the snake in the buffer.
func drawSnake(sn snake) {
	for _, pos := range sn.pos {
		termbox.SetCell(pos.x, pos.y, snakeBody, snakeFgColor, snakeBgColor)
	}
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

// Makes a move for a snake. Returns a snake with an updated position.
func moveSnake(s snake, v coord, fw, fh int) snake {
	copy(s.pos[1:], s.pos[:])
	s.pos[0] = coord{s.pos[0].x + v.x, s.pos[0].y + v.y}
	return s
}

func hitTheBorder(g game) bool {
	termbox.Clear(snakeFgColor, snakeBgColor)
	termbox.Clear(appleFgColor, appleBgColor)
	termbox.Clear(borderColor, appleIcon)
	if g.sn.pos[0].x == 0 || g.sn.pos[0].x == g.fieldWidth-1 || g.sn.pos[0].y == 0 || g.sn.pos[0].y == g.fieldHeight-1 {
		writeText(g.fieldWidth/2-19, g.fieldHeight/2, "SNAKE HIT THE BORDER, THE GAME IS OVER", termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
		termbox.Flush()
		time.Sleep(4 * time.Second)
		return true
	}
	return false
}

func biteItself(g game) bool {
	termbox.Clear(snakeFgColor, snakeBgColor)
	termbox.Clear(appleFgColor, appleBgColor)
	termbox.Clear(borderColor, appleIcon)
	for _, s := range g.sn.pos[2:] {
		if g.sn.pos[0] == s {
			writeText(g.fieldWidth/2-19, g.fieldHeight/2, "THE SNAKE BIT ITSELF, THE GAME IS OVER", termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
			termbox.Flush()
			time.Sleep(4 * time.Second)
			return true
		}
	}
	return false
}

func step(g game) game {
	for hitTheBorder(g) {
		os.Exit(0)
	}
	for biteItself(g) {
		os.Exit(0)
	}
	w, h := termbox.Size()
	if g.sn.pos[0].x == g.ap.po.x && g.sn.pos[0].y == g.ap.po.y {
		g.ap = newApple(w, h)
		g.count++
		// a = append([]T{x}, a...)
		g.sn.pos = append([]coord{{g.sn.pos[0].x, g.sn.pos[0].y}}, g.sn.pos...)
	}
	g.sn = moveSnake(g.sn, g.v, g.fieldWidth, g.fieldHeight)
	draw(g)
	return g
}

func drawCount(g game) {
	writeText(g.fieldWidth/2-3, 0, fmt.Sprint("Points:", g.count), termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault)
}

func moveLeft(g game) game {
	vv := coord{1, 0}
	if g.v != vv {
		g.v = coord{-1, 0}
	}
	return g
}
func moveRight(g game) game {
	vv := coord{-1, 0}
	if g.v != vv {
		g.v = coord{1, 0}
	}
	return g
}
func moveUp(g game) game {
	vv := coord{0, 1}
	if g.v != vv {
		g.v = coord{0, -1}
	}
	return g
}
func moveDown(g game) game {
	vv := coord{0, -1}
	if g.v != vv {
		g.v = coord{0, 1}
	}
	return g
}

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
