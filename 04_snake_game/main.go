package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	gameOvertext = `________  ________  _____ ______   _______           ________  ___      ___ _______   ________     
|\   ____\|\   __  \|\   _ \  _   \|\  ___ \         |\   __  \|\  \    /  /|\  ___ \ |\   __  \    
\ \  \___|\ \  \|\  \ \  \\\__\ \  \ \   __/|        \ \  \|\  \ \  \  /  / | \   __/|\ \  \|\  \   
 \ \  \  __\ \   __  \ \  \\|__| \  \ \  \_|/__       \ \  \\\  \ \  \/  / / \ \  \_|/_\ \   _  _\  
  \ \  \|\  \ \  \ \  \ \  \    \ \  \ \  \_|\ \       \ \  \\\  \ \    / /   \ \  \_|\ \ \  \\  \| 
   \ \_______\ \__\ \__\ \__\    \ \__\ \_______\       \ \_______\ \__/ /     \ \_______\ \__\\ _\ 
    \|_______|\|__|\|__|\|__|     \|__|\|_______|        \|_______|\|__|/       \|_______|\|__|\|__|`
	apple        = 'A'
	snakeBody    = '*'
	snakeFgColor = termbox.ColorRed
	// Use the default background color for the snake.
	snakeBgColor = termbox.ColorDefault
)

// writeText writes a string to the buffer.
// func writeText(x, y int, s string, fg, bg termbox.Attribute) {
// 	for i, ch := range s {
// 		termbox.SetCell(x+i, y, ch, fg, bg)
// 	}
// }

// coord is a coordinate on a plane.
type coord struct {
	x, y int
}

// snake is a struct with fields representing a snake.
type snake struct {
	// Position of a snake.
	body []coord
	// Movement direction of a snake.
	dir coord
}

// game represents a state of the game.
type game struct {
	sn    snake
	apple coord
	// Game field dimensions.
	fieldWidth, fieldHeight int
	score                   int
}

func printText(c coord, text string) {
	i := 0
	for _, v := range text {
		if v == '\n' {
			c.y++
			i = 0
			continue
		}
		termbox.SetCell(c.x+i, c.y, v, snakeFgColor, snakeBgColor)
		i++
	}

}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

// newSnake returns a new struct instance representing a snake.
// The snake is placed in a random position in the game field.
// The movement direction is right.
func newSnake(maxX, maxY int) snake {
	// rand.Intn generates a pseudo-random number:
	// https://pkg.go.dev/math/rand#Intn
	x, y := random(3, maxX), random(0, maxY)
	return snake{
		body: []coord{{x, y}, {x - 1, y}, {x - 2, y}, {x - 3, y}},
		dir:  coord{-1, 0},
	}
}

// newGame returns a new game state.
func newGame() game {
	// Sets game field dimensions to the size of the terminal.
	w, h := termbox.Size()
	return game{
		sn:          newSnake(w, h),
		apple:       coord{rand.Intn(w), rand.Intn(h)},
		fieldWidth:  w,
		fieldHeight: h,
	}
}

// drawSnakePosition draws the current snake position (as a debugging
// information) in the buffer.
// func drawSnakePosition(g game) {
// 	str := fmt.Sprintf("(%d, %d)", g.sn.pos.x, g.sn.pos.y)
// 	writeText(g.fieldWidth-len(str), 0, str, snakeFgColor, snakeBgColor)
// }

// drawSnake draws the snake in the buffer.
func drawSnake(sn snake) {
	for i, v := range sn.body {
		if i != len(sn.body)-1 {
			termbox.SetCell(v.x, v.y, snakeBody, snakeFgColor, snakeBgColor)
		}
	}
}

// Redraws the terminal.
func draw(g game) {
	// Clear the old "frame".
	termbox.Clear(snakeFgColor, snakeBgColor)
	// draw apple
	termbox.SetCell(g.apple.x, g.apple.y, apple, snakeFgColor, snakeBgColor)
	// drawSnakePosition(g)
	drawSnake(g.sn)
	// Update the "frame".
	termbox.Flush()
}

// mod is a custom implementation of the '%' (modulo) operator that always
// returns positive numbers.
func mod(n, m int) int {

	return (n%m + m) % m
}

// Makes a move for a snake. Returns a game with an updated position.
func moveSnake(g game) game {
	for i := range g.sn.body {
		if i != len(g.sn.body)-1 {
			g.sn.body[i] = g.sn.body[i+1]
		} else {
			g.sn.body[i].x = mod(g.sn.body[i].x+g.sn.dir.x, g.fieldWidth)
			g.sn.body[i].y = mod(g.sn.body[i].y+g.sn.dir.y, g.fieldHeight)
		}
	}

	return g
}

func collision(g game) bool {
	m := g.sn.body[0]
	for _, v := range g.sn.body[1 : len(g.sn.body)-1] {
		if v.x == m.x && v.y == m.y {
			return true
		}
	}
	return false
}

func GameOver(score int) {
	termbox.Clear(snakeFgColor, snakeBgColor)
	printText(coord{0, 0}, gameOvertext)
	printText(coord{46, 10}, "Your score: "+strconv.Itoa(score))
	termbox.Flush()
	time.Sleep(5 * time.Second)
}

func step(g game) game {
	w, h := termbox.Size()
	if g.apple == g.sn.body[len(g.sn.body)-1] {
		g.score++
		g.apple = coord{rand.Intn(w), rand.Intn(h)}
		g.sn.body = append([]coord{{g.sn.body[0].x, g.sn.body[0].y}}, g.sn.body...)
	}
	g = moveSnake(g)
	draw(g)
	return g
}

func moveLeft(g game) game {
	if (g.sn.dir != coord{1, 0}) {
		g.sn.dir = coord{-1, 0}
	}
	return g
}
func moveRight(g game) game {
	if (g.sn.dir != coord{-1, 0}) {
		g.sn.dir = coord{1, 0}
	}
	return g
}
func moveUp(g game) game {
	if (g.sn.dir != coord{0, 1}) {
		g.sn.dir = coord{0, -1}
	}
	return g
}
func moveDown(g game) game {
	if (g.sn.dir != coord{0, -1}) {
		g.sn.dir = coord{0, 1}
	}
	return g
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
			if collision(g) {
				GameOver(g.score)
				return
			}
			g = step(g)
		}
	}
}
