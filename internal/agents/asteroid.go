package agents

import (
	"math"
	"math/rand"

	// anonymous import for png decoder
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/jtbonhomme/asteboids/internal/vector"
	"github.com/sirupsen/logrus"
)

const (
	rubbleSplit           int     = 3
	asteroidMaxVelocity   float64 = 0.8
	asteroidRotationSpeed float64 = 0.02
)

// Asteroid is a PhysicalBody agent
// It represents a bullet shot by a starship agent.
type Asteroid struct {
	physics.Body
	rubbleImages []*ebiten.Image
}

// NewAsteroid creates a new Asteroid (PhysicalBody agent)
func NewAsteroid(
	log *logrus.Logger,
	x, y,
	screenWidth, screenHeight float64,
	cbr physics.AgentRegister,
	cbu physics.AgentUnregister,
	asteroidImage *ebiten.Image,
	rubbleImages []*ebiten.Image,
	debug bool) *Asteroid {
	a := Asteroid{}
	a.AgentType = physics.AsteroidAgent
	a.Register = cbr
	a.Unregister = cbu

	a.Orientation = math.Pi / 16 * float64(rand.Intn(32))

	a.Init(vector.Vector2D{
		X: asteroidMaxVelocity * math.Cos(a.Orientation),
		Y: asteroidMaxVelocity * math.Sin(a.Orientation),
	})
	a.Log = log
	a.LimitVelocity(asteroidMaxVelocity)

	a.Move(vector.Vector2D{
		X: x,
		Y: y,
	})
	a.PhysicWidth = 100
	a.PhysicHeight = 100
	a.ScreenWidth = screenWidth
	a.ScreenHeight = screenHeight

	a.Image = asteroidImage
	a.rubbleImages = rubbleImages
	a.Debug = debug
	return &a
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (a *Asteroid) Update() {
	defer a.Body.Update()
	//defer a.Body.UpdatePosition()
	a.Rotate(asteroidRotationSpeed)
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (a *Asteroid) Draw(screen *ebiten.Image) {
	defer a.Body.Draw(screen)
}

// Explode proceeds the asteroid explosion and termination.
// For the first explosion, it splits into smaller asteroids (e.g. rubble).
// For the second explosion, each rubble disapear from game.
func (a *Asteroid) Explode() {
	defer a.Unregister(a.ID(), a.Type())

	for i := 0; i < rubbleSplit; i++ {
		rubble := NewRubble(a.Log,
			a.Position().X,
			a.Position().Y,
			a.ScreenWidth, a.ScreenHeight,
			a.Unregister,
			a.rubbleImages[rand.Intn(5)],
			a.Debug)
		a.Register(rubble)
	}
}
