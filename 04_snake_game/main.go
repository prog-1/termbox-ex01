package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	snakeBody    = '*'
	snakeFgColor = termbox.ColorGreen
	// Use the default background color for the snake.
	snakeBgColor       = termbox.ColorDefault
	appleBody          = 'O'
	appleFgColor       = termbox.ColorRed
	appleBgColor       = termbox.ColorDefault
	borderBody         = '#'
	borderFgColor      = termbox.ColorWhite
	borderBgColor      = termbox.ColorDefault
	gameover1Body      = '-'
	gameover2Body      = '|'
	gameover1FgColor   = termbox.ColorWhite
	gameover1BgColor   = termbox.ColorDefault
	gameover2FgColor   = termbox.ColorWhite
	gameover2BgColor   = termbox.ColorDefault
	gameoverCornerBody = '+'
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
type apple struct {
	pos coord
}

// snake is a struct with fields representing a snake.
type snake struct {
	// Position of a snake.
	pos []coord
}

// game represents a state of the game.
type game struct {
	sn                      snake
	v                       coord
	ap                      apple
	fieldWidth, fieldHeight int
	score                   int
}

func newAplle(maxX, maxY int) apple {
	return apple{coord{rand.Intn(maxX), rand.Intn(maxY)}}
}

// newSnake returns a new struct instance representing a snake.
// The snake is placed in a random position in the game field.
// The movement direction is right.
func newSnake(maxX, maxY int) snake {
	var snake snake
	y := rand.Intn(maxY)
	x := rand.Intn(maxX)
	for x == 0 || y == 0 {
		y = rand.Intn(maxY)
		x = rand.Intn(maxX)
	}

	snake.pos = append(snake.pos, coord{x, y})
	snake.pos = append(snake.pos, coord{snake.pos[0].x - 1, snake.pos[0].y})
	snake.pos = append(snake.pos, coord{snake.pos[1].x - 1, snake.pos[0].y})

	return snake

}

// newGame returns a new game state.
func newGame() game {
	var sc int
	// Sets game field dimensions to the size of the terminal.
	w, h := termbox.Size()
	return game{
		fieldWidth:  w,
		fieldHeight: h,
		sn:          newSnake(w, h),
		v:           coord{1, 0},
		ap:          newAplle(w, h),
		score:       sc,
	}
}

// drawSnakePosition draws the current snake position (as a debugging
// information) in the buffer.
func drawSnakePosition(g game) {
	str := fmt.Sprintf("(%d, %d)", g.sn.pos[0].x, g.sn.pos[0].y)
	writeText(g.fieldWidth-len(str), 0, str, snakeFgColor, snakeBgColor)
}
func drawApllePosition(g game) {
	str := fmt.Sprintf("(%d, %d)", g.ap.pos.x, g.ap.pos.y)
	writeText(g.fieldWidth-len(str), 0, str, appleFgColor, appleBgColor)
}

// drawSnake draws the snake in the buffer.
func drawSnake(sn snake) {
	for _, pos := range sn.pos {
		termbox.SetCell(pos.x, pos.y, snakeBody, snakeFgColor, snakeBgColor)
	}

}
func drawApple(ap apple) {
	termbox.SetCell(ap.pos.x, ap.pos.y, appleBody, appleFgColor, appleBgColor)
}

func drawborder() {
	w, h := termbox.Size()
	for i := 0; i < w-1; i++ {
		termbox.SetCell(i, 0, borderBody, borderFgColor, borderBgColor)
		termbox.SetCell(i, h-1, borderBody, borderFgColor, borderBgColor)
	}
	for i := 0; i < h-1; i++ {
		termbox.SetCell(0, i, borderBody, borderFgColor, borderBgColor)
		termbox.SetCell(w-1, i, borderBody, borderFgColor, borderBgColor)
	}

}

func borderCrash(g game) bool {
	if g.sn.pos[0].x == 0 || g.sn.pos[0].y == 0 || g.sn.pos[0].x == g.fieldWidth-1 || g.sn.pos[0].y == g.fieldHeight-1 {
		return true
	}
	return false

}
func drawscore(g game) {
	w, _ := termbox.Size()
	score := fmt.Sprintf("--YOUR SCORE: %d--", g.score)
	writeText(w/2-7, 1, score, gameover1FgColor, gameover1BgColor)
}
func gameover(g game) {
	w, h := termbox.Size()
	for i := w/2 - 8; i < w/2+9; i++ {
		termbox.SetCell(i, h/2-3, gameover1Body, gameover1FgColor, gameover1BgColor)
		termbox.SetCell(i, h/2+3, gameover1Body, gameover1FgColor, gameover1BgColor)

	}
	for i := h/2 - 2; i < h/2+3; i++ {
		termbox.SetCell(w/2+9, i, gameover2Body, gameover2FgColor, gameover2BgColor)
		termbox.SetCell(w/2-9, i, gameover2Body, gameover2FgColor, gameover2BgColor)
	}
	termbox.SetCell(w/2+9, h/2+3, gameoverCornerBody, gameover2FgColor, gameover2BgColor)
	termbox.SetCell(w/2+9, h/2-3, gameoverCornerBody, gameover2FgColor, gameover2BgColor)
	termbox.SetCell(w/2-9, h/2+3, gameoverCornerBody, gameover2FgColor, gameover2BgColor)
	termbox.SetCell(w/2-9, h/2-3, gameoverCornerBody, gameover2FgColor, gameover2BgColor)

	writeText(w/2-4, h/2, "GAME OVER", gameover1FgColor, gameover1BgColor)
	score := fmt.Sprintf("%s%d", "YOUR SCORE: ", g.score)
	writeText(w/2-7, h/2+1, score, gameover1FgColor, gameover1BgColor)
}

// Redraws the terminal.
func draw(g game) {
	// Clear the old "frame".
	termbox.Clear(snakeFgColor, snakeBgColor)
	drawSnakePosition(g)
	drawSnake(g.sn)
	drawApllePosition(g)
	drawApple(g.ap)
	drawscore(g)
	drawborder()
	// Update the "frame".
	termbox.Flush()
}

// mod is a custom implementation of the '%' (modulo) operator that always
// returns positive numbers.
// func mod(n, m int) int {
// 	return (n%m + m) % m
// }

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

func moveLeft(g game) game {
	if g.v.x != 1 {
		g.v = coord{-1, 0}
	}
	return g
}
func moveRight(g game) game {
	if g.v.x != -1 {
		g.v = coord{1, 0}
	}
	return g
}
func moveUp(g game) game {
	if g.v.y != 1 {
		g.v = coord{0, -1}
	}
	return g
}
func moveDown(g game) game {
	if g.v.y != -1 {
		g.v = coord{0, 1}
	}
	return g
}

func aplleEaten(g game) (apple, []coord, int) {
	w, h := termbox.Size()
	if g.ap.pos.x == g.sn.pos[0].x && g.ap.pos.y == g.sn.pos[0].y {
		g.ap = newAplle(w, h)
		g.sn.pos = append(g.sn.pos, coord{g.sn.pos[len(g.sn.pos)-1].x - g.v.x, g.sn.pos[len(g.sn.pos)-1].y - g.v.y})
		g.score++
	}

	return g.ap, g.sn.pos, g.score
}
func appleborder(g game) apple {
	w, h := termbox.Size()
	if g.ap.pos.x == 0 || g.ap.pos.y == 0 || g.ap.pos.x == w-1 || g.ap.pos.y == h-1 {
		g.ap = newAplle(w, h)
	}
	return g.ap
}
func snakehitsnake(g game) bool {
	tmp := 0
	for _, v := range g.sn.pos {
		if tmp == 1 {
			if v.x == g.sn.pos[0].x && v.y == g.sn.pos[0].y {
				return true
			}
		}

		tmp = 1
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
			if borderCrash(g) || snakehitsnake(g) {
				gameover(g)
				termbox.Flush()
				time.Sleep(3 * time.Second)
				return

			}
			g = step(g)
			g.ap, g.sn.pos, g.score = aplleEaten(g)
			g.ap = appleborder(g)

		}
	}
}
