package agents

import (
	"math"

	// anonymous import for png decoder
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/jtbonhomme/asteboids/internal/vector"
	"github.com/sirupsen/logrus"
)

const (
	bulletVelocity float64 = 20.0
	bulletTTL      int     = 30
)

// Bullet is a PhysicalBody agent
// It represents a bullet shot by a starship agent.
type Bullet struct {
	physics.Body
	lifespan int
}

// NewBullet creates a new Bullet (PhysicalBody agent)
func NewBullet(log *logrus.Logger,
	x, y float64,
	orientation float64,
	screenWidth, screenHeight float64,
	cb physics.AgentUnregister,
	bulletImage *ebiten.Image) *Bullet {
	b := Bullet{
		lifespan: bulletTTL,
	}
	b.AgentType = physics.BulletAgent
	b.Unregister = cb

	b.Orientation = orientation

	b.Init(vector.Vector2D{
		X: bulletVelocity * math.Cos(b.Orientation),
		Y: bulletVelocity * math.Sin(b.Orientation),
	})
	b.Log = log

	b.Move(vector.Vector2D{
		X: x,
		Y: y,
	})
	b.PhysicWidth = 16
	b.PhysicHeight = 16
	b.ScreenWidth = screenWidth
	b.ScreenHeight = screenHeight
	b.Image = bulletImage
	return &b
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
// Update maintains a TTL counter to limit live of bullets.
func (b *Bullet) Update() {
	defer b.Body.UpdatePosition()
	b.lifespan--
	if b.lifespan == 0 {
		b.SelfDestroy()
	}
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (b *Bullet) Draw(screen *ebiten.Image) {
	b.Body.Draw(screen)
}

// SelfDestroy removes the agent from the game
func (b *Bullet) SelfDestroy() {
	b.Unregister(b.ID(), b.Type())
}
