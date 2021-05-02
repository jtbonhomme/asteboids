package physics

import (

	// anonymous import for png decoder
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/vector"
)

const (
	defaultMaxVelocity float64 = 3.5
)

const (
	StarshipAgent string = "starship"
	AsteroidAgent string = "asteroid"
	RubbleAgent   string = "rubble"
	BulletAgent   string = "bullet"
	BoidAgent     string = "boid"
	AIAgent       string = "ai"
)

// Size represents height and width of a physical body.
type Size struct {
	H float64
	W float64
}

type Physic interface {
	// Draw draws the agent on screen.
	Draw(*ebiten.Image)
	// Update proceeds the agent state.
	Update()
	// Init initializes the physic body.
	Init(vector.Vector2D)
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
	// position returns physical body position.
	Position() vector.Vector2D
	// FuturePosition return position the physic body will be in z time iteration.
	FuturePosition(float64) vector.Vector2D
	// Dimension returns physical body dimension.
	Dimension() Size
	// Type returns physical body agent type as a string.
	Type() string
	// Explode proceeds the agent explosion and termination.
	Explode()
	// Velocity returns physical body velocity.
	Velocity() vector.Vector2D
	// Dump write out internal agent's state.
	Dump(*os.File) error
}

// AgentRegister is a function to register an agent.
type AgentRegister func(Physic)

// AgentUnregister is a function to unregister an agent.
type AgentUnregister func(string, string)

// Todo change float64, float64 parameter by a Position
// AgentVision is a function used by agents to "see" around them.
type AgentVision func(float64, float64, float64) []Physic
