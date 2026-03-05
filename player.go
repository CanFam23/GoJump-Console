package main

import "github.com/gdamore/tcell/v3"

// Player stores position, size, and jump physics state.
type Player struct {
	x                 float64
	y, w, h           int
	isJumping         bool
	velocity, gravity float64
}

// NewPlayer creates a player with default jump physics.
func NewPlayer(x float64, y, w, h int) (p *Player) {
	return &Player{
		x:         x,
		y:         y,
		w:         w,
		h:         h,
		isJumping: false,
		velocity:  0.0,
		gravity:   0.5,
	}
}

// Draw renders the player as a rectangle.
func (p *Player) Draw(s tcell.Screen) {
	drawBox(s, int(p.x), p.y, int(p.x)+p.w, p.y-p.h, playerStyle)
}

// InitJump starts a jump when the player is grounded.
func (p *Player) InitJump() {
	if !p.isJumping {
		p.isJumping = true
		p.velocity = 5.0
	}
}

// Update applies jump and horizontal movement, then clamps to ground.
func (p *Player) Update(groundY int, gameSpeed float64, screenWidth int) {
	dx, dy := p.CalcDisplacement()

	// Terminal y increases downward, so jumping subtracts y.
	p.y -= int(dy)

	nextX := p.x + gameSpeed + float64(dx)

	if nextX > 0 && nextX < float64(screenWidth-p.w) {
		p.x = nextX
	}

	if p.y >= groundY {
		p.y = groundY
		p.isJumping = false
		p.Reset(0, 0)
	}
}

// CalcDisplacement returns the per-frame displacement from jump physics.
func (p *Player) CalcDisplacement() (int, float64) {
	if p.isJumping {
		p.velocity -= p.gravity
		return 1, p.velocity
	}
	return 0, 0.0
}

// Reset optionally resets position and always resets jump state.
func (p *Player) Reset(x, y int) {
	if x != 0 || y != 0 {
		p.x = float64(x)
		p.y = y
	}
	p.isJumping = false
	p.velocity = 5.0
}
