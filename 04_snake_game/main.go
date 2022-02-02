package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	appleBody    = '🍎'
	appleFg      = termbox.ColorLightRed
	appleBg      = termbox.ColorDefault
	snakeBody    = '*'
	snakeFgColor = termbox.ColorGreen
	snakeBgColor = termbox.ColorDefault
)

func writeText(x, y int, s string, fg, bg termbox.Attribute) { // writeText writes a string to the buffer.
	for i, ch := range s {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

type coord struct { // coord is a coordinate on a plane.
	x, y int
}

type snake struct { // snake is a struct with fields representing a snake.
	pos []coord // Position of a snake.
}
type apple struct {
	pos coord
}

type game struct { // game represents a state of the game.
	sn                      snake
	v                       coord
	ap                      apple
	fieldWidth, fieldHeight int // Game field dimensions.
}

// newSnake returns a new struct instance representing a snake.
// The snake is placed in a random position in the game field.
// The movement direction is right.
func newSnake(maxX, maxY int) snake {
	// rand.Intn generates a pseudo-random number:
	//return snake{[5]coord{{5, 5}, {4, 5}, {3, 5}, {2, 5}}}
	// https://pkg.go.dev/math/rand#Intn
	return snake{[]coord{{rand.Intn(maxX - 1), rand.Intn(maxY - 1)}, {rand.Intn(maxX - 2), rand.Intn(maxY - 1)}, {rand.Intn(maxX - 3), rand.Intn(maxY - 1)}, {rand.Intn(maxX - 4), rand.Intn(maxY - 1)}}}
}
func newApple(maxX, maxY int) apple {
	return apple{coord{rand.Intn(maxX), rand.Intn(maxY)}}
}

func newGame() game { // newGame returns a new game state.
	w, h := termbox.Size() // Sets game field dimensions to the size of the terminal.
	return game{
		fieldWidth:  w,
		fieldHeight: h,
		sn:          newSnake(w, h),
		v:           coord{1, 0},
		ap:          newApple(w, h),
	}
}

// drawSnakePosition draws the current snake position (as a debugging
// information) in the buffer.
func writeGameOver(g game) {
	gOC := termbox.ColorRed
	gBc := termbox.ColorDefault
	str := "Game Over"
	writeText((g.fieldWidth/2)-len(str), (g.fieldHeight/2)-1, str, gOC, gBc)
}
func drawSnakePosition(g game) {
	str := fmt.Sprintf("(%d, %d)", g.sn.pos[0].x, g.sn.pos[0].y)
	writeText(g.fieldWidth-len(str), 0, str, snakeFgColor, snakeBgColor)
}
func drawScore(g game) {
	temporaryScore := 10
	str := fmt.Sprintf("Score: %d", temporaryScore)
	writeText((g.fieldWidth/2)-len(str), 0, str, snakeFgColor, snakeBgColor)
}
func drawHearts(g game) { // maybe in future
	heart := "♥"
	str := fmt.Sprintf("%s%s%s", heart, heart, heart)
	writeText(0, 0, str, snakeFgColor, snakeBgColor)
}
func drawApple(ap apple) {
	termbox.SetCell(ap.pos.x, ap.pos.y, appleBody, appleFg, appleBg)
}
func drawSnake(sn snake) { // drawSnake draws the snake in the buffer.
	for _, pos := range sn.pos {
		termbox.SetCell(pos.x, pos.y, snakeBody, snakeFgColor, snakeBgColor)
	}
}
func drawBorder(g game) {
	ch := rune('#')
	bColor := termbox.ColorCyan
	bBackC := termbox.ColorDefault
	for i := 0; i != g.fieldWidth; i++ {
		termbox.SetCell(i, 1, ch, bColor, bBackC)
	}
	for i := 1; i != g.fieldHeight; i++ {
		termbox.SetCell(0, i, ch, bColor, bBackC)
	}
	for i := g.fieldWidth; i != 0; i-- {
		termbox.SetCell(i, g.fieldHeight-1, ch, bColor, bBackC)
	}
	for i := g.fieldHeight; i != 0; i-- {
		termbox.SetCell(g.fieldWidth-1, i, ch, bColor, bBackC)
	}
}
func borderExists(g game) {
	if g.sn.pos[0].x == 1 || g.sn.pos[0].y == 2 || g.sn.pos[0].x == g.fieldWidth-2 || g.sn.pos[0].y == g.fieldHeight-2 {
		writeGameOver(g)
		time.Sleep(100 * time.Millisecond)
		return
	}
}
func draw(g game) { // Redraws the terminal.
	termbox.Clear(snakeFgColor, snakeBgColor) // Clear the old "frame".
	drawApple(g.ap)
	drawBorder(g)
	writeGameOver(g)
	drawHearts(g)
	drawScore(g)
	drawSnakePosition(g)
	drawSnake(g.sn)
	termbox.Flush() // Update the "frame".
}

func moveSnake(s snake, v coord, fw, fh int) snake { // Makes a move for a snake. Returns a snake with an updated position.
	copy(s.pos[1:], s.pos[:])
	s.pos[0] = coord{s.pos[0].x + v.x, s.pos[0].y + v.y}
	return s
}

func step(g game) game {
	g.sn = moveSnake(g.sn, g.v, g.fieldWidth, g.fieldHeight)
	if g.ap.pos == g.sn.pos[0] {
		ap := newApple(g.fieldWidth, g.fieldHeight)
		drawApple(ap)
		g.sn.pos = append([]coord{{g.sn.pos[0].x, g.sn.pos[0].y}}, g.sn.pos...)
		//score++
		termbox.Flush()
	}
	draw(g)
	borderExists(g)
	return g
}

func moveLeft(g game) game  { g.v = coord{-1, 0}; return g }
func moveRight(g game) game { g.v = coord{1, 0}; return g }
func moveUp(g game) game    { g.v = coord{0, -1}; return g }
func moveDown(g game) game  { g.v = coord{0, 1}; return g }

func main() {
	err := termbox.Init() // Initialize termbox.
	if err != nil {
		log.Fatalf("failed to init termbox: %v", err)
	}
	defer termbox.Close()

	rand.Seed(time.Now().UnixNano()) // Other initialization.

	g := newGame()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(150 * time.Millisecond)
	defer ticker.Stop()

	for { // This is the main event loop.

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
