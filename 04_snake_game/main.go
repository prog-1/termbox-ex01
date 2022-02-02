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
	snakeTail    = 'o'
	appleIcon    = '🍎'
	snakeFgColor = termbox.ColorGreen
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
	pos []coord
}
type apple struct {
	pos coord
}

type game struct {
	x                       int
	score                   int
	s                       snake
	a                       apple
	dir                     coord
	fieldWidth, fieldHeight int
}

func newSnake(maxX, maxY int) snake {
	x := rand.Intn(maxX-1) + 1
	y := rand.Intn(maxY-1) + 1
	return snake{[]coord{{x, y}}}
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
		x:           1,
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
	str := fmt.Sprintf("(%d, %d)", g.s.pos[0].x, g.s.pos[0].y)
	writeText(g.fieldWidth-len(str), 0, str, termbox.ColorRed, borderColor)
}

func drawSnake(s snake) {
	for _, pos := range s.pos {
		termbox.SetCell(pos.x, pos.y, snakeBody, snakeFgColor, snakeBgColor)
	}
}

func newApple(fw, fh int) apple {
	return apple{coord{rand.Intn(fw-1) + 1, rand.Intn(fh-1) + 1}}
}
func newAppleCoord(g game) game {
	g.a.pos.x = rand.Intn(g.fieldWidth-1) + 1
	g.a.pos.y = rand.Intn(g.fieldHeight-1) + 1
	return g
}
func drawApple(g game) {
	termbox.SetCell(g.a.pos.x, g.a.pos.y, appleIcon, termbox.ColorRed, termbox.ColorDefault)
	termbox.Flush()
}
func drawTail(s snake) {
	termbox.SetCell(s.pos[0].x, s.pos[0].y, snakeTail, snakeFgColor, snakeBgColor)
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

func moveSnake(s snake, dir coord, fw, fh int) snake {
	copy(s.pos[1:], s.pos[:])
	s.pos[0] = coord{s.pos[0].x + dir.x, s.pos[0].y + dir.y}
	return s
}
func moveAndAdd(s snake, dir coord, fw, fh int) snake {
	var new coord = coord{s.pos[len(s.pos)-1].x, s.pos[len(s.pos)-1].y}
	copy(s.pos[1:], s.pos[:])
	s.pos[0] = coord{s.pos[0].x + dir.x, s.pos[0].y + dir.y}
	s.pos = append(s.pos, new)
	return s
}

func step(g game) game {
	if g.score == g.x {
		g.x++
		g.s = moveAndAdd(g.s, g.dir, g.fieldWidth, g.fieldHeight)
	} else {
		g.s = moveSnake(g.s, g.dir, g.fieldWidth, g.fieldHeight)
	}
	draw(g)
	return g
}
func checkb(g game) bool {
	if g.s.pos[0].x == 0 || g.s.pos[0].x == g.fieldWidth-1 || g.s.pos[0].y == 0 || g.s.pos[0].y == g.fieldHeight-1 {
		return true
	}
	return false
}
func checka(g game) bool {
	if g.s.pos[0].x == g.a.pos.x && g.s.pos[0].y == g.a.pos.y {
		return true
	}
	return false
}
func checkSnakeCoord(g game) bool {
	for i := 1; i != len(g.s.pos)-1; i++ {
		if g.s.pos[0] == g.s.pos[i] {
			return true
		}
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
					if g.dir.y == -1 {
						g.dir = coord{0, -1}
					} else {
						g.dir = coord{0, 1}
					}
				case termbox.KeyArrowUp:
					if g.dir.y == 1 {
						g.dir = coord{0, 1}
					} else {
						g.dir = coord{0, -1}
					}
				case termbox.KeyArrowLeft:
					if g.dir.x == 1 {
						g.dir = coord{1, 0}
					} else {
						g.dir = coord{-1, 0}
					}
				case termbox.KeyArrowRight:
					if g.dir.x == -1 {
						g.dir = coord{-1, 0}
					} else {
						g.dir = coord{1, 0}
					}
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
			if g.score > 3 {
				if checkSnakeCoord(g) {
					termbox.Clear(termbox.ColorWhite, termbox.ColorBlue)
					writeText(70, 5, `Game Over`, termbox.ColorWhite, termbox.ColorBlue)
					termbox.Flush()
					time.Sleep(2 * time.Second)
					return
				}
			}
			if checka(g) {
				g.score++
				g = newAppleCoord(g)
			}

			g = step(g)
		}
	}
}
