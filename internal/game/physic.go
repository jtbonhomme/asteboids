package game

import (
	"math"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

const (
	AccelerationFactor float64 = 0.3
	velocityFactor     float64 = 1.8
	maxVelocity        float64 = 5.5
	rotationAngle      float64 = math.Pi / 36 // rotation of 5°
	frictionFactor     float64 = 0.03
)

type Position struct {
	X int
	Y int
}

type Vector struct {
	X float64
	Y float64
}

type Physic interface {
	// Draw draws the agent on screen
	Draw(*ebiten.Image)
	// Update updates the agent state
	Update()
	// Init initializes the physic body
	Init()
	// ID displays physic body unique ID
	ID() string
}

// AgentUnregister is a function to unregister an agent
type AgentUnregister func(string)

type PhysicBody struct {
	Position

	id          uuid.UUID
	Log         *logrus.Logger
	Orientation float64 // theta (radian)
	Size        float64

	PhysicWidth  int
	PhysicHeight int
	ScreenWidth  int
	ScreenHeight int

	Velocity     Vector
	Acceleration Vector

	Unregister AgentUnregister
	Image      *ebiten.Image
}

// Init initializes the physic body
func (a *PhysicBody) Init() {
	a.id = uuid.New()
}

// Draw draws the agent.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (a *PhysicBody) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(-a.PhysicWidth/2), float64(-a.PhysicHeight/2))
	op.GeoM.Rotate(a.Orientation)
	op.GeoM.Translate(float64(a.X), float64(a.Y))
	screen.DrawImage(a.Image, op)
}

func (a *PhysicBody) Rotate(i float64) {
	a.Orientation += i * rotationAngle
	if a.Orientation > 2*math.Pi {
		a.Orientation -= 2 * math.Pi
	}
	if a.Orientation < 0 {
		a.Orientation += 2 * math.Pi
	}
}

func (a *PhysicBody) UpdateAcceleration(i float64) {
	a.Acceleration.X = AccelerationFactor * i * math.Cos(a.Orientation)
	a.Acceleration.Y = AccelerationFactor * i * math.Sin(a.Orientation)
}

func (a *PhysicBody) UpdateVelocity() {
	a.Velocity.X += a.Acceleration.X - frictionFactor*a.Velocity.X
	a.Velocity.Y += a.Acceleration.Y - frictionFactor*a.Velocity.Y

	velocityValue := math.Sqrt(a.Velocity.X*a.Velocity.X + a.Velocity.Y*a.Velocity.Y)
	if velocityValue > maxVelocity {
		a.Velocity.X = maxVelocity * math.Cos(a.Orientation)
		a.Velocity.Y = maxVelocity * math.Sin(a.Orientation)
	}
	if velocityValue < 0 {
		a.Velocity.X = 0
		a.Velocity.Y = 0
	}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (a *PhysicBody) Update() {
	a.UpdateVelocity()

	// update position
	a.X += int(velocityFactor * a.Velocity.X)
	a.Y += int(velocityFactor * a.Velocity.Y)

	if a.X > a.ScreenWidth {
		a.X = 0
	} else if a.X < 0 {
		a.X = a.ScreenWidth
	}
	if a.Y > a.ScreenHeight {
		a.Y = 0
	} else if a.Y < 0 {
		a.Y = a.ScreenHeight
	}
}

// ID displays physic body unique ID
func (a *PhysicBody) ID() string {
	return a.id.String()
}
