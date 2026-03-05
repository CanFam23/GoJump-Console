package main

import (
	"github.com/gdamore/tcell/v3"
)

const playerWidth = 6
const playerHeight = 4
const floorHeight = 3
const gameSpeed = -0.25

type Game struct {
	scrW   int
	scrH   int
	walls  []*Wall
	floor  *Wall
	player *Player
}

func NewGame(w, h int) *Game {
	game := &Game{
		scrW:   w,
		scrH:   h,
		floor:  NewWall(0, h-3, w, 3),
		player: NewPlayer(10, h-floorHeight-1, playerWidth, playerHeight),
	}

	for i := range 2 {
		game.walls = append(
			game.walls,
			NewWallDefault(float64(w+((w/2)*i)), h-floorHeight-1-10),
		)
	}

	return game
}

// drawBox(s, 10, h-pHeight-floorHeight-1, 10+pWidth, h-floorHeight-1, playerStyle, "p")

func (g *Game) OnResize(newScrW, newScrH int) {
	g.scrW = newScrW
	g.scrH = newScrH

	g.player.y = newScrH - floorHeight - 1

	g.floor.y = newScrH - 3

	g.floor.w = newScrW

	for _, w := range g.walls {
		w.y = g.scrH - floorHeight - 1 - w.h
	}
}

func (g *Game) Reset() {
	g.player.Reset(10, g.scrH-floorHeight-1)

	for i := range g.walls {
		g.walls[i].x = float64(g.scrW + ((g.scrW / 2) * i))
	}
}

func (g *Game) Draw(s tcell.Screen) {
	g.floor.Draw(s)
	g.player.Draw(s)

	for _, w := range g.walls {
		w.Draw(s)
	}
}

func (g *Game) Update() {
	groundY := g.scrH - floorHeight - 1
	g.player.Update(groundY, gameSpeed, g.scrW)

	// move walls left only (no dy)
	for _, wall := range g.walls {
		if wall.IsColliding(g.player) {
			g.Reset()
		}

		if int(wall.x)+wall.w < 0 {
			wall.x = float64(g.scrW) + float64(wall.w)
		}
		wall.Update(gameSpeed, 0)
	}
}
