package agents

import (
	"image/color"
	"time"

	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jtbonhomme/asteboids/internal/fonts"
	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/jtbonhomme/asteboids/internal/vector"
	"github.com/sirupsen/logrus"
)

const (
	bulletThrottle      time.Duration = 200 * time.Millisecond
	rotationAngle       float64       = math.Pi / 36 // rotation of 5°
	starshipMaxVelocity float64       = 5.5
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
	starshipImage *ebiten.Image,
	bulletImage *ebiten.Image,
	debug bool) *Starship {
	s := Starship{
		lastBulletTime: time.Now(),
	}
	s.AgentType = physics.StarshipAgent
	s.Register = cbr
	s.Unregister = cbu
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
	defer s.Body.UpdatePosition()

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		s.Rotate(-rotationAngle)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		s.Rotate(rotationAngle)
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		s.UpdateAcceleration(1)
	} else {
		s.UpdateAcceleration(0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.SelfDestroy()
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		s.Shot()
	}

	s.UpdateVelocity()

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

	if s.Debug {
		msg := s.String()
		textDim := text.BoundString(fonts.MonoSansRegularFont, msg)
		textWidth := textDim.Max.X - textDim.Min.X
		text.Draw(screen, msg, fonts.MonoSansRegularFont, int(s.Position().X)-textWidth/2, int(s.Position().Y+s.PhysicHeight/2+5), color.Gray16{0x999f})
	}
}

// SelfDestroy removes the agent from the game
func (s *Starship) SelfDestroy() {
	defer s.Explode()
	s.Log.Debugf("SelfDestroy starship %s", s.ID())
}
