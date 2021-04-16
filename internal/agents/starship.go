package agents

import (
	"image/color"
	"time"

	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jtbonhomme/asteboids/internal/game"
	"github.com/sirupsen/logrus"
)

const (
	bulletThrottle time.Duration = 100 * time.Millisecond
)

// Starship is a PhysicalBody agent.
// It represents a playable star ship.
type Starship struct {
	game.PhysicBody
	bullets        map[string]*Bullet
	lastBulletTime time.Time
}

// NewStarship creates a new Starship (PhysicalBody agent)
func NewStarship(log *logrus.Logger, x, y, screenWidth, screenHeight int, cb game.AgentUnregister, debug bool) *Starship {
	s := Starship{
		bullets:        make(map[string]*Bullet),
		lastBulletTime: time.Now(),
	}
	s.Type = "starship"
	s.Unregister = cb
	s.Init()

	s.Orientation = math.Pi / 2
	s.Velocity = game.Vector{
		X: 0,
		Y: 0,
	}
	s.Acceleration = game.Vector{
		X: 0,
		Y: 0,
	}
	s.PhysicWidth = 50
	s.PhysicHeight = 50
	s.ScreenWidth = screenWidth
	s.ScreenHeight = screenHeight
	s.X = x
	s.Y = y
	s.Log = log

	err := s.LoadImage("./assets/ship.png")
	if err != nil {
		log.Errorf("error when loading image from file: %s", err.Error())
	}
	s.Debug = debug
	return &s
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (s *Starship) Update() {
	defer s.PhysicBody.Update()

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		s.Rotate(-1)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		s.Rotate(1)
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
		s.RegisterBullet()
	}

	// Update bullets
	for _, b := range s.bullets {
		b.Update()
	}
}

// Register adds a new bullet to the game.
func (s *Starship) RegisterBullet() {
	// throtlle call to avoid continuous shooting
	if time.Since(s.lastBulletTime) < bulletThrottle {
		return
	}
	s.lastBulletTime = time.Now()
	bullet := NewBullet(s.Log, s.X, s.Y, s.Orientation, s.ScreenWidth, s.ScreenHeight, s.UnregisterBullet)
	s.bullets[bullet.ID()] = bullet
}

// Unregister deletes an bullet from the game.
func (s *Starship) UnregisterBullet(id string) {
	delete(s.bullets, id)
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (s *Starship) Draw(screen *ebiten.Image) {
	defer s.PhysicBody.Draw(screen)
	// Update bullets
	for _, b := range s.bullets {
		b.Draw(screen)
	}
	if s.Debug {
		text.Draw(screen, s.String(), game.MplusNormalFont, s.X, s.Y-s.PhysicHeight/2+5, color.White)
	}
}

// SelfDestroy removes the agent from the game
func (s *Starship) SelfDestroy() {
	defer s.Unregister(s.ID())
	s.Log.Infof("SelfDestroy starship %s", s.ID())
}

// todo unregister bullet as cb
