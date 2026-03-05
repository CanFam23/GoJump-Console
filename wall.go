package main

import "github.com/gdamore/tcell/v3"

const defaultWallWidth = 2
const defaultWallHeight = 10

type Wall struct {
	w, h, y int
	x       float64
}

func NewWall(x float64, y, w, h int) *Wall {
	return &Wall{
		w: w,
		h: h,
		x: x,
		y: y,
	}
}

func NewWallDefault(x float64, y int) *Wall {
	return &Wall{
		w: defaultWallWidth,
		h: defaultWallHeight,
		x: x,
		y: y,
	}
}

func (w *Wall) Update(dx float64, dy int) {
	w.x += dx
	w.y += dy
}

func (w *Wall) Draw(s tcell.Screen) {
	drawBox(s, int(w.x), w.y, int(int(w.x)+w.w), w.y+w.h, floorStyle, "")
}

func (w *Wall) IsGroundFor(player Player) bool {
	return player.y+player.h > w.y
}

func (wall *Wall) IsColliding(p *Player) bool {
	// X overlap (float64)
	pLeft := p.x
	pRight := p.x + float64(p.w)

	wLeft := wall.x
	wRight := wall.x + float64(wall.w)

	// Y overlap (int). Player is drawn from (y-h) to y
	pTop := p.y - p.h
	pBottom := p.y

	wTop := wall.y
	wBottom := wall.y + wall.h

	return pRight > wLeft &&
		pLeft < wRight &&
		pBottom > wTop &&
		pTop < wBottom
}
