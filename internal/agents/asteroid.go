package agents

import (
	"fmt"
	"math"
	"math/rand"

	// anonymous import for png decoder
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/game"
	"github.com/sirupsen/logrus"
)

const (
	asteroidVelocity float64 = 2.0
	rotationSpeed    float64 = 0.05
)

// Asteroid is a PhysicalBody agent
// It represents a bullet shot by a starship agent.
type Asteroid struct {
	game.PhysicBody
}

// NewAsteroid creates a new Asteroid (PhysicalBody agent)
func NewAsteroid(log *logrus.Logger, screenWidth, screenHeight int, cb game.AgentUnregister) *Asteroid {
	a := Asteroid{}
	a.Type = "asteroid"
	a.Unregister = cb

	a.Init()
	a.Log = log

	a.Orientation = math.Pi / 8 * float64(rand.Intn(16))
	a.Velocity = game.Vector{
		X: asteroidVelocity * math.Cos(a.Orientation),
		Y: asteroidVelocity * math.Sin(a.Orientation),
	}
	a.Size = 3
	a.PhysicWidth = 100
	a.PhysicHeight = 100
	a.ScreenWidth = screenWidth
	a.ScreenHeight = screenHeight
	a.X = rand.Intn(screenWidth)
	a.Y = rand.Intn(screenHeight)

	err := a.LoadImage(fmt.Sprintf("./assets/asteroid%d.png", rand.Intn(5)))
	if err != nil {
		log.Errorf("error when loading image from file: %s", err.Error())
	}
	return &a
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
// Update maintains a TTL counter to limit live of bullets.
func (a *Asteroid) Update() {
	a.Orientation += rotationSpeed
	// update position
	a.X += int(a.Velocity.X)
	a.Y += int(a.Velocity.Y)

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

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (a *Asteroid) Draw(screen *ebiten.Image) {
	a.PhysicBody.Draw(screen)
}

// SelfDestroy removes the agent from the game
func (a *Asteroid) SelfDestroy() {
	a.Unregister(a.ID())
}
