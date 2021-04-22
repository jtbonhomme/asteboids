package physics

import (
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/vector"
	"github.com/sirupsen/logrus"
)

type Body struct {
	position    vector.Vector2D
	AgentType   string
	id          uuid.UUID
	Log         *logrus.Logger
	Orientation float64 // theta (radian)

	PhysicWidth  float64
	PhysicHeight float64
	ScreenWidth  float64
	ScreenHeight float64

	velocity     vector.Vector2D
	maxVelocity  float64
	acceleration vector.Vector2D

	Register   AgentRegister
	Unregister AgentUnregister
	Vision     AgentVision
	Image      *ebiten.Image

	Debug bool
}

// Init initializes the physic body
func (pb *Body) Init(velocity vector.Vector2D) {
	pb.id = uuid.New()
	pb.velocity = velocity
	pb.maxVelocity = defaultMaxVelocity
}

// Draw draws the agent.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (pb *Body) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	defer screen.DrawImage(pb.Image, op)

	op.GeoM.Translate(-pb.PhysicWidth/2, -pb.PhysicHeight/2)
	op.GeoM.Rotate(pb.Orientation)
	op.GeoM.Translate(pb.position.X, pb.position.Y)
	if pb.Debug {
		pb.DrawBodyBoundaryBox(screen)
	}
}
