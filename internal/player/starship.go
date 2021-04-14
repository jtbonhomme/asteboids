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

const (
	accelerationFactor float64 = 0.3
	velocityFactor     float64 = 1.8
	maxVelocity        float64 = 5.5
	rotationAngle      float64 = math.Pi / 36 // rotation of 5°
	frictionFactor     float64 = 0.03
)

type Starship struct {
	log         *logrus.Logger
	position    game.Position
	orientation float64 // theta (radian)
	size        float64

	starshipWidth  int
	starshipHeight int
	screenWidth    int
	screenHeight   int

	velocity     game.Vector
	acceleration game.Vector

	//vertices        []ebiten.Vertex

	image *ebiten.Image
}

func NewStarship(log *logrus.Logger, x, y, screenWidth, screenHeight int) *Starship {
	s := Starship{
		orientation: math.Pi / 2,
		velocity: game.Vector{
			X: 0,
			Y: 0,
		},
		acceleration: game.Vector{
			X: 0,
			Y: 0,
		},
		size:           20,
		starshipWidth:  50,
		starshipHeight: 50,
		screenWidth:    screenWidth,
		screenHeight:   screenHeight,
		position: game.Position{
			X: x,
			Y: y,
		},
		log: log,
		//vertices: []ebiten.Vertex{},
	}
	//s.updateVertices()
	s.image = ebiten.NewImage(s.starshipWidth, s.starshipHeight)

	f, err := os.Open("./assets/ship.png")
	if err != nil {
		s.log.Errorf("error when opening file: %s", err.Error())
		s.image.Fill(color.White)
		return &s
	}

	defer f.Close()
	rawImage, _, err := image.Decode(f)
	if err != nil {
		s.log.Errorf("error when decoding image from file: %s", err.Error())
		s.image.Fill(color.White)
		return &s
	}

	newImage := ebiten.NewImageFromImage(rawImage)
	if newImage != nil {
		s.image.DrawImage(newImage, nil)
	} else {
		s.log.Errorf("error when creating image from raw: %s", err.Error())
		s.image.Fill(color.White)
	}
	return &s
}

// Draw draws the agent.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (s *Starship) Draw(screen *ebiten.Image) {
	/*op := &ebiten.DrawTrianglesOptions{}
	op.Address = ebiten.AddressUnsafe
	screen.DrawTriangles(s.vertices, []uint16{0, 1, 2}, s.image.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)*/

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(-s.starshipWidth/2), float64(-s.starshipHeight/2))
	op.GeoM.Rotate(s.orientation)
	op.GeoM.Translate(float64(s.position.X), float64(s.position.Y))

	screen.DrawImage(s.image, op)
}

/*
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
	vs[0].DstX = float32(centerX + int(s.size*math.Cos(s.orientation)))
	vs[0].DstY = float32(centerY + int(s.size*math.Sin(s.orientation)))
	// spaceship base
	vs[1].DstX = float32(centerX + int(s.size*math.Cos(s.orientation+2*math.Pi/3)))
	vs[1].DstY = float32(centerY + int(s.size*math.Sin(s.orientation+2*math.Pi/3)))
	vs[2].DstX = float32(centerX + int(s.size*math.Cos(s.orientation+4*math.Pi/3)))
	vs[2].DstY = float32(centerY + int(s.size*math.Sin(s.orientation+4*math.Pi/3)))

	s.vertices = vs
}*/

func (s *Starship) rotate(i float64) {
	s.orientation += i * rotationAngle
	if s.orientation > 2*math.Pi {
		s.orientation -= 2 * math.Pi
	}
	if s.orientation < 0 {
		s.orientation += 2 * math.Pi
	}
}

func (s *Starship) updateAcceleration(i float64) {
	s.acceleration.X = accelerationFactor * i * math.Cos(s.orientation)
	s.acceleration.Y = accelerationFactor * i * math.Sin(s.orientation)
}

func (s *Starship) updateVelocity() {
	s.velocity.X += s.acceleration.X - frictionFactor*s.velocity.X
	s.velocity.Y += s.acceleration.Y - frictionFactor*s.velocity.Y

	velocityValue := math.Sqrt(s.velocity.X*s.velocity.X + s.velocity.Y*s.velocity.Y)
	if velocityValue > maxVelocity {
		s.velocity.X = maxVelocity * math.Cos(s.orientation)
		s.velocity.Y = maxVelocity * math.Sin(s.orientation)
	}
	if velocityValue < 0 {
		s.velocity.X = 0
		s.velocity.Y = 0
	}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (s *Starship) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		s.rotate(-1)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		s.rotate(1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		s.updateAcceleration(1)
		/*	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
			s.updateAcceleration(-1)*/
	} else {
		s.updateAcceleration(0)
	}

	s.updateVelocity()

	// update position
	s.position.X += int(velocityFactor * s.velocity.X)
	s.position.Y += int(velocityFactor * s.velocity.Y)

	if s.position.X > s.screenWidth {
		s.position.X = 0
	} else if s.position.X < 0 {
		s.position.X = s.screenWidth
	}
	if s.position.Y > s.screenHeight {
		s.position.Y = 0
	} else if s.position.Y < 0 {
		s.position.Y = s.screenHeight
	}
	//s.updateVertices()
}

// Position returns the current agent position
func (s *Starship) Position() game.Position {
	return s.position
}

// Velocity returns the current agent speed vector
func (s *Starship) Velocity() (v game.Vector, n float64) {
	return s.velocity, math.Sqrt(s.velocity.X*s.velocity.X + s.velocity.Y*s.velocity.Y)
}

// Acceleration returns the current agent acceleration vector
func (s *Starship) Acceleration() (v game.Vector, n float64) {
	return s.acceleration, math.Sqrt(s.acceleration.X*s.acceleration.X + s.acceleration.Y*s.acceleration.Y)
}

// orientation returns the current agent orientation in degre
func (s *Starship) Orientation() float64 {
	return s.orientation * 360 / (2 * math.Pi)
}
