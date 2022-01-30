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
	snakeBody    = '*'
	snakeFgColor = termbox.ColorRed
	// Use the default background color for the snake.
	snakeBgColor = termbox.ColorDefault
	appleBody    = 'O'
	appleFgColor = termbox.ColorGreen
	appleBgColor = termbox.ColorDefault
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
	pos   coord
	score int
}

// game represents a state of the game.
type game struct {
	sn snake
	v  coord
	a  apple
	// Game field dimensions.
	fieldWidth, fieldHeight int
}

// newSnake returns a new struct instance representing a snake.
// The snake is placed in a random position in the game field.
// The movement direction is right.
func newSnake(g game) snake {
	// rand.Intn generates a pseudo-random number:
	// https://pkg.go.dev/math/rand#Intn
	// return snake{coord{rand.Intn(maxX), rand.Intn(maxY)}}
	g.sn.pos = coord{rand.Intn(g.fieldWidth), rand.Intn(g.fieldHeight)}
	for {
		if g.sn.pos.x == 0 || g.sn.pos.x == g.fieldWidth-1 || g.sn.pos.y == 0 || g.sn.pos.y == 1 || g.sn.pos.y == g.fieldHeight-1 {
			g.sn.pos = coord{rand.Intn(g.fieldWidth), rand.Intn(g.fieldHeight)}
		} else {
			break
		}
	}
	return snake{g.sn.pos}
}

func newApple(g game) apple {
	g.a.pos = coord{rand.Intn(g.fieldWidth), rand.Intn(g.fieldHeight)}
	for {
		if g.a.pos.x == 0 || g.a.pos.x == g.fieldWidth-1 || g.a.pos.y == 0 || g.a.pos.y == 1 || g.a.pos.y == g.fieldHeight-1 {
			g.a.pos = coord{rand.Intn(g.fieldWidth), rand.Intn(g.fieldHeight)}
		} else {
			break
		}
	}
	return apple{g.a.pos, g.a.score}
}

// newGame returns a new game state.
func newGame() game {
	// Sets game field dimensions to the size of the terminal.
	w, h := termbox.Size()
	return game{
		fieldWidth:  w,
		fieldHeight: h,
		v:           coord{1, 0},
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

func drawWalls(g game) {
	for i := 0; i < g.fieldWidth; i++ {
		termbox.SetBg(i, 1, termbox.ColorWhite)
		termbox.SetBg(i, g.fieldHeight-1, termbox.ColorWhite)
	}
	for j := 0; j < g.fieldHeight; j++ {
		termbox.SetBg(0, j+1, termbox.ColorWhite)
		termbox.SetBg(g.fieldWidth-1, j+1, termbox.ColorWhite)
	}
}

func drawApple(a apple) {
	termbox.SetCell(a.pos.x, a.pos.y, appleBody, appleFgColor, appleBgColor)
}

func drawScore(g game) {
	writeText(g.fieldWidth/2, 0, fmt.Sprint("Score:", g.a.score), termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
}

// Redraws the terminal.
func draw(g game) {
	// Clear the old "frame".
	termbox.Clear(snakeFgColor, snakeBgColor)
	drawWalls(g)
	drawApple(g.a)
	drawSnakePosition(g)
	drawSnake(g.sn)
	drawScore(g)
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

func collisions(g game) {
	if g.sn.pos.x == 0 || g.sn.pos.x == g.fieldWidth-1 || g.sn.pos.y == 1 || g.sn.pos.y == g.fieldHeight-1 {
		writeText(g.fieldWidth/2, g.fieldHeight/2, "Game Over", termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
		termbox.Flush()
		time.Sleep(5 * time.Second)
		termbox.Clear(snakeFgColor, snakeBgColor)
		termbox.Flush()
		os.Exit(0)
	}
}

func step(g game) game {
	g.sn = moveSnake(g.sn, g.v, g.fieldWidth, g.fieldHeight)
	collisions(g)
	if g.sn.pos == g.a.pos {
		g.a.score++
		g.a = newApple(g)
	}
	draw(g)
	return g
}

func moveLeft(g game) game  { g.v = coord{-1, 0}; return g }
func moveRight(g game) game { g.v = coord{1, 0}; return g }
func moveUp(g game) game    { g.v = coord{0, -1}; return g }
func moveDown(g game) game  { g.v = coord{0, 1}; return g }

func mainMenu() (choice int) {
	fmt.Println(`Select difficulty:
1) Easy
2) Medium
3) Hard`)
	fmt.Scanln(&choice)
	return
}

// Tasks:
func main() {
	var ticker *time.Ticker
	choice := mainMenu()
	if choice == 1 {
		ticker = time.NewTicker(250 * time.Millisecond)
	} else if choice == 2 {
		ticker = time.NewTicker(150 * time.Millisecond)
	} else if choice == 3 {
		ticker = time.NewTicker(100 * time.Millisecond)
	} else {
		fmt.Println("ERR: wrong choice")
	}
	defer ticker.Stop()

	// Initialize termbox.
	err := termbox.Init()
	if err != nil {
		log.Fatalf("failed to init termbox: %v", err)
	}
	defer termbox.Close()

	// Other initialization.
	rand.Seed(time.Now().UnixNano())
	g := newGame()
	g.sn = newSnake(g)
	g.a = newApple(g)

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
