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
	dir coord
}

// newSnake returns a new struct instance representing a snake.
func newSnake() snake {
	return snake{coord{3, 3}, coord{1, 0}}
}

// Redraws the terminal.
func drawSnake(s snake) {
	termbox.Clear(snakeFgColor, snakeBgColor)
	termbox.SetCell(s.pos.x, s.pos.y, snakeBody, snakeFgColor, snakeBgColor)
	termbox.Flush()
}

// Makes a move for a snake and returns a snake with an updated position.
func moveSnake(s snake, v coord) snake {
	return snake{pos: coord{x: s.pos.x + v.x, y: s.pos.y + v.y}}
}

// Makes a single iteration for a snake.

func step(s snake, dir coord) snake {
	x, y := s.pos.x, s.pos.y
	if x == 30 && y == 3 {
		s.dir.x = 0
		s.dir.y = 1
	}
	if x == 30 && y == 10 {
		s.dir.x = -1
		s.dir.y = 0
	}
	if x == 3 && y == 10 {
		s.dir.x = 0
		s.dir.y = -1
	}
	if x == 3 && y == 3 {
		s.dir.x = 1
		s.dir.y = 0
	}
	s = moveSnake(s, dir)
	drawSnake(s)
	return s
}

// Tasks:
// 1. Change the initial position of the snake and the movement direction.
// 2. Modify the snake trajectory, so it moves along the sides of a rectangle.
//    Hint: You may need to introduce an additional field to the snake struct.
func main() {
	// Initialize termbox.
	err := termbox.Init()
	if err != nil {
		log.Fatalf("failed to init termbox: %v", err)
	}
	defer termbox.Close()

	s := snake{pos: coord{5, 5}}
	// This is the main event loop.
	for {
		s = step(s, coord{1, 0})
		time.Sleep(100 * time.Millisecond)
	}
}
