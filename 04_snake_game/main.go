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
	snakeTail    = '#'
	appleicon    = '@'
	snakeFgColor = termbox.ColorRed
	snakeBgColor = termbox.ColorDefault
	borderColor  = termbox.ColorWhite
)

func writeText(x, y int, s string, fg, bg termbox.Attribute) {
	for i, ch := range s {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

type coord struct {
	x, y int
}

type snake struct {
	pos coord
}
type apple struct {
	pos coord
}

type game struct {
	score                   int
	s                       snake
	a                       apple
	dir                     coord
	fieldWidth, fieldHeight int
}

func newSnake(maxX, maxY int) snake {
	return snake{coord{rand.Intn(maxX - 1), rand.Intn(maxY - 1)}}
}

func newGame() game {
	w, h := termbox.Size()
	return game{
		fieldWidth:  w,
		fieldHeight: h,
		s:           newSnake(w, h),
		a:           newApple(w, h),
		dir:         coord{1, 0},
		score:       0,
	}
}
func drawBorder(fw, fh int) {
	for i := 0; i <= fw; i++ {
		termbox.SetBg(i, 0, borderColor)
	}
	for i := 0; i <= fh; i++ {
		termbox.SetBg(fw-1, i, borderColor)
	}
	for i := fw; i >= 0; i-- {
		termbox.SetBg(i, fh-1, borderColor)
	}
	for i := fh; i >= 0; i-- {
		termbox.SetBg(0, i, borderColor)
	}
}

func drawScore(g game) {
	str := fmt.Sprintf("Score: %v", g.score)
	writeText(0, 0, str, termbox.ColorRed, borderColor)
}
func drawSnakeCord(g game) {
	str := fmt.Sprintf("(%d, %d)", g.s.pos.x, g.s.pos.y)
	writeText(g.fieldWidth-len(str), 0, str, snakeFgColor, borderColor)
}

func drawSnake(s snake) {
	termbox.SetCell(s.pos.x, s.pos.y, snakeBody, snakeFgColor, snakeBgColor)
}
func drawTail(s snake) {
	termbox.SetCell(s.pos.x, s.pos.y, snakeTail, snakeFgColor, snakeBgColor)
}
func newApple(fw, fh int) apple {
	apx := rand.Intn(fw - 1)
	apy := rand.Intn(fh - 1)
	return apple{coord{apx, apy}}
}
func drawApple(g game) {
	termbox.SetCell(g.a.pos.x, g.a.pos.y, appleicon, termbox.ColorRed, termbox.ColorDefault)
	termbox.Flush()
}

func draw(g game) {
	termbox.Clear(snakeFgColor, snakeBgColor)
	drawBorder(g.fieldWidth, g.fieldHeight)
	drawSnakeCord(g)
	drawScore(g)
	drawApple(g)
	drawSnake(g.s)
	drawTail(g.s)
	termbox.Flush()
}

func mod(n, m int) int {
	return (n%m + m) % m
}

func moveSnake(s snake, dir coord, fw, fh int) snake {
	s.pos.x = mod(s.pos.x+dir.x, fw)
	s.pos.y = mod(s.pos.y+dir.y, fh)
	return s
}

func step(g game) game {
	g.s = moveSnake(g.s, g.dir, g.fieldWidth, g.fieldHeight)
	draw(g)
	return g
}
func checkb(g game) bool {
	if g.s.pos.x == 0 || g.s.pos.x == g.fieldWidth-1 || g.s.pos.y == 0 || g.s.pos.y == g.fieldHeight-1 {
		return true
	}
	return false
}
func checka(g game) bool {
	if g.s.pos.x == g.a.pos.x && g.s.pos.y == g.a.pos.y {
		return true
	}
	return false
}

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatalf("failed to init termbox: %v", err)
	}
	defer termbox.Close()

	rand.Seed(time.Now().UnixNano())
	g := newGame()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(90 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyArrowDown:
					g.dir = coord{0, 1}
				case termbox.KeyArrowUp:
					g.dir = coord{0, -1}
				case termbox.KeyArrowLeft:
					g.dir = coord{-1, 0}
				case termbox.KeyArrowRight:
					g.dir = coord{1, 0}
				case termbox.KeyEsc:
					return
				}
			}
		case <-ticker.C:
			if checkb(g) {
				termbox.Clear(termbox.ColorWhite, termbox.ColorBlue)
				writeText(70, 5, `Game Over`, termbox.ColorWhite, termbox.ColorBlue)
				termbox.Flush()
				time.Sleep(2 * time.Second)
				return
			}
			if checka(g) {
				g.score++
				newApple(g.fieldWidth, g.fieldHeight)
			}
			g = step(g)
		}
	}
}
