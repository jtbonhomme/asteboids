package physics

import (
	"math"
)

// Rotate change physical body orientation.
func (pb *Body) Rotate(rotationAngle float64) {
	pb.Orientation += rotationAngle
	if pb.Orientation > 2*math.Pi {
		pb.Orientation -= 2 * math.Pi
	}
	if pb.Orientation < 0 {
		pb.Orientation += 2 * math.Pi
	}
}

// Update is called every tick (1/60 [s] by default).
// Update proceeds the agent state.
func (pb *Body) Update() {
	pb.UpdateAcceleration(1)
	pb.UpdateVelocity()
	pb.UpdatePosition()
}

// UpdatePosition compute new acceleration.
func (pb *Body) UpdateAcceleration(i float64) {
	pb.acceleration.X = accelerationFactor * i * math.Cos(pb.Orientation)
	pb.acceleration.Y = accelerationFactor * i * math.Sin(pb.Orientation)
}

// UpdatePosition compute new velocity.
func (pb *Body) UpdateVelocity() {
	pb.velocity.X += pb.acceleration.X - frictionFactor*pb.velocity.X
	pb.velocity.Y += pb.acceleration.Y - frictionFactor*pb.velocity.Y

	velocityValue := pb.velocity.X*pb.velocity.X + pb.velocity.Y*pb.velocity.Y
	if velocityValue > pb.maxVelocity*pb.maxVelocity {
		pb.velocity.X = pb.maxVelocity * math.Cos(pb.Orientation)
		pb.velocity.Y = pb.maxVelocity * math.Sin(pb.Orientation)
	}
	if velocityValue < 0 {
		pb.velocity.X = 0
		pb.velocity.Y = 0
	}
}

// UpdatePosition compute new position.
func (pb *Body) UpdatePosition() {
	pb.position.X += velocityFactor * pb.velocity.X
	pb.position.Y += velocityFactor * pb.velocity.Y

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
