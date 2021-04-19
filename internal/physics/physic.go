package physics

import (

	// anonymous import for png decoder
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	accelerationFactor float64 = 0.3
	velocityFactor     float64 = 1.8
	maxVelocity        float64 = 5.5
	frictionFactor     float64 = 0.03
	collisionPrecision float64 = 1.0
)

const (
	StarshipAgent string = "starship"
	AsteroidAgent string = "asteroid"
	RubbleAgent   string = "rubble"
	BulletAgent   string = "bullet"
	BoidAgent     string = "boid"
)

// Size represents coordonnates (X, Y) of a physical body.
type Position struct {
	X float64
	Y float64
}

// Size represents height and width of a physical body.
type Size struct {
	H float64
	W float64
}

// Vector represents a vector composantes.
type Vector struct {
	X float64
	Y float64
}

// Block is a dimension and position helper structure.
type Block struct {
	Position
	Size
}

type Physic interface {
	// Draw draws the agent on screen.
	Draw(*ebiten.Image)
	// Update proceeds the agent state.
	Update()
	// Init initializes the physic body.
	Init()
	// ID displays physic body unique ID.
	ID() string
	// String displays physic body information as a string.
	String() string
	// Intersect returns true if the physical body collide another one.
	// Collision is computed based on Axis-Aligned Bounding Boxes.
	// https://developer.mozilla.org/en-US/docs/Games/Techniques/2D_collision_detection
	Intersect(Physic) bool
	// IntersectMultiple checks if multiple physical bodies are colliding with the first
	IntersectMultiple(map[string]Physic) (string, bool)
	// Dimensions returns physical body dimensions.
	Dimension() Block
	// Type returns physical body agent type as a string.
	Type() string
	// Explode proceeds the agent explosion and termination.
	Explode()
}

// AgentRegister is a function to register an agent.
type AgentRegister func(Physic)

// AgentUnregister is a function to unregister an agent.
type AgentUnregister func(string, string)

// AgentVision is a function used by agents to "see" around them.
type AgentVision func(float64, float64, float64) []Physic
