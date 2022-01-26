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
	// changing size of terminal, you also change size of rectangle where we can see snake
	//  So everyone will have their own game size and snakes position depending on the size of their terminal...
	return snake{coord{50, 6}, coord{1, 0}}
}

// Redraws the terminal.
func drawSnake(s snake) {
	termbox.Clear(snakeFgColor, snakeBgColor)
	termbox.SetCell(s.pos.x, s.pos.y, snakeBody, snakeFgColor, snakeBgColor)
	termbox.Flush()
}

// Makes a move for a snake and returns a snake with an updated position
// and
// Makes a single iteration for a snake
func stepForward(s snake) snake {
	s.dir = coord{1, 0}
	s.pos.x += s.dir.x
	s.pos.y += s.dir.y
	drawSnake(s)
	return s
}

func stepBack(s snake) snake {
	s.dir = coord{1, 0}
	s.pos.x -= s.dir.x
	s.pos.y -= s.dir.y
	drawSnake(s)
	return s
}

func stepUp(s snake) snake {
	s.dir = coord{0, 1}
	s.pos.x += s.dir.x
	s.pos.y += s.dir.y
	drawSnake(s)
	return s
}

func stepDown(s snake) snake {
	s.dir = coord{0, 1}
	s.pos.x -= s.dir.x
	s.pos.y -= s.dir.y
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

	s := newSnake()
	// This is the main event loop.
	for { // its perfectly works even without any if statements
		for s.pos.x != 5 {
			s = stepBack(s)
			time.Sleep(100 * time.Millisecond)
		}
		for s.pos.y != 5 {
			s = stepDown(s)
			time.Sleep(100 * time.Millisecond)
		}
		for s.pos.x != 15 {
			s = stepForward(s)
			time.Sleep(100 * time.Millisecond)
		}
		for s.pos.y != 10 {
			s = stepUp(s)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
