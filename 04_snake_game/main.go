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
	pos []coord
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
func newSnake(g game) coord {
	// rand.Intn generates a pseudo-random number:
	// https://pkg.go.dev/math/rand#Intn
	// return snake{coord{rand.Intn(maxX), rand.Intn(maxY)}}
	g.sn.pos[0] = coord{rand.Intn(g.fieldWidth), rand.Intn(g.fieldHeight)}
	for {
		if g.sn.pos[0].x == 0 || g.sn.pos[0].x == g.fieldWidth-1 || g.sn.pos[0].y == 0 || g.sn.pos[0].y == 1 || g.sn.pos[0].y == g.fieldHeight-1 {
			g.sn.pos[0] = coord{rand.Intn(g.fieldWidth), rand.Intn(g.fieldHeight)}
		} else {
			break
		}
	}
	return g.sn.pos[0]
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
	str := fmt.Sprintf("(%d, %d)", g.sn.pos[0].x, g.sn.pos[0].y)
	writeText(g.fieldWidth-len(str), 0, str, snakeFgColor, snakeBgColor)
}

// drawSnake draws the snake in the buffer.
func drawSnake(sn snake) {
	for _, pos := range sn.pos {
		termbox.SetCell(pos.x, pos.y, snakeBody, snakeFgColor, snakeBgColor)
	}
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
	writeText(g.fieldWidth/2-6/2, 0, fmt.Sprint("Score:", g.a.score), termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
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
// func mod(n, m int) int {
// 	return (n%m + m) % m
// }

// Makes a move for a snake. Returns a snake with an updated position.
func moveSnake(s snake, v coord, fw, fh int) snake {
	// s.pos.x = mod(s.pos.x+v.x, fw)
	// s.pos.y = mod(s.pos.y+v.y, fh)
	copy(s.pos[1:], s.pos[:])
	s.pos[0] = coord{s.pos[0].x + v.x, s.pos[0].y + v.y}
	return s
}

func gameOver(g game, s string) {
	writeText(g.fieldWidth/2-9/2, g.fieldHeight/2, "Game Over", termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
	writeText(g.fieldWidth/2-len([]rune(s))/2, g.fieldHeight/2+1, s, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
	termbox.Flush()
	time.Sleep(5 * time.Second)
	termbox.Clear(snakeFgColor, snakeBgColor)
	termbox.Flush()
}

func collisions(g game) (isTerminated bool) {
	if g.sn.pos[0].x == 0 || g.sn.pos[0].x == g.fieldWidth-1 || g.sn.pos[0].y == 1 || g.sn.pos[0].y == g.fieldHeight-1 {
		gameOver(g, "I guess it wasn't a pleasant feeling")
		return true
	}
	for i := 3; i < len(g.sn.pos)-1; i++ {
		if g.sn.pos[0] == g.sn.pos[i] {
			gameOver(g, "Your tail doesn't look appetizing, does it?")
			return true
		}
	}
	return false
}

func step(g game) game {
	g.sn = moveSnake(g.sn, g.v, g.fieldWidth, g.fieldHeight)
	if g.sn.pos[0] == g.a.pos {
		g.a.score++
		g.a = newApple(g)
		g.sn.pos = append(g.sn.pos, coord{g.sn.pos[len(g.sn.pos)-1].x - g.sn.pos[len(g.sn.pos)-2].x, g.sn.pos[len(g.sn.pos)-1].y - g.sn.pos[len(g.sn.pos)-2].y})
	}
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

	// Initialize termbox.
	err := termbox.Init()
	if err != nil {
		log.Fatalf("failed to init termbox: %v", err)
	}

	// Other initialization.
	rand.Seed(time.Now().UnixNano())
	g := newGame()
	g.sn.pos = make([]coord, 2)
	g.sn.pos[0] = newSnake(g)
	g.a = newApple(g)
	prevKey := termbox.KeyArrowRight // Because the starting direction is always to the right(line 97)

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
					if prevKey != termbox.KeyArrowUp {
						g = moveDown(g)
						prevKey = termbox.KeyArrowDown
					}
				case termbox.KeyArrowUp:
					if prevKey != termbox.KeyArrowDown {
						g = moveUp(g)
						prevKey = termbox.KeyArrowUp
					}
				case termbox.KeyArrowLeft:
					if prevKey != termbox.KeyArrowRight {
						g = moveLeft(g)
						prevKey = termbox.KeyArrowLeft
					}
				case termbox.KeyArrowRight:
					if prevKey != termbox.KeyArrowLeft {
						g = moveRight(g)
						prevKey = termbox.KeyArrowRight
					}
				case termbox.KeyEsc:
					return
				}
			}
		case <-ticker.C:
			g = step(g)
			if collisions(g) {
				ticker.Stop()
				termbox.Close()
				os.Exit(0)
			}
			draw(g)
		}
	}
}
