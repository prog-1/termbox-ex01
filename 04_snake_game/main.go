package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	snakeBody    = '&'
	snakeFgColor = termbox.ColorRed
	// Use the default background color for the snake.
	snakeBgColor  = termbox.ColorDefault
	appleBody     = '0'
	appleFgColor  = termbox.ColorLightGreen
	appleBgColor  = termbox.ColorDefault
	borderColor   = termbox.ColorBlue
	energyBody    = 'E' // idea: if snake eats E, tail should become x2 longer
	energyFgColor = termbox.ColorBlue
	energyBgColor = termbox.ColorDefault
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

type applePos struct {
	pos coord
}

type eee struct {
	pos coord
}

// game represents a state of the game.
type game struct {
	sn    snake
	v     coord
	apple applePos
	e     eee
	score int
	// Game field dimensions.
	fieldWidth, fieldHeight int
}

func newApple(maxX, maxY int) applePos {
	return applePos{coord{rand.Intn(maxX), rand.Intn(maxY)}}
}

// newSnake returns a new struct instance representing a snake.
// The snake is placed in a random position in the game field.
// The movement direction is right.
func newSnake(maxX, maxY int) snake {
	// rand.Intn generates a pseudo-random number:
	// https://pkg.go.dev/math/rand#Intn
	return snake{[]coord{{5, 2}, {4, 2}}}
}

// newGame returns a new game state.
func newGame() game {
	// Sets game field dimensions to the size of the terminal.
	w, h := termbox.Size()
	var sc int
	return game{
		fieldWidth:  w,
		fieldHeight: h,
		sn:          newSnake(w, h),
		v:           coord{1, 0},
		apple:       newApple(w, h),
		score:       sc,
	}
}

func drawscore(g game) {
	writeText(g.fieldWidth/2, 0, fmt.Sprint("Your points:", g.score), termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault)

}

// drawSnakePosition draws the current snake position (as a debugging
// information) in the buffer.
func drawSnakePosition(g game) {
	str := fmt.Sprintf("(%d, %d)", g.sn.pos[0].x, g.sn.pos[0].y)
	writeText(g.fieldWidth-len(str), 0, str, snakeFgColor, snakeBgColor)
}
func drawApllePosition(g game) {
	str := fmt.Sprintf("(%d, %d)", g.apple.pos.x, g.apple.pos.y)
	writeText(g.fieldWidth-len(str), 0, str, appleFgColor, appleBgColor)
}

// drawSnake draws the snake in the buffer.
func drawSnake(sn snake) {
	for _, pos := range sn.pos {
		termbox.SetCell(pos.x, pos.y, snakeBody, snakeFgColor, snakeBgColor)
	}
}

func drawApple(apple applePos) {
	termbox.SetCell(apple.pos.x, apple.pos.y, appleBody, appleFgColor, appleBgColor)
}

// Redraws the terminal.
func draw(g game) {
	// Clear the old "frame".
	termbox.Clear(snakeFgColor, snakeBgColor)
	drawSnakePosition(g)
	drawSnake(g.sn)
	drawApllePosition(g)
	drawApple(g.apple)
	newBorder(g.fieldWidth, g.fieldHeight)
	drawscore(g)
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
	copy(s.pos[1:], s.pos[:])
	s.pos[0] = coord{s.pos[0].x + v.x, s.pos[0].y + v.y}
	return s
}

func step(g game) game {
	g.sn = moveSnake(g.sn, g.v, g.fieldWidth, g.fieldHeight)
	draw(g)
	return g
}

func aplleEaten(g game) (applePos, []coord, int) {
	w, h := termbox.Size()
	if g.apple.pos.x == g.sn.pos[0].x && g.apple.pos.y == g.sn.pos[0].y {
		g.apple = newApple(w, h)
		g.sn.pos = append(g.sn.pos, coord{g.sn.pos[len(g.sn.pos)-1].x - g.v.x, g.sn.pos[len(g.sn.pos)-1].y - g.v.y})
		g.score++
	}
	return g.apple, g.sn.pos, g.score
}

func moveLeft(g game) game  { g.v = coord{-1, 0}; return g }
func moveRight(g game) game { g.v = coord{1, 0}; return g }
func moveUp(g game) game    { g.v = coord{0, -1}; return g }
func moveDown(g game) game  { g.v = coord{0, 1}; return g }

func newBorder(width, height int) {
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

func borderCrash(g game) bool {
	if g.sn.pos[0].x == 0 || g.sn.pos[0].y == 0 || g.sn.pos[0].x == g.fieldWidth-1 || g.sn.pos[0].y == g.fieldHeight-1 {
		return true
	}
	return false
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

	ticker := time.NewTicker(70 * time.Millisecond)
	defer ticker.Stop()

	// This is the main event loop.
	for {
		termbox.Flush()
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
			if borderCrash(g) {
				termbox.Clear(termbox.ColorRed, termbox.ColorDefault)
				writeText(65, 3, `GAME OVER`, termbox.ColorRed, termbox.ColorDefault)
				termbox.Flush()
				time.Sleep(5 * time.Second)
				return
			}
			g = step(g)
		}
		g.apple, g.sn.pos, g.score = aplleEaten(g)
	}
}

func appleAtBorders(g game) applePos { //It is an important function, because I have a problem: apples can exist at borders
	w, h := termbox.Size()
	if g.apple.pos.x == 0 || g.apple.pos.x == w-1 || g.apple.pos.y == 0 || g.apple.pos.y == h-1 {
		g.apple = newApple(w, h)
	}
	return g.apple
}
