package vector

import "math"

type Vector2D struct {
	X float64
	Y float64
}

func (v *Vector2D) Add(v2 Vector2D) {
	v.X += v2.X
	v.Y += v2.Y
}

func (v *Vector2D) Subtract(v2 Vector2D) {
	v.X -= v2.X
	v.Y -= v2.Y
}

func (v *Vector2D) Limit(max float64) {
	magSq := v.MagnitudeSquared()
	if magSq > max*max {
		v.Divide(math.Sqrt(magSq))
		v.Multiply(max)
	}
}

func (v *Vector2D) Normalize() {
	mag := math.Sqrt(v.X*v.X + v.Y*v.Y)
	v.X /= mag
	v.Y /= mag
}

func (v *Vector2D) SetMagnitude(z float64) {
	v.Normalize()
	v.X *= z
	v.Y *= z
}

func (v *Vector2D) MagnitudeSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vector2D) Divide(z float64) {
	v.X /= z
	v.Y /= z
}

func (v *Vector2D) Multiply(z float64) {
	v.X *= z
	v.Y *= z
}

func (v Vector2D) Theta() float64 {
	v2 := v
	v2.Normalize()
	theta := math.Atan(v2.Y / v2.X)
	if theta < 0 {
		theta += math.Pi * 2
	}
	if v2.X < 0 {
		theta += math.Pi
	} else if v2.Y < 0 {
		theta += math.Pi * 2
	}

	return theta
}

func (v Vector2D) Distance(v2 Vector2D) float64 {
	return math.Sqrt(math.Pow(v2.X-v.X, 2) + math.Pow(v2.Y-v.Y, 2))
}

func (v Vector2D) IsNil() bool {
	return (v.X == 0 && v.Y == 0)
}
