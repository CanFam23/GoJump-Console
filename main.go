package main

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

var floorStyle = tcell.StyleDefault.
	Foreground(color.White).
	Background(color.Purple)

var playerStyle = tcell.StyleDefault.Foreground(color.White).Background(color.Red)

// drawBox renders a bordered rectangle with a filled background.
// Adapted from the tcell tutorial drawBox example:
// https://github.com/gdamore/tcell/blob/main/TUTORIAL.md
func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.Put(col, row, " ", style)
		}
	}

	for col := x1; col <= x2; col++ {
		s.Put(col, y1, string(tcell.RuneHLine), style)
		s.Put(col, y2, string(tcell.RuneHLine), style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.Put(x1, row, string(tcell.RuneVLine), style)
		s.Put(x2, row, string(tcell.RuneVLine), style)
	}

	if y1 != y2 && x1 != x2 {
		s.Put(x1, y1, string(tcell.RuneULCorner), style)
		s.Put(x2, y1, string(tcell.RuneURCorner), style)
		s.Put(x1, y2, string(tcell.RuneLLCorner), style)
		s.Put(x2, y2, string(tcell.RuneLRCorner), style)
	}
}

// drawLine writes a single line of text starting at x,y.
func drawLine(s tcell.Screen, x, y int, style tcell.Style, text string) {
	for _, r := range text {
		s.SetContent(x, y, r, nil, style)
		x++
	}
}

// drawHome renders the instructions screen before gameplay starts.
func drawHome(s tcell.Screen, style tcell.Style) {
	w, h := s.Size()
	lines := []string{
		"GO JUMP CONSOLE",
		"",
		"How to play:",
		"J: Jump",
		"D: Move right",
		"Avoid the walls",
		"",
		"Optimal terminal size: 75x25",
		"Changing terminal size changes gameplay and may break it.",
		"(Fixed terminal size cannot be enforced here.)",
		"",
		"S: Start game",
		"Esc or Ctrl+C: Quit",
	}

	s.Clear()
	startY := max((h/2)-(len(lines)/2), 0)

	for i, line := range lines {
		y := startY + i
		if y >= h {
			break
		}
		x := (w - len(line)) / 2
		if x < 0 {
			x = 0
		}
		drawLine(s, x, y, style, line)
	}

	s.Show()
}

// main initializes tcell and runs the event/update loop.
// Screen setup and event-loop structure are adapted from:
// https://github.com/gdamore/tcell/blob/main/TUTORIAL.md
func main() {
	defStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

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

	// Adapted panic-safe shutdown pattern from the tcell tutorial demo.
	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	w, h := s.Size()

	game := NewGame(w, h)
	showHome := true

	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	for {
		select {
		case ev := <-s.EventQ():
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
				w, h = s.Size()
				game.OnResize(w, h)
				if showHome {
					drawHome(s, defStyle)
				}
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					return
				} else if ev.Key() == tcell.KeyCtrlL {
					s.Sync()
					if showHome {
						drawHome(s, defStyle)
					}
				} else if showHome {
					if ev.Str() == "S" || ev.Str() == "s" {
						showHome = false
						game.Reset()
					}
				} else if ev.Str() == "C" || ev.Str() == "c" {
					s.Clear()
				} else if ev.Str() == "J" || ev.Str() == "j" {
					game.player.InitJump()
				} else if ev.Str() == "D" || ev.Str() == "d" {
					game.moveRight = true
				} else {
					s.Clear()
					s.Put(0, 0, ev.Str(), defStyle)
				}
			}
		case <-ticker.C:
			if showHome {
				drawHome(s, defStyle)
			} else {
				Update(s, game)
			}
		}
	}
}

// Update advances and redraws one gameplay frame.
func Update(s tcell.Screen, g *Game) {
	g.Update()
	s.Clear()
	g.Draw(s)
	s.Show()
}
