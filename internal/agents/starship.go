package agents

import (
	"image"
	"image/color"

	// anonymous import for png decoder
	_ "image/png"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/game"
	"github.com/sirupsen/logrus"
)

// Starship is a PhysicalBody agent.
// It represents a playable star ship.
type Starship struct {
	game.PhysicBody
	bullets map[string]*Bullet
}

// NewStarship creates a new Starship (PhysicalBody agent)
func NewStarship(log *logrus.Logger, x, y, screenWidth, screenHeight int, cb game.AgentUnregister) *Starship {
	s := Starship{
		bullets: make(map[string]*Bullet),
	}
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

	s.Image = ebiten.NewImage(s.ScreenWidth, s.ScreenHeight)

	f, err := os.Open("./assets/ship.png")
	if err != nil {
		log.Errorf("error when opening file: %s", err.Error())
		s.Image.Fill(color.White)
		return &s
	}

	defer f.Close()
	rawImage, _, err := image.Decode(f)
	if err != nil {
		log.Errorf("error when decoding image from file: %s", err.Error())
		s.Image.Fill(color.White)
		return &s
	}

	newImage := ebiten.NewImageFromImage(rawImage)
	if newImage != nil {
		s.Image.DrawImage(newImage, nil)
	} else {
		log.Errorf("error when creating image from raw: %s", err.Error())
		s.Image.Fill(color.White)
	}
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
		s.RegisterBullet(NewBullet(s.Log, s.X, s.Y, s.Orientation, s.ScreenWidth, s.ScreenHeight, s.UnregisterBullet))
	}

	// Update bullets
	for _, b := range s.bullets {
		b.Update()
	}
}

// Register adds a new bullet to the game.
func (s *Starship) RegisterBullet(bullet *Bullet) {
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
}

// SelfDestroy removes the agent from the game
func (s *Starship) SelfDestroy() {
	defer s.Unregister(s.ID())
	s.Log.Infof("SelfDestroy starship %s", s.ID())
}

// todo unregister bullet as cb
