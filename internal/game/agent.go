package game

import "github.com/hajimehoshi/ebiten/v2"

type Agent interface {
	// Draw draws the agent on screen
	Draw(*ebiten.Image)
	// Update updates the agent state
	Update()
	// Position returns the current agent position
	Position() Position
	// Velocity returns the current agent velocity
	Velocity() (Vector, float64)
	// Acceleration returns the current agent acceleration
	Acceleration() (Vector, float64)
	// Orientation returns the current agent direction
	Orientation() float64
}

type Position struct {
	X int
	Y int
}

type Vector struct {
	X float64
	Y float64
}
