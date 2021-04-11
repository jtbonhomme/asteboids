package player

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/game"
	"github.com/sirupsen/logrus"
)

type Starship struct {
	log           *logrus.Logger
	position      game.Position
	direction     float64
	size          float64
	rotationAngle float64

	speed        float64
	acceleration float64
	vertices     []ebiten.Vertex

	emptyImage *ebiten.Image
}

func NewStarship(log *logrus.Logger, x, y int) *Starship {
	s := Starship{
		direction:     90.0,
		speed:         0.0,
		acceleration:  0.0,
		size:          20,
		rotationAngle: 5,
		position: game.Position{
			X: x,
			Y: y,
		},
		log:      log,
		vertices: []ebiten.Vertex{},
	}
	s.updateVertices()
	s.emptyImage = ebiten.NewImage(3, 3)
	s.emptyImage.Fill(color.White)

	return &s
}

// Draw draws the agent.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (s *Starship) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawTrianglesOptions{}
	op.Address = ebiten.AddressUnsafe
	screen.DrawTriangles(s.vertices, []uint16{0, 1, 2}, s.emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)
}

func (s *Starship) updateVertices() {
	vs := []ebiten.Vertex{}
	for i := 0; i < 3; i++ {
		vs = append(vs, ebiten.Vertex{
			DstX:   0,
			DstY:   0,
			SrcX:   0,
			SrcY:   0,
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		})
	}
	centerX := s.position.X
	centerY := s.position.Y
	// spaceship head
	vs[0].DstX = float32(centerX + int(s.size*math.Cos(s.direction*2*math.Pi/360)))
	vs[0].DstY = float32(centerY + int(s.size*math.Sin(s.direction*2*math.Pi/360)))
	// spaceship base
	vs[1].DstX = float32(centerX + int(s.size*math.Cos(s.direction*2*math.Pi/360+2*math.Pi/3)))
	vs[1].DstY = float32(centerY + int(s.size*math.Sin(s.direction*2*math.Pi/360+2*math.Pi/3)))
	vs[2].DstX = float32(centerX + int(s.size*math.Cos(s.direction*2*math.Pi/360+4*math.Pi/3)))
	vs[2].DstY = float32(centerY + int(s.size*math.Sin(s.direction*2*math.Pi/360+4*math.Pi/3)))

	s.vertices = vs
}

func (s *Starship) rotate(i float64) {
	s.direction += i * s.rotationAngle
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (s *Starship) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		s.rotate(-1)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		s.rotate(1)
	}
	s.updateVertices()
}

// Position returns the current agent position
func (s *Starship) Position() game.Position {
	return s.position
}

// Speed returns the current agent speed
func (s *Starship) Speed() float64 {
	return s.speed
}

// Direction returns the current agent direction
func (s *Starship) Direction() float64 {
	return s.direction
}
