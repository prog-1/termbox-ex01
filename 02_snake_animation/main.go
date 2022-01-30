package main

import (
	"log"
	"time"

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
	return snake{coord{5, 7}, coord{1, 0}}
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

// Makes a single iteration for a snake.
func step(s snake) snake {
	for _, pivot := range []struct {
		pos, dir coord
	}{
		{coord{5, 7}, coord{0, 1}},
		{coord{5, 10}, coord{-1, 0}},
		{coord{1, 10}, coord{0, -1}},
		{coord{1, 7}, coord{1, 0}},
	} {
		if (coord{s.pos.x, s.pos.y} == pivot.pos) {
			s.dir = pivot.dir
		}
	}
	s = moveSnake(s)
	drawSnake(s)
	return s
}

// Tasks:
// 1. Change the initial position of the snake and the movement direction.
// 2. Modify the snake behavior, so it makes 3 steps up, 3 steps down, then repeats.
//    Hint: You may need to introduce an additional field to the snake struct.
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
		s = step(s)
		time.Sleep(100 * time.Millisecond)
	}
}
