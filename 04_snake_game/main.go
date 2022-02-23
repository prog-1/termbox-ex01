package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	snakeBody    = '*'
	snakeFgColor = termbox.ColorGreen
	// Use the default background color for the snake.
	snakeBgColor = termbox.ColorDefault

	appleBody    = '🍏'
	appleFgColor = termbox.ColorYellow
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
	ap apple
	v  coord
	// Game field dimensions.
	fieldWidth, fieldHeight int
}

// newSnake returns a new struct instance representing a snake.
// The snake is placed in a random position in the game field.
// The movement direction is right.
func newSnake(g game) snake {
	// rand.Intn generates a pseudo-random number:
	// https://pkg.go.dev/math/rand#Intn
	x, y := rand.Intn(g.fieldWidth), rand.Intn(g.fieldHeight)

	g.sn.pos = []coord{{x, y}, {x - 1, y}, {x - 2, y}}
	for {
		if x == 0 || y == 0 {
			break
		}
		if g.sn.pos[0].x == 0 || g.sn.pos[0].x == g.fieldWidth-1 || g.sn.pos[0].y == 0 || g.sn.pos[0].y == 1 || g.sn.pos[0].y == g.fieldHeight-1 {
			g.sn.pos = []coord{{x, y}, {x - 1, y}, {x - 2, y}}
		} else {
			break
		}
	}
	return snake{g.sn.pos}
}

func newApple(g game) apple {
	// g.ap.pos.x, g.ap.pos.y = rand.Intn(g.fieldWidth), rand.Intn(g.fieldHeight)
	// return apple{coord{g.ap.pos.x, g.ap.pos.y}, g.ap.score}

	g.ap.pos = coord{rand.Intn(g.fieldWidth), rand.Intn(g.fieldHeight)}
	for {
		if g.ap.pos.x == 0 || g.ap.pos.x == g.fieldWidth-1 || g.ap.pos.y == 0 || g.ap.pos.y == 1 || g.ap.pos.y == g.fieldHeight-1 {
			g.ap.pos = coord{rand.Intn(g.fieldWidth), rand.Intn(g.fieldHeight)}
		} else {
			break
		}
	}
	return apple{g.ap.pos, g.ap.score}
}

// newGame returns a new game state.
func newGame() game {
	// Sets game field dimensions to the size of the terminal.
	w, h := termbox.Size()
	return game{
		fieldWidth:  w,
		fieldHeight: h,
		// sn:       newSnake(w, h),
		// ap:       newApple(w, h),
		v: coord{1, 0},
	}
}

// drawSnakePosition draws the current snake position (as a debugging
// information) in the buffer.
func drawSnakePosition(g game) {
	str := fmt.Sprintf("(%d, %d)", g.sn.pos[0].x, g.sn.pos[0].y)
	writeText(g.fieldWidth-len(str), 0, str, snakeFgColor, snakeBgColor)
}

func drawScore(g game) {
	str := fmt.Sprintf("Score: %d", g.ap.score)
	if g.ap.score < 5 { // the color will change depending on the score
		writeText((g.fieldWidth-len(str))/2, 0, str, termbox.ColorWhite, snakeBgColor)
	} else if g.ap.score >= 5 && g.ap.score < 10 {
		writeText((g.fieldWidth-len(str))/2, 0, str, termbox.ColorLightCyan, snakeBgColor)
	} else if g.ap.score >= 10 && g.ap.score < 20 {
		writeText((g.fieldWidth-len(str))/2, 0, str, termbox.ColorLightGreen, snakeBgColor)
	} else if g.ap.score >= 20 {
		writeText((g.fieldWidth-len(str))/2, 0, str, termbox.ColorLightMagenta, snakeBgColor)
	}
}

// drawSnake draws the snake in the buffer.
func drawSnake(sn snake) {
	// termbox.SetCell(sn.pos.x, sn.pos.y, snakeBody, snakeFgColor, snakeBgColor)

	for _, pos := range sn.pos {
		termbox.SetCell(pos.x, pos.y, snakeBody, snakeFgColor, snakeBgColor)
	}
}

func drawApple(ap apple) {
	termbox.SetCell(ap.pos.x, ap.pos.y, appleBody, appleFgColor, appleBgColor)
}

func drawBorders(g game) {
	for i := 0; i < g.fieldWidth; i++ {
		termbox.SetBg(i, 1, termbox.ColorLightRed)
		termbox.SetBg(i, g.fieldHeight-1, termbox.ColorLightCyan)
	}
	for j := 1; j < g.fieldHeight; j++ {
		termbox.SetBg(0, j, termbox.ColorLightMagenta)
		termbox.SetBg(g.fieldWidth-1, j, termbox.ColorLightGreen)
	}
	termbox.SetBg(0, 1, termbox.ColorLightYellow)
	termbox.SetBg(0, g.fieldHeight-1, termbox.ColorLightYellow)
	termbox.SetBg(g.fieldWidth-1, 1, termbox.ColorLightYellow)
	termbox.SetBg(g.fieldWidth-1, g.fieldHeight-1, termbox.ColorLightYellow)
}

// Redraws the terminal.
func draw(g game) {
	// Clear the old "frame".
	termbox.Clear(snakeFgColor, snakeBgColor)
	drawSnakePosition(g)
	drawScore(g)
	drawSnake(g.sn)
	drawApple(g.ap)
	drawBorders(g)
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

func HitsTheWalls(g game) bool {
	str := fmt.Sprintln("Game Over")
	if g.sn.pos[0].x == 0 || g.sn.pos[0].x == g.fieldWidth-1 || g.sn.pos[0].y == 1 || g.sn.pos[0].y == g.fieldHeight-1 {
		termbox.Clear(snakeFgColor, snakeBgColor)
		writeText((g.fieldWidth-len(str))/2, g.fieldHeight/2, str, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
		termbox.Flush()
		time.Sleep(5 * time.Second)
		termbox.Clear(snakeFgColor, snakeBgColor)
		termbox.Flush()
		return true
	}
	return false
}

func HitsItself(g game) bool {
	str := fmt.Sprintln("Game Over")
	for _, pos := range g.sn.pos[2:] {
		if g.sn.pos[0] == pos {
			termbox.Clear(snakeFgColor, snakeBgColor)
			writeText((g.fieldWidth-len(str))/2, g.fieldHeight/2, str, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
			termbox.Flush()
			time.Sleep(5 * time.Second)
			termbox.Clear(snakeFgColor, snakeBgColor)
			termbox.Flush()
			return true
		}
	}
	return false
}

func step(g game) game {
	g.sn = moveSnake(g.sn, g.v, g.fieldWidth, g.fieldHeight)

	if g.sn.pos[0] == g.ap.pos {
		g.ap.score++
		g.ap = newApple(g)
		g.sn.pos = append([]coord{{g.sn.pos[0].x, g.sn.pos[0].y}}, g.sn.pos...)
	}

	if HitsTheWalls(g) || HitsItself(g) {
		termbox.Close()
		os.Exit(0)
	}

	draw(g)
	return g
}

func moveLeft(g game) game {
	if !reflect.DeepEqual(g.v, coord{1, 0}) {
		g.v = coord{-1, 0}
	}
	return g
}

func moveRight(g game) game {
	if !reflect.DeepEqual(g.v, coord{-1, 0}) {
		g.v = coord{1, 0}
	}
	return g
}

func moveUp(g game) game {
	if !reflect.DeepEqual(g.v, coord{0, 1}) {
		g.v = coord{0, -1}
	}
	return g
}

func moveDown(g game) game {
	if !reflect.DeepEqual(g.v, coord{0, -1}) {
		g.v = coord{0, 1}
	}
	return g
}

func mainMenuDifficulty() (choice int) {
	fmt.Println(`Select difficulty:
1) Easy
2) Medium
3) Hard
4) Exit`)
	fmt.Scanln(&choice)
	return
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
	g.sn = newSnake(g)
	g.ap = newApple(g)

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	var ticker *time.Ticker
	choice := mainMenuDifficulty()
	if choice == 1 {
		ticker = time.NewTicker(120 * time.Millisecond)
	} else if choice == 2 {
		ticker = time.NewTicker(80 * time.Millisecond)
	} else if choice == 3 {
		ticker = time.NewTicker(50 * time.Millisecond)
	} else if choice == 4 {
		termbox.Close()
		os.Exit(0)
	} else {
		fmt.Println("ERROR: wrong choice")
	}
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
				case termbox.KeyEnter: // I also decided to add a 5 second pause
					str := fmt.Sprintln("Pause")
					writeText((g.fieldWidth-len(str))/2, g.fieldHeight/2, str, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault)
					timeleft := "5"
					writeText((g.fieldWidth-1)/2, (g.fieldHeight+2)/2, timeleft, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault)
					termbox.Flush()
					time.Sleep(1 * time.Second)
					timeleft = "4"
					writeText((g.fieldWidth-1)/2, (g.fieldHeight+2)/2, timeleft, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault)
					termbox.Flush()
					time.Sleep(1 * time.Second)
					timeleft = "3"
					writeText((g.fieldWidth-1)/2, (g.fieldHeight+2)/2, timeleft, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault)
					termbox.Flush()
					time.Sleep(1 * time.Second)
					timeleft = "2"
					writeText((g.fieldWidth-1)/2, (g.fieldHeight+2)/2, timeleft, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault)
					termbox.Flush()
					time.Sleep(1 * time.Second)
					timeleft = "1"
					writeText((g.fieldWidth-1)/2, (g.fieldHeight+2)/2, timeleft, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault)
					termbox.Flush()
					time.Sleep(1 * time.Second)
				}
			}
		case <-ticker.C:
			g = step(g)
		}
	}
}
