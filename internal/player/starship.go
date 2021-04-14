package player

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

type Starship struct {
	game.AgentBody
}

func NewStarship(log *logrus.Logger, x, y, screenWidth, screenHeight int) *Starship {
	s := Starship{}
	s.Orientation = math.Pi / 2
	s.Velocity = game.Vector{
		X: 0,
		Y: 0,
	}
	s.Acceleration = game.Vector{
		X: 0,
		Y: 0,
	}
	s.Size = 20
	s.AgentWidth = 50
	s.AgentHeight = 50
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
	s.AgentBody.Update()
}
