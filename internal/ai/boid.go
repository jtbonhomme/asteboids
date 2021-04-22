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
	boidVelocity float64 = 5.0
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
		X: boidVelocity * math.Cos(b.Orientation),
		Y: boidVelocity * math.Sin(b.Orientation),
	})
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

func (b *Boid) Avoid(agents []physics.Physic, agentType string) float64 {
	newOrientation := b.Orientation + math.Pi/16

	return newOrientation
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
// Update maintains a TTL counter to limit live of bullets.
func (b *Boid) Update() {
	defer b.Body.UpdatePosition()
	nearestAgent := b.Vision(b.Position().X, b.Position().Y, 400.0)
	b.Orientation = b.Avoid(nearestAgent, physics.AsteroidAgent)

	b.Accelerate(vector.Vector2D{
		X: boidVelocity * math.Cos(b.Orientation),
		Y: boidVelocity * math.Sin(b.Orientation),
	})
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
