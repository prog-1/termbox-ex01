package main

import (
	"log"

	"github.com/nsf/termbox-go"
)

const (
	snakeBody    = '🐍'
	snakeFgColor = termbox.ColorRed
	// Use the default background color for the snake.
	snakeBgColor = termbox.ColorDefault
)

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

// newSnake returns a new struct instance representing a snake.
func newSnake() snake {
	return snake{coord{5, 5}, coord{1, 0}}
}

// Redraws the terminal.
func drawSnake(s snake) {
	termbox.Clear(snakeFgColor, snakeBgColor)
	termbox.SetCell(s.pos.x, s.pos.y, snakeBody, snakeFgColor, snakeBgColor)
	termbox.Flush()
}

// Makes a move for a snake and returns a snake with an updated position.
func moveSnake(s snake) snake {
	s.pos.x += s.dir.x
	s.pos.y += s.dir.y
	return s
}

// Tasks:
// 1. Modify the snake behavior, so it reacts on keyboard 'Up' and 'Down' arrows
// presses (termbox.KeyArrowUp and termbox.KeyArrowDown events).
func main() {
	// Initialize termbox.
	err := termbox.Init()
	if err != nil {
		log.Fatalf("failed to init termbox: %v", err)
	}
	defer termbox.Close()

	s := newSnake()
	// This is the main event loop.
	for {
		drawSnake(s)
		// Wait for an event.
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			// This switch is equivalent to
			// if ev.Key == termbox.KeyArrowLeft {
			//   s.dir = coord{-1, 0}
			// } else if ev.Key == termbox.KeyArrowRight {
			//   s.dir = coord{1, 0}
			// } else if ev.Key == termbox.KeyEsc {
			//   return
			// }
			switch ev.Key {
			case termbox.KeyArrowLeft:
				s.dir = coord{-1, 0}
				s = moveSnake(s)
			case termbox.KeyArrowRight:
				s.dir = coord{1, 0}
				s = moveSnake(s)
			// The program exits when a user presses 'Esc'.
			case termbox.KeyEsc:
				return
			}
		}
	}
}
