package physics

import (
	"fmt"
	"image/color"
	"math"
	"os"

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
	return fmt.Sprintf("%s: pos [%d, %d]\nsize [%d, %d] orient %0.2f rad (%0.0f °)\nvel {%0.2f %0.2f} acc {%0.2f %0.2f}",
		pb.Type(),
		int(pb.position.X),
		int(pb.position.Y),
		int(pb.PhysicWidth),
		int(pb.PhysicHeight),
		pb.Orientation,
		pb.Orientation*180/math.Pi,
		pb.Velocity().X,
		pb.Velocity().Y,
		pb.Acceleration().X,
		pb.Acceleration().Y)
}

// Intersect returns true if the physical body collide another one.
// Collision is computed based on Axis-Aligned Bounding Boxes.
// https://developer.mozilla.org/en-US/docs/Games/Techniques/2D_collision_detection
func (pb *Body) Intersect(p Physic) bool {
	ax, ay := pb.position.X, pb.position.Y
	aw, ah := pb.Dimension().W, pb.Dimension().H

	bx, by := p.Position().X, p.Position().Y
	bw, bh := p.Dimension().W, p.Dimension().H

	return (ax < bx+bw && ay < by+bh) && (ax+aw > bx && ay+ah > by)
}

// IntersectMultiple checks if multiple physical bodies are colliding with the first
func (pb *Body) IntersectMultiple(physics map[string]Physic) (string, bool) {
	for _, p := range physics {
		if pb.Intersect(p) {
			pb.Log.Warnf("%s [%d , %d] (%dx%d) intersect with %s [%d , %d] (%dx%d)",
				pb.ID(),
				int(pb.Position().X), int(pb.Position().Y),
				int(pb.Dimension().W), int(pb.Dimension().H),
				p.ID(),
				int(p.Position().X), int(p.Position().Y),
				int(p.Dimension().W), int(p.Dimension().H))
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

// LimitVelocity limits the physical body maximum velocity.
func (pb *Body) LimitVelocity(maxVelocity float64) {
	pb.maxVelocity = maxVelocity
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

// DrawBodyBoundaryBox draws a box around the body, based on its dimension.
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

func isElementOf(elt string, arr []string) bool {
	for i := 0; i < len(arr); i++ {
		if elt == arr[i] {
			return true
		}
	}
	return false
}

// LinkAgents draws a perimter around the body, based on a given radius.
func (pb *Body) LinkAgents(screen *ebiten.Image, agents []Physic, agentTypes []string) {
	for _, a := range agents {
		if isElementOf(a.Type(), agentTypes) {
			// Draw line between agents
			ebitenutil.DrawLine(
				screen,
				pb.Position().X, pb.Position().Y,
				a.Position().X, a.Position().Y,
				color.Gray16{0x2264},
			)
		}
	}
}

// Type returns physical body agent type as a string.
func (pb *Body) Type() string {
	return pb.AgentType
}

// Explode proceeds the agent explosion and termination.
func (pb *Body) Explode() {
	pb.Unregister(pb.ID(), pb.Type())
}

// Dump write out internal agent's state.
func (pb *Body) Dump(f *os.File) error {
	_, err := f.Write([]byte("\n *** " + pb.ID() + " ***\n" + pb.String() + "\n"))
	return err
}

// NewBody creates a body
func NewBody(x, y, w, h float64) *Body {
	return &Body{
		position: vector.Vector2D{
			X: x,
			Y: y,
		},
		PhysicWidth:  w,
		PhysicHeight: h,
	}
}

// FuturePosition return position the physic body will be in t time iteration,
// if velocity is constant.
func (pb *Body) FuturePosition(t float64) vector.Vector2D {
	return futurePosition(t, pb.Position(), pb.Velocity())
}

// futurePosition return position the physic body will be in t time iteration,
// if velocity is constant.
func futurePosition(t float64, pos, vel vector.Vector2D) vector.Vector2D {
	result := vel
	result.Multiply(t)
	result.Add(pos)
	return result
}
