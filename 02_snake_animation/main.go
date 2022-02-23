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

	s := snake{pos: coord{5, 8}}
	// This is the main event loop.
	dir := coord{1, 0}
	for {
		x, y := s.pos.x, s.pos.y
		if x == 5 && y == 8 {
			dir = coord{1, 0}
		} else if x == 30 && y == 8 {
			dir = coord{0, 1}
		} else if x == 30 && y == 13 {
			dir = coord{-1, 0}
		} else if x == 5 && y == 13 {
			dir = coord{0, -1}
		}
		s = step(s, dir)
		time.Sleep(100 * time.Millisecond)
	}
}
