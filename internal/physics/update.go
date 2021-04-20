package physics

import (
	"math"
)

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
	pb.UpdateVelocity()
}

// UpdatePosition compute new acceleration.
func (pb *Body) UpdateAcceleration(i float64) {
	pb.Acceleration.X = accelerationFactor * i * math.Cos(pb.Orientation)
	pb.Acceleration.Y = accelerationFactor * i * math.Sin(pb.Orientation)
}

// UpdatePosition compute new velocity.
func (pb *Body) UpdateVelocity() {
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

// UpdatePosition compute new position.
func (pb *Body) UpdatePosition() {
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
