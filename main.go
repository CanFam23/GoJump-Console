package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

var floorStyle = tcell.StyleDefault.
	Foreground(color.White).
	Background(color.Purple)

var playerStyle = tcell.StyleDefault.Foreground(color.White).Background(color.Red)

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	var width int
	for text != "" {
		text, width = s.Put(col, row, text, style)
		col += width
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
		if width == 0 {
			// incomplete grapheme at end of string
			break
		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.Put(col, row, " ", style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.Put(col, y1, string(tcell.RuneHLine), style)
		s.Put(col, y2, string(tcell.RuneHLine), style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.Put(x1, row, string(tcell.RuneVLine), style)
		s.Put(x2, row, string(tcell.RuneVLine), style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.Put(x1, y1, string(tcell.RuneULCorner), style)
		s.Put(x2, y1, string(tcell.RuneURCorner), style)
		s.Put(x1, y2, string(tcell.RuneLLCorner), style)
		s.Put(x2, y2, string(tcell.RuneLRCorner), style)
	}

	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
}

func main() {
	defStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()
	s.SetSize(75, 25)

	// Draw initial boxes
	// drawBox(s, 1, 1, 42, 7, boxStyle, "Click and drag to draw a box")
	// drawBox(s, 5, 9, 32, 14, boxStyle, "Press C to reset")

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	w, h := s.Size()

	game := NewGame(w, h)

	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	for {
		select {
		case ev := <-s.EventQ():
			// Process event
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
				w, h = s.Size()

				sXStr := strconv.Itoa(w)
				sYStr := strconv.Itoa(h)

				s.Clear()

				drawText(s, 0, 0, len(sXStr), 0, defStyle, sXStr)
				drawText(s, 0, 1, len(sYStr), 1, defStyle, sYStr)

				game.OnResize(w, h)
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					return
				} else if ev.Key() == tcell.KeyCtrlL {
					s.Sync()
				} else if ev.Str() == "C" || ev.Str() == "c" {
					s.Clear()
				} else if ev.Str() == "J" || ev.Str() == "j" {
					game.player.InitJump()
				} else if ev.Str() == "D" || ev.Str() == "d" {
					game.player.x += 4
				} else {
					s.Clear()
					s.Put(0, 0, ev.Str(), defStyle)
				}
			}
		case <-ticker.C:
			Update(s, *game)
		}
	}
}

func Update(s tcell.Screen, g Game) {
	g.Update()

	// Player
	s.Clear()

	g.Draw(s)
	// Update screen
	s.Show()
}
