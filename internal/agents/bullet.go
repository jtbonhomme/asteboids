package agents

import (
	"math"

	// anonymous import for png decoder
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/game"
	"github.com/sirupsen/logrus"
)

const (
	bulletVelocity float64 = 19.0
	bulletTTL      int     = 30
)

// Bullet is a PhysicalBody agent
// It represents a bullet shot by a starship agent.
type Bullet struct {
	game.PhysicBody
	ttl int
}

// NewBullet creates a new Bullet (PhysicalBody agent)
func NewBullet(log *logrus.Logger, x, y int, orientation float64, screenWidth, screenHeight int, cb game.AgentUnregister) *Bullet {
	b := Bullet{
		ttl: bulletTTL,
	}
	b.Type = "bullet"
	b.Unregister = cb

	b.Init()
	b.Log = log

	b.Orientation = orientation
	b.Velocity = game.Vector{
		X: bulletVelocity * math.Cos(b.Orientation),
		Y: bulletVelocity * math.Sin(b.Orientation),
	}
	b.Size = 3
	b.PhysicWidth = 16
	b.PhysicHeight = 16
	b.ScreenWidth = screenWidth
	b.ScreenHeight = screenHeight
	b.X = x
	b.Y = y

	err := b.LoadImage("./assets/bullet.png")
	if err != nil {
		log.Errorf("error when loading image from file: %s", err.Error())
	}
	return &b
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
// Update maintains a TTL counter to limit live of bullets.
func (b *Bullet) Update() {
	b.ttl--
	if b.ttl == 0 {
		b.SelfDestroy()
	}
	// update position
	b.X += int(b.Velocity.X)
	b.Y += int(b.Velocity.Y)

	if b.X > b.ScreenWidth {
		b.X = 0
	} else if b.X < 0 {
		b.X = b.ScreenWidth
	}
	if b.Y > b.ScreenHeight {
		b.Y = 0
	} else if b.Y < 0 {
		b.Y = b.ScreenHeight
	}
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (b *Bullet) Draw(screen *ebiten.Image) {
	b.PhysicBody.Draw(screen)
}

// SelfDestroy removes the agent from the game
func (b *Bullet) SelfDestroy() {
	b.Unregister(b.ID())
}
