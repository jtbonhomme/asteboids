package physics

import (
	"math"
)

func (pb *Body) normalizeOrientation() {
	if pb.Orientation > 2*math.Pi {
		pb.Orientation -= 2 * math.Pi
	}
	if pb.Orientation < 0 {
		pb.Orientation += 2 * math.Pi
	}
}

// Rotate change physical body orientation.
func (pb *Body) Rotate(rotationAngle float64) {
	pb.Orientation += rotationAngle
	pb.normalizeOrientation()
}

// Update is called every tick (1/60 [s] by default).
// Update proceeds the agent state.
func (pb *Body) Update() {
	pb.UpdateVelocity()
	pb.UpdatePosition()
}

// UpdateVelocity computes new velocity.
func (pb *Body) UpdateVelocity() {
	// update velocity from acceleration
	pb.velocity.Add(pb.acceleration)

	// limit velocity to max value
	pb.velocity.Limit(pb.maxVelocity)
}

// UpdateOrientation computes orientation from velocity.
func (pb *Body) UpdateOrientation() {
	// Update orientation from velocity
	if !pb.velocity.IsNil() {
		pb.Orientation = pb.velocity.Theta()
	}
	pb.normalizeOrientation()
}

// UpdatePosition compute new position.
func (pb *Body) UpdatePosition() {
	pb.position.Add(pb.velocity)

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
