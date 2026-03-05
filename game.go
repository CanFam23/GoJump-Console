package main

import (
	"github.com/gdamore/tcell/v3"
)

const playerWidth = 6
const playerHeight = 4
const floorHeight = 3
const gameSpeed = -0.25

// Game holds runtime state for the current run.
type Game struct {
	scrW      int
	scrH      int
	walls     []*Wall
	floor     *Wall
	player    *Player
	moveRight bool
}

// NewGame constructs the initial game state.
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

// OnResize updates dimensions and repositions entities tied to screen height.
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

// Reset returns player and walls to their starting layout.
func (g *Game) Reset() {
	g.player.Reset(10, g.scrH-floorHeight-1)

	for i := range g.walls {
		g.walls[i].x = float64(g.scrW + ((g.scrW / 2) * i))
	}
}

// Draw renders the floor, player, and all walls.
func (g *Game) Draw(s tcell.Screen) {
	g.floor.Draw(s)
	g.player.Draw(s)

	for _, w := range g.walls {
		w.Draw(s)
	}
}

// Update advances one simulation tick.
func (g *Game) Update() {
	groundY := g.scrH - floorHeight - 1
	g.player.Update(groundY, gameSpeed, g.scrW)

	if g.moveRight {
		g.moveRightStep()
		g.moveRight = false
	}

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

// moveRightStep applies one right movement step and collision response.
func (g *Game) moveRightStep() {
	const step = 3.0

	g.player.x += step

	if g.player.x+float64(g.player.w) > float64(g.scrW) {
		g.player.x = float64(g.scrW - g.player.w)
	}

	// Keep collision behavior consistent with wall movement collisions.
	for _, wall := range g.walls {
		if wall.IsColliding(g.player) {
			g.Reset()
		}
	}
}
