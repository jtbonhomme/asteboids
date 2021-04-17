package physics

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	// anonymous import for png decoder
	_ "image/png"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sirupsen/logrus"
)

const (
	accelerationFactor float64 = 0.3
	velocityFactor     float64 = 1.8
	maxVelocity        float64 = 5.5
	frictionFactor     float64 = 0.03
	collisionPrecision float64 = 1.0
)

const (
	StarshipAgent string = "starship"
	AsteroidAgent string = "asteroid"
	BulletAgent   string = "bullet"
)

// Size represents coordonnates (X, Y) of a physical body.
type Position struct {
	X int
	Y int
}

// Size represents height and width of a physical body.
type Size struct {
	H int
	W int
}

// Vector represents a vector composantes.
type Vector struct {
	X float64
	Y float64
}

// Block is a dimension and position helper structure.
type Block struct {
	Position
	Size
}

type Physic interface {
	// Draw draws the agent on screen.
	Draw(*ebiten.Image)
	// Update proceeds the agent state.
	Update()
	// Init initializes the physic body.
	Init()
	// ID displays physic body unique ID.
	ID() string
	// LoadImage loads a picture in an ebiten image.
	LoadImage(string) error
	// String displays physic body information as a string.
	String() string
	// Intersect returns true if the physical body collide another one.
	// Collision is computed based on Axis-Aligned Bounding Boxes.
	// https://developer.mozilla.org/en-US/docs/Games/Techniques/2D_collision_detection
	Intersect(Physic) bool
	// IntersectMultiple checks if multiple physical bodies are colliding with the first
	IntersectMultiple([]Physic) bool
	// Dimensions returns physical body dimensions.
	Dimension() Block
	// Type returns physical body agent type as a string.
	Type() string
	// Explode proceeds the agent explosion and termination.
	Explode()
}

// AgentUnregister is a function to unregister an agent
type AgentUnregister func(string, string)

type PhysicBody struct {
	Position
	AgentType   string
	id          uuid.UUID
	Log         *logrus.Logger
	Orientation float64 // theta (radian)
	Size        float64

	PhysicWidth  int
	PhysicHeight int
	ScreenWidth  int
	ScreenHeight int

	Velocity     Vector
	Acceleration Vector

	Unregister AgentUnregister
	Image      *ebiten.Image

	Debug bool
}

// Init initializes the physic body
func (pb *PhysicBody) Init() {
	pb.id = uuid.New()
}

// Draw draws the agent.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (pb *PhysicBody) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	defer screen.DrawImage(pb.Image, op)

	op.GeoM.Translate(float64(-pb.PhysicWidth/2), float64(-pb.PhysicHeight/2))
	op.GeoM.Rotate(pb.Orientation)
	op.GeoM.Translate(float64(pb.X), float64(pb.Y))
	if pb.Debug {
		pb.DrawBodyBoundaryBox(screen)
	}
}

func (pb *PhysicBody) Rotate(rotationAngle float64) {
	pb.Orientation += rotationAngle
	if pb.Orientation > 2*math.Pi {
		pb.Orientation -= 2 * math.Pi
	}
	if pb.Orientation < 0 {
		pb.Orientation += 2 * math.Pi
	}
}

func (pb *PhysicBody) UpdateAcceleration(i float64) {
	pb.Acceleration.X = accelerationFactor * i * math.Cos(pb.Orientation)
	pb.Acceleration.Y = accelerationFactor * i * math.Sin(pb.Orientation)
}

func (pb *PhysicBody) updateVelocity() {
	pb.Velocity.X += pb.Acceleration.X - frictionFactor*pb.Velocity.X
	pb.Velocity.Y += pb.Acceleration.Y - frictionFactor*pb.Velocity.Y

	velocityValue := math.Sqrt(pb.Velocity.X*pb.Velocity.X + pb.Velocity.Y*pb.Velocity.Y)
	if velocityValue > maxVelocity {
		pb.Velocity.X = maxVelocity * math.Cos(pb.Orientation)
		pb.Velocity.Y = maxVelocity * math.Sin(pb.Orientation)
	}
	if velocityValue < 0 {
		pb.Velocity.X = 0
		pb.Velocity.Y = 0
	}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (pb *PhysicBody) Update() {
	pb.updateVelocity()

	// update position
	pb.X += int(velocityFactor * pb.Velocity.X)
	pb.Y += int(velocityFactor * pb.Velocity.Y)

	if pb.X > pb.ScreenWidth {
		pb.X = 0
	} else if pb.X < 0 {
		pb.X = pb.ScreenWidth
	}
	if pb.Y > pb.ScreenHeight {
		pb.Y = 0
	} else if pb.Y < 0 {
		pb.Y = pb.ScreenHeight
	}
}

// ID displays physic body unique ID.
func (pb *PhysicBody) ID() string {
	return pb.id.String()
}

// String displays physic body information as a string.
func (pb *PhysicBody) String() string {
	return fmt.Sprintf("%s: [%d, %d] [%d, %d]\n%0.2f rad (%0.0f °) {%0.2f %0.2f}",
		pb.Type(),
		pb.X,
		pb.Y,
		pb.PhysicWidth,
		pb.PhysicHeight,
		pb.Orientation,
		pb.Orientation*180/math.Pi,
		pb.Velocity.X,
		pb.Velocity.Y)
}

// LoadImage loads a picture in an ebiten image.
func (pb *PhysicBody) LoadImage(file string) error {
	pb.Image = ebiten.NewImage(pb.ScreenWidth, pb.ScreenHeight)

	f, err := os.Open(file)
	if err != nil {
		pb.Image.Fill(color.White)
		return errors.New("error when opening file " + err.Error())
	}

	defer f.Close()
	rawImage, _, err := image.Decode(f)
	if err != nil {
		pb.Image.Fill(color.White)
		return errors.New("error when decoding image from file " + err.Error())
	}

	newImage := ebiten.NewImageFromImage(rawImage)
	if newImage == nil {
		pb.Image.Fill(color.White)
		return errors.New("error when creating image from raw " + err.Error())
	}
	pb.Image.DrawImage(newImage, nil)
	return nil
}

// Intersect returns true if the physical body collide another one.
// Collision is computed based on Axis-Aligned Bounding Boxes.
// https://developer.mozilla.org/en-US/docs/Games/Techniques/2D_collision_detection
func (pb *PhysicBody) Intersect(p Physic) bool {
	ax, ay := pb.Dimension().X, pb.Dimension().Y
	aw, ah := int(float64(pb.Dimension().W)*collisionPrecision), int(float64(pb.Dimension().H)*collisionPrecision)

	bx, by := p.Dimension().X, p.Dimension().Y
	bw, bh := int(float64(p.Dimension().W)*collisionPrecision), int(float64(p.Dimension().H)*collisionPrecision)

	intersect := (ax < bx+bw && ay < by+bh) && (ax+aw > bx && ay+ah > by)
	if intersect {
		pb.Log.Debugf("%s agent %s [%d, %d] (%d x %d) has collide with %s agent %s [%d, %d] (%d x %d)",
			pb.Type(), pb.ID(),
			ax, ay, aw, ah,
			p.Type(), p.ID(),
			bx, by, bw, bh,
		)
	}
	return intersect
}

// IntersectMultiple checks if multiple physical bodies are colliding with the first
func (pb *PhysicBody) IntersectMultiple(physics []Physic) bool {
	for _, p := range physics {
		if pb.Intersect(p) {
			return true
		}
	}
	return false
}

// Dimensions returns physical body dimensions.
func (pb *PhysicBody) Dimension() Block {
	return Block{
		Position{
			X: pb.X,
			Y: pb.Y,
		},
		Size{
			H: pb.PhysicHeight,
			W: pb.PhysicWidth,
		},
	}
}

func (pb *PhysicBody) DrawBodyBoundaryBox(screen *ebiten.Image) {
	// Top boundary
	ebitenutil.DrawLine(
		screen,
		float64(pb.X-pb.PhysicWidth/2),
		float64(pb.Y-pb.PhysicHeight/2),
		float64(pb.X+pb.PhysicWidth/2),
		float64(pb.Y-pb.PhysicHeight/2),
		color.Gray16{0x6666},
	)
	// Right boundary
	ebitenutil.DrawLine(
		screen,
		float64(pb.X+pb.PhysicWidth/2),
		float64(pb.Y-pb.PhysicHeight/2),
		float64(pb.X+pb.PhysicWidth/2),
		float64(pb.Y+pb.PhysicHeight/2),
		color.Gray16{0x6666},
	)
	// Bottom boundary
	ebitenutil.DrawLine(
		screen,
		float64(pb.X-pb.PhysicWidth/2),
		float64(pb.Y+pb.PhysicHeight/2),
		float64(pb.X+pb.PhysicWidth/2),
		float64(pb.Y+pb.PhysicHeight/2),
		color.Gray16{0x6666},
	)
	// Left boundary
	ebitenutil.DrawLine(
		screen,
		float64(pb.X-pb.PhysicWidth/2),
		float64(pb.Y-pb.PhysicHeight/2),
		float64(pb.X-pb.PhysicWidth/2),
		float64(pb.Y+pb.PhysicHeight/2),
		color.Gray16{0x6666},
	)
}

// Type returns physical body agent type as a string.
func (pb *PhysicBody) Type() string {
	return pb.AgentType
}

// Explode proceeds the agent explosion and termination.
func (pb *PhysicBody) Explode() {
	defer pb.Unregister(pb.ID(), pb.Type())
	pb.Log.Infof("%s agent %s exploded !", pb.Type(), pb.ID())
}
