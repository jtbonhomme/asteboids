package physics

import (
	"fmt"
	"image/color"
	"math"

	// anonymous import for png decoder
	_ "image/png"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jtbonhomme/asteboids/internal/vector"
	"github.com/sirupsen/logrus"
)

type PhysicBody struct {
	position    vector.Vector2D
	AgentType   string
	id          uuid.UUID
	Log         *logrus.Logger
	Orientation float64 // theta (radian)

	PhysicWidth  float64
	PhysicHeight float64
	ScreenWidth  float64
	ScreenHeight float64

	Velocity     vector.Vector2D
	Acceleration vector.Vector2D

	Register   AgentRegister
	Unregister AgentUnregister
	Vision     AgentVision
	Image      *ebiten.Image

	Debug bool
}

// Init initializes the physic body
func (pb *PhysicBody) Init(x, y float64) {
	pb.id = uuid.New()
	pb.position.X = x
	pb.position.Y = y
}

// Draw draws the agent.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (pb *PhysicBody) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	defer screen.DrawImage(pb.Image, op)

	op.GeoM.Translate(-pb.PhysicWidth/2, -pb.PhysicHeight/2)
	op.GeoM.Rotate(pb.Orientation)
	op.GeoM.Translate(pb.position.X, pb.position.Y)
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

func (pb *PhysicBody) UpdateVelocity() {
	pb.Velocity.X += pb.Acceleration.X - frictionFactor*pb.Velocity.X
	pb.Velocity.Y += pb.Acceleration.Y - frictionFactor*pb.Velocity.Y

	velocityValue := pb.Velocity.X*pb.Velocity.X + pb.Velocity.Y*pb.Velocity.Y
	if velocityValue > maxVelocity*maxVelocity {
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
	// update position
	pb.position.X += velocityFactor * pb.Velocity.X
	pb.position.Y += velocityFactor * pb.Velocity.Y

	if pb.position.X > pb.ScreenWidth {
		pb.position.X = 0
	} else if pb.position.X < 0 {
		pb.position.X = pb.ScreenWidth
	}
	if pb.position.Y > pb.ScreenHeight {
		pb.position.Y = 0
	} else if pb.position.Y < 0 {
		pb.position.Y = pb.ScreenHeight
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
		int(pb.position.X),
		int(pb.position.Y),
		int(pb.PhysicWidth),
		int(pb.PhysicHeight),
		pb.Orientation,
		pb.Orientation*180/math.Pi,
		pb.Velocity.X,
		pb.Velocity.Y)
}

// Intersect returns true if the physical body collide another one.
// Collision is computed based on Axis-Aligned Bounding Boxes.
// https://developer.mozilla.org/en-US/docs/Games/Techniques/2D_collision_detection
func (pb *PhysicBody) Intersect(p Physic) bool {
	ax, ay := pb.position.X, pb.position.Y
	aw, ah := pb.Dimension().W*collisionPrecision, pb.Dimension().H*collisionPrecision

	bx, by := p.Position().X, p.Position().Y
	bw, bh := p.Dimension().W*collisionPrecision, p.Dimension().H*collisionPrecision

	return (ax < bx+bw && ay < by+bh) && (ax+aw > bx && ay+ah > by)
}

// IntersectMultiple checks if multiple physical bodies are colliding with the first
func (pb *PhysicBody) IntersectMultiple(physics map[string]Physic) (string, bool) {
	for _, p := range physics {
		if pb.Intersect(p) {
			return p.ID(), true
		}
	}
	return "", false
}

// Dimension returns physical body dimension.
func (pb *PhysicBody) Dimension() Size {
	return Size{
		H: pb.PhysicHeight,
		W: pb.PhysicWidth,
	}
}

// Position returns physical body position.
func (pb *PhysicBody) Position() vector.Vector2D {
	return pb.position
}

func (pb *PhysicBody) DrawBodyBoundaryBox(screen *ebiten.Image) {
	// Top boundary
	ebitenutil.DrawLine(
		screen,
		pb.position.X-pb.PhysicWidth/2,
		pb.position.Y-pb.PhysicHeight/2,
		pb.position.X+pb.PhysicWidth/2,
		pb.position.Y-pb.PhysicHeight/2,
		color.Gray16{0x6666},
	)
	// Right boundary
	ebitenutil.DrawLine(
		screen,
		pb.position.X+pb.PhysicWidth/2,
		pb.position.Y-pb.PhysicHeight/2,
		pb.position.X+pb.PhysicWidth/2,
		pb.position.Y+pb.PhysicHeight/2,
		color.Gray16{0x6666},
	)
	// Bottom boundary
	ebitenutil.DrawLine(
		screen,
		pb.position.X-pb.PhysicWidth/2,
		pb.position.Y+pb.PhysicHeight/2,
		pb.position.X+pb.PhysicWidth/2,
		pb.position.Y+pb.PhysicHeight/2,
		color.Gray16{0x6666},
	)
	// Left boundary
	ebitenutil.DrawLine(
		screen,
		pb.position.X-pb.PhysicWidth/2,
		pb.position.Y-pb.PhysicHeight/2,
		pb.position.X-pb.PhysicWidth/2,
		pb.position.Y+pb.PhysicHeight/2,
		color.Gray16{0x6666},
	)
}

// Type returns physical body agent type as a string.
func (pb *PhysicBody) Type() string {
	return pb.AgentType
}

// Explode proceeds the agent explosion and termination.
func (pb *PhysicBody) Explode() {
	pb.Unregister(pb.ID(), pb.Type())
}
