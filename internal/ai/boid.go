package ai

import (
	"image/color"
	"math"
	"math/rand"

	// anonymous import for png decoder
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jtbonhomme/asteboids/internal/fonts"
	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/jtbonhomme/asteboids/internal/vector"
	"github.com/sirupsen/logrus"
)

const (
	boidAcceleration float64 = 4.0
	boidMaxVelocity  float64 = 5.0
	separationFactor float64 = 1.0
	cohesionFactor   float64 = 1.0
	alignmentFactor  float64 = 1.5
)

// Boid is a PhysicalBody agent.
// It represents a single autonomous agent.
type Boid struct {
	physics.Body
}

// NewBoid creates a new Boid (PhysicalBody agent)
func NewBoid(
	log *logrus.Logger,
	x, y,
	screenWidth, screenHeight float64,
	boidImage *ebiten.Image,
	vision physics.AgentVision,
	debug bool) *Boid {
	b := Boid{}
	b.AgentType = physics.BoidAgent

	b.Orientation = math.Pi / 32 * float64(rand.Intn(64))

	b.Init(vector.Vector2D{
		X: boidMaxVelocity * math.Cos(b.Orientation),
		Y: boidMaxVelocity * math.Sin(b.Orientation),
	})
	b.LimitVelocity(boidMaxVelocity)
	b.Log = log

	b.Move(vector.Vector2D{
		X: x,
		Y: y,
	})
	b.PhysicWidth = 10
	b.PhysicHeight = 10
	b.ScreenWidth = screenWidth
	b.ScreenHeight = screenHeight

	b.Image = boidImage
	b.Vision = vision
	b.Debug = debug
	return &b
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (b *Boid) Update() {
	defer b.Body.UpdatePosition()

	acceleration := vector.Vector2D{}
	nearestAgent := b.Vision(b.Position().X, b.Position().Y, 400.0)

	cohesion := b.cohesion(nearestAgent)
	cohesion.Multiply(cohesionFactor)
	acceleration.Add(cohesion)

	separation := b.separate(nearestAgent)
	separation.Multiply(separationFactor)
	acceleration.Add(separation)

	alignment := b.align(nearestAgent)
	alignment.Multiply(alignmentFactor)
	acceleration.Add(alignment)

	b.Accelerate(acceleration)
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (b *Boid) Draw(screen *ebiten.Image) {
	defer b.Body.Draw(screen)

	if b.Debug {
		msg := b.String()
		textDim := text.BoundString(fonts.MonoSansRegularFont, msg)
		textWidth := textDim.Max.X - textDim.Min.X
		text.Draw(screen, msg, fonts.MonoSansRegularFont, int(b.Position().X)-textWidth/2, int(b.Position().Y+b.PhysicHeight/2+5), color.Gray16{0x999f})
	}
}

// cohesion returns the force imposed by flocking cohesion rule.
func (b *Boid) cohesion(flock []physics.Physic) vector.Vector2D {
	return vector.Vector2D{
		X: 0,
		Y: 0,
	}
}

// separate returns the force imposed by flocking separation rule.
func (b *Boid) separate(flock []physics.Physic) vector.Vector2D {
	return vector.Vector2D{
		X: 0,
		Y: 0,
	}
}

// align returns the force imposed by flocking alignment rule.
func (b *Boid) align(flock []physics.Physic) vector.Vector2D {
	return vector.Vector2D{
		X: 0,
		Y: 0,
	}
}
