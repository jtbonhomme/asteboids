package agents

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
	rubbleSplit           int     = 3
	asteroidVelocity      float64 = 1.0
	asteroidRotationSpeed float64 = 0.05
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
		X: asteroidVelocity * math.Cos(a.Orientation),
		Y: asteroidVelocity * math.Sin(a.Orientation),
	})
	a.Log = log

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
// Update maintains a TTL counter to limit live of bullets.
func (a *Asteroid) Update() {
	defer a.Body.UpdatePosition()
	a.Rotate(asteroidRotationSpeed)
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (a *Asteroid) Draw(screen *ebiten.Image) {
	defer a.Body.Draw(screen)

	if a.Debug {
		msg := a.String()
		textDim := text.BoundString(fonts.MonoSansRegularFont, msg)
		textWidth := textDim.Max.X - textDim.Min.X
		text.Draw(screen, msg, fonts.MonoSansRegularFont, int(a.Position().X)-textWidth/2, int(a.Position().Y+a.PhysicHeight/2+5), color.Gray16{0x999f})
	}
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
