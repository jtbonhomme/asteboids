package physics

import (
	"fmt"
	"image/color"
	"math"

	// anonymous import for png decoder
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jtbonhomme/asteboids/internal/vector"
)

// ID displays physic body unique ID.
func (pb *Body) ID() string {
	return pb.id.String()
}

// String displays physic body information as a string.
func (pb *Body) String() string {
	return fmt.Sprintf("%s: [%d, %d] [%d, %d]\n%0.2f rad (%0.0f °) {%0.2f %0.2f}",
		pb.Type(),
		int(pb.position.X),
		int(pb.position.Y),
		int(pb.PhysicWidth),
		int(pb.PhysicHeight),
		pb.Orientation,
		pb.Orientation*180/math.Pi,
		pb.Velocity().X,
		pb.Velocity().Y)
}

// Intersect returns true if the physical body collide another one.
// Collision is computed based on Axis-Aligned Bounding Boxes.
// https://developer.mozilla.org/en-US/docs/Games/Techniques/2D_collision_detection
func (pb *Body) Intersect(p Physic) bool {
	ax, ay := pb.position.X, pb.position.Y
	aw, ah := pb.Dimension().W*collisionPrecision, pb.Dimension().H*collisionPrecision

	bx, by := p.Position().X, p.Position().Y
	bw, bh := p.Dimension().W*collisionPrecision, p.Dimension().H*collisionPrecision

	return (ax < bx+bw && ay < by+bh) && (ax+aw > bx && ay+ah > by)
}

// IntersectMultiple checks if multiple physical bodies are colliding with the first
func (pb *Body) IntersectMultiple(physics map[string]Physic) (string, bool) {
	for _, p := range physics {
		if pb.Intersect(p) {
			return p.ID(), true
		}
	}
	return "", false
}

// Dimension returns physical body dimension.
func (pb *Body) Dimension() Size {
	return Size{
		H: pb.PhysicHeight,
		W: pb.PhysicWidth,
	}
}

// Position returns physical body position.
func (pb *Body) Position() vector.Vector2D {
	return pb.position
}

// Velocity returns physical body velocity.
func (pb *Body) Velocity() vector.Vector2D {
	return pb.velocity
}

// Move set physical body positiion.
func (pb *Body) Move(position vector.Vector2D) {
	pb.position = position
}

// Accelerate set physical body acceleration.
func (pb *Body) Accelerate(acceleration vector.Vector2D) {
	pb.acceleration = acceleration
}

// Acceleration returns physical body acceleration.
func (pb *Body) Acceleration() vector.Vector2D {
	return pb.acceleration
}

func (pb *Body) DrawBodyBoundaryBox(screen *ebiten.Image) {
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
func (pb *Body) Type() string {
	return pb.AgentType
}

// Explode proceeds the agent explosion and termination.
func (pb *Body) Explode() {
	pb.Unregister(pb.ID(), pb.Type())
}
