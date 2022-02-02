package main

import (
	"log"
	"time"

	"github.com/nsf/termbox-go"
)

// writeText writes a string to the buffer.
func writeText(x, y int, s string, fg, bg termbox.Attribute) {
	i := 0
	for _, ch := range s {
		termbox.SetCell(x+i, y, ch, fg, bg)
		i++
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatalf("failed to init termbox: %v", err)
	}
	defer termbox.Close()

	// Set a rune with a given foreground/background color in the buffer.
	termbox.SetCell(0, 0, 'A', termbox.ColorBlue, termbox.ColorCyan)
	// Runes can be formatted with bold, underline or cursive by combining a
	// color and the formatting attribute using '|'.
	termbox.SetCell(5, 0, 'X', termbox.ColorGreen|termbox.AttrBold, termbox.ColorLightMagenta)
	// Writing a text to the buffer.
	writeText(10, 5, "Hello, world!", termbox.ColorLightRed|termbox.AttrUnderline, termbox.ColorLightGreen)
	// Synchronize the buffer with the terminal.
	termbox.Flush()
	// Wait a few seconds.
	time.Sleep(3 * time.Second)
	// Clear the buffer.
	termbox.Clear(termbox.ColorCyan, termbox.ColorBlue)
	termbox.Flush()
	time.Sleep(3 * time.Second)
}
