package agents

import (
	"time"

	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/jtbonhomme/asteboids/internal/vector"
	"github.com/sirupsen/logrus"
)

const (
	bulletThrottle       time.Duration = 200 * time.Millisecond
	rotationAngle        float64       = math.Pi / 36 // rotation of 5°
	starshipMaxVelocity  float64       = 3.0
	starshipAcceleration float64       = 0.2
)

// Starship is a PhysicalBody agent.
// It represents a playable star ship.
type Starship struct {
	physics.Body
	lastBulletTime time.Time
	bulletImage    *ebiten.Image
}

// NewStarship creates a new Starship (PhysicalBody agent)
func NewStarship(
	log *logrus.Logger,
	x, y,
	screenWidth, screenHeight float64,
	cbr physics.AgentRegister,
	cbu physics.AgentUnregister,
	vision physics.AgentVision,
	starshipImage *ebiten.Image,
	bulletImage *ebiten.Image,
	debug bool) *Starship {
	s := Starship{
		lastBulletTime: time.Now(),
	}
	s.AgentType = physics.StarshipAgent
	s.Register = cbr
	s.Unregister = cbu
	s.Vision = vision
	s.Init(vector.Vector2D{
		X: 0,
		Y: 0,
	})
	s.LimitVelocity(starshipMaxVelocity)
	s.Orientation = math.Pi / 2
	s.Move(vector.Vector2D{
		X: x,
		Y: y,
	})
	s.Accelerate(vector.Vector2D{
		X: 0,
		Y: 0,
	})
	s.PhysicWidth = 50
	s.PhysicHeight = 50
	s.ScreenWidth = screenWidth
	s.ScreenHeight = screenHeight
	s.Log = log

	s.Image = starshipImage
	s.bulletImage = bulletImage

	s.Debug = debug
	return &s
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (s *Starship) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		s.Rotate(-rotationAngle)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		s.Rotate(rotationAngle)
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		acceleration := vector.Vector2D{
			X: math.Cos(s.Orientation),
			Y: math.Sin(s.Orientation),
		}
		acceleration.Multiply(starshipAcceleration)
		s.Accelerate(acceleration)
	} else {
		s.Accelerate(vector.Vector2D{})
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.SelfDestroy()
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		s.Shot()
	}

	s.UpdateVelocity()
	s.UpdatePosition()
}

// Shot adds a new bullet to the game.
func (s *Starship) Shot() {
	// throtlle call to avoid continuous shooting
	if time.Since(s.lastBulletTime) < bulletThrottle {
		return
	}
	s.lastBulletTime = time.Now()

	bullet := NewBullet(s.Log,
		s.Position().X, s.Position().Y,
		s.Orientation,
		s.ScreenWidth,
		s.ScreenHeight,
		s.Unregister,
		s.bulletImage)
	s.Register(bullet)
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (s *Starship) Draw(screen *ebiten.Image) {
	defer s.Body.Draw(screen)
	nearestAgent := s.Vision(s.Position().X, s.Position().Y)
	s.LinkAgents(screen, nearestAgent, []string{physics.AsteroidAgent, physics.RubbleAgent})
}

// SelfDestroy removes the agent from the game
func (s *Starship) SelfDestroy() {
	s.Explode()
}
