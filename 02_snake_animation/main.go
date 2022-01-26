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
	dir      coord
	movedown int
}

// newSnake returns a new struct instance representing a snake.
func newSnake() snake {
	return snake{coord{0, 3}, coord{2, 3}, 0}
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
	if s.movedown == 1 {
		s.pos.y += s.dir.y
		s.movedown = 0
	} else {
		s.pos.y -= s.dir.y
		s.movedown = 1
	}
	return s
}

// Makes a single iteration for a snake.
func step(s snake) snake {
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

//go run ./02_snake_animation/main.go
