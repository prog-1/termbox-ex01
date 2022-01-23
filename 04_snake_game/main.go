package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	snakeBody    = '*'
	snakeFgColor = termbox.ColorRed
	// Use the default background color for the snake.
	snakeBgColor = termbox.ColorDefault
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
	// Movement direction of a snake.
	dir coord
}

// game represents a state of the game.
type game struct {
	sn snake
	// Game field dimensions.
	fieldWidth, fieldHeight int
}

// newSnake returns a new struct instance representing a snake.
// The snake is placed in a random position in the game field.
// The movement direction is right.
func newSnake(maxX, maxY int) snake {
	// rand.Intn generates a pseudo-random number:
	// https://pkg.go.dev/math/rand#Intn
	return snake{coord{rand.Intn(maxX), rand.Intn(maxY)}, coord{1, 0}}
}

// newGame returns a new game state.
func newGame() game {
	// Sets game field dimensions to the size of the terminal.
	w, h := termbox.Size()
	return game{
		fieldWidth:  w,
		fieldHeight: h,
		sn:          newSnake(w, h),
	}
}

// drawSnakePosition draws the current snake position (as a debugging
// information) in the buffer.
func drawSnakePosition(g game) {
	str := fmt.Sprintf("(%d, %d)", g.sn.pos.x, g.sn.pos.y)
	writeText(g.fieldWidth-len(str), 0, str, snakeFgColor, snakeBgColor)
}

// drawSnake draws the snake in the buffer.
func drawSnake(sn snake) {
	termbox.SetCell(sn.pos.x, sn.pos.y, snakeBody, snakeFgColor, snakeBgColor)
}

// Redraws the terminal.
func draw(g game) {
	// Clear the old "frame".
	termbox.Clear(snakeFgColor, snakeBgColor)
	drawSnakePosition(g)
	drawSnake(g.sn)
	// Update the "frame".
	termbox.Flush()
}

// mod is a custom implementation of the '%' (modulo) operator that always
// returns positive numbers.
func mod(n, m int) int {
	return (n%m + m) % m
}

// Makes a move for a snake. Returns a game with an updated position.
func moveSnake(g game) game {
	g.sn.pos.x = mod(g.sn.pos.x+g.sn.dir.x, g.fieldWidth)
	g.sn.pos.y = mod(g.sn.pos.y+g.sn.dir.y, g.fieldHeight)
	return g
}

func step(g game) game {
	g = moveSnake(g)
	draw(g)
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
	// This is the main event loop.
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyArrowDown:
					g.sn.dir = coord{0, 1}
				case termbox.KeyArrowUp:
					g.sn.dir = coord{0, -1}
				case termbox.KeyArrowLeft:
					g.sn.dir = coord{-1, 0}
				case termbox.KeyArrowRight:
					g.sn.dir = coord{1, 0}
				case termbox.KeyEsc:
					return
				}
			}
		default:
			g = step(g)
			time.Sleep(70 * time.Millisecond)
		}
	}
}
