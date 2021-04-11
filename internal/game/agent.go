package game

import "github.com/hajimehoshi/ebiten/v2"

type Agent interface {
	// Draw draws the agent on screen
	Draw(*ebiten.Image)
	// Update updates the agent state
	Update()
	// Position returns the current agent position
	Position() Position
	// Speed returns the current agent speed
	Speed() float64
	// Direction returns the current agent direction
	Direction() float64
}

type Position struct {
	X int
	Y int
}
