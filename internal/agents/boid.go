package agents

import (
	"image/color"
	"math"
	"math/rand"

	// anonymous import for png decoder
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/jtbonhomme/asteboids/internal/vector"
	"github.com/sirupsen/logrus"
)

const (
	boidMaxForce     float64 = 0.3
	boidMaxVelocity  float64 = 3.0
	separationFactor float64 = 1.9
	cohesionFactor   float64 = 1.5
	alignmentFactor  float64 = 1.3
)

// Boid is a PhysicalBody agent.
// It represents a single autonomous agent.
type Boid struct {
	physics.Body
}

// NewBoid creates a new Boid (PhysicalBody agent)
func NewBoid(
	log *logrus.Logger,
	x, y,
	screenWidth, screenHeight float64,
	boidImage *ebiten.Image,
	vision physics.AgentVision,
	debug bool) *Boid {
	b := Boid{}
	b.AgentType = physics.BoidAgent

	b.Orientation = math.Pi / 32 * float64(rand.Intn(64))

	b.Init(vector.Vector2D{
		X: boidMaxVelocity * math.Cos(b.Orientation),
		Y: boidMaxVelocity * math.Sin(b.Orientation),
	})
	b.LimitVelocity(boidMaxVelocity)
	b.Log = log

	b.Move(vector.Vector2D{
		X: x,
		Y: y,
	})
	b.PhysicWidth = 10
	b.PhysicHeight = 10
	b.ScreenWidth = screenWidth
	b.ScreenHeight = screenHeight

	emptyImage := ebiten.NewImage(10, 10)
	emptyImage.Fill(color.RGBA{100, 100, 200, 255})
	op := &ebiten.DrawTrianglesOptions{}
	op.Address = ebiten.AddressRepeat
	op.CompositeMode = ebiten.CompositeModeSourceOut
	b.Image = ebiten.NewImage(10, 10)
	b.Image.DrawTriangles(
		[]ebiten.Vertex{
			{DstX: float32(0), DstY: float32(0), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
			{DstX: float32(10), DstY: float32(5), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
			{DstX: float32(5), DstY: float32(5), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
			{DstX: float32(0), DstY: float32(10), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		},
		[]uint16{0, 1, 2, 3, 1, 2},
		emptyImage,
		op,
	)
	b.Vision = vision
	b.VisionRadius = 100
	b.Debug = debug
	return &b
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (b *Boid) Update() {

	acceleration := vector.Vector2D{}
	nearestAgent := b.Vision(b.Position(), b.VisionRadius)

	cohesion := b.cohesion(nearestAgent)
	cohesion.Multiply(cohesionFactor)
	acceleration.Add(cohesion)

	separation := b.separate(nearestAgent)
	separation.Multiply(separationFactor)
	acceleration.Add(separation)

	alignment := b.align(nearestAgent)
	alignment.Multiply(alignmentFactor)
	acceleration.Add(alignment)

	b.Accelerate(acceleration)
	b.UpdateVelocity()
	b.UpdateOrientation()
	b.UpdatePosition()
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (b *Boid) Draw(screen *ebiten.Image) {
	defer b.Body.Draw(screen)
	nearestAgent := b.Vision(b.Position(), b.VisionRadius)
	b.LinkAgents(screen, nearestAgent, []string{physics.BoidAgent})
}

func (b *Boid) seek(target vector.Vector2D) vector.Vector2D {
	desired := b.Position()
	desired.Subtract(target)
	desired.Normalize()
	desired.Multiply(boidMaxVelocity)
	steer := b.Velocity()
	steer.Subtract(desired)
	steer.Limit(boidMaxForce)
	return steer
}

// cohesion returns the force imposed by flocking cohesion rule.
func (b *Boid) cohesion(agents []physics.Physic) vector.Vector2D {
	result := vector.Vector2D{
		X: 0,
		Y: 0,
	}
	var nBoids float64 = 0.0
	for _, agent := range agents {
		if agent.Type() == physics.BoidAgent && agent.ID() != b.ID() {
			nBoids++
			result.Add(agent.Position())
		}
	}
	if nBoids > 0 {
		result.Divide(nBoids)
		result = b.seek(result)
	}
	return result
}

// separate returns the force imposed by flocking separation rule.
func (b *Boid) separate(agents []physics.Physic) vector.Vector2D {
	result := vector.Vector2D{
		X: 0,
		Y: 0,
	}
	var nBoids float64 = 0.0
	for _, agent := range agents {
		if agent.ID() == b.ID() {
			continue
		}
		nBoids++
		d := b.Position().Distance(agent.Position())
		diff := b.Position()
		diff.Subtract(agent.Position())
		diff.Normalize()
		diff.Divide(d)
		result.Add(diff)
	}
	if nBoids > 0 {
		result.Divide(nBoids)
		result.Normalize()
		result.Multiply(boidMaxVelocity)
		result.Subtract(b.Velocity())
		result.Limit(boidMaxForce)
	}
	return result
}

// align returns the force imposed by flocking alignment rule.
func (b *Boid) align(agents []physics.Physic) vector.Vector2D {
	result := vector.Vector2D{
		X: 0,
		Y: 0,
	}
	var nBoids float64 = 0.0
	for _, agent := range agents {
		if agent.Type() == physics.BoidAgent && agent.ID() != b.ID() {
			nBoids++
			result.Add(agent.Velocity())
		}
	}
	if nBoids > 0 {
		result.Divide(nBoids)
		result.Multiply(boidMaxVelocity)
		result.Subtract(b.Velocity())
		result.Limit(boidMaxForce)
	}
	return result
}
