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
	maxSteps     = 3
	startX       = 10
	startY       = 10
	// I don't know why, but with startPos coord in snake struct the program won't work.
)

// coord is a coordinate on a plane.
type coord struct {
	x, y int
}

// snake is a struct with fields representing a snake.
type snake struct {
	// Position of a snake.
	pos coord
	// startPos coord
	// maxSteps int
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

	s := snake{pos: coord{startX, startY}}
	// s := snake{pos: coord{10, 10}, maxSteps: 3}
	// s.startPos = coord{s.pos.x, s.pos.y}
	dir := coord{0, -1}
	// This is the main event loop.
	for {
		// if s.pos.x == s.startPos.x && s.pos.y == s.startPos.y {
		// 	dir = coord{0, -1}
		// }
		// if s.pos.y == s.startPos.y-s.maxSteps && s.pos.x == s.startPos.x {
		// 	dir = coord{1, 0}
		// }
		// if s.pos.x == s.startPos.x+s.maxSteps && s.pos.y == s.startPos.y-s.maxSteps {
		// 	dir = coord{0, 1}
		// }
		// if s.pos.y == s.startPos.y && s.pos.x == s.startPos.x+s.maxSteps {
		// 	dir = coord{-1, 0}
		// }
		if s.pos.x == startX && s.pos.y == startY {
			dir = coord{0, -1}
		}
		if s.pos.y == startY-maxSteps && s.pos.x == startX {
			dir = coord{1, 0}
		}
		if s.pos.x == startX+maxSteps && s.pos.y == startY-maxSteps {
			dir = coord{0, 1}
		}
		if s.pos.y == startY && s.pos.x == startX+maxSteps {
			dir = coord{-1, 0}
		}
		s = step(s, dir)
		time.Sleep(100 * time.Millisecond)
	}
}
