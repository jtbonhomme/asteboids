package agents

import (
	"time"

	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/jtbonhomme/asteboids/internal/sounds"
	"github.com/jtbonhomme/asteboids/internal/vector"
	"github.com/sirupsen/logrus"
)

const (
	aiRotationAngle      float64 = math.Pi / 36 // rotation of 5°
	aiMaxVelocity        float64 = 3.0
	aiAcceleration       float64 = 0.2
	avoidCollisionFactor float64 = 1.0
)

// AI is a PhysicalBody agent.
// It represents a playable star ship.
type AI struct {
	physics.Body
	lastBulletTime time.Time
	bulletImage    *ebiten.Image
	nearestAgents  []physics.Physic
}

// NewAI creates a new AI (PhysicalBody agent)
func NewAI(
	log *logrus.Logger,
	x, y,
	screenWidth, screenHeight float64,
	cbr physics.AgentRegister,
	cbu physics.AgentUnregister,
	vision physics.AgentVision,
	aiImage *ebiten.Image,
	bulletImage *ebiten.Image,
	debug bool) *AI {
	s := AI{
		lastBulletTime: time.Now(),
	}
	s.AgentType = physics.AIAgent
	s.Register = cbr
	s.Unregister = cbu
	s.Vision = vision
	s.Init(vector.Vector2D{
		X: 0,
		Y: 0,
	})
	s.LimitVelocity(aiMaxVelocity)
	s.Orientation = math.Pi / 2
	s.Move(vector.Vector2D{
		X: x,
		Y: y,
	})
	s.Accelerate(vector.Vector2D{
		X: 0,
		Y: 0,
	})
	s.PhysicWidth = 50
	s.PhysicHeight = 50
	s.VisionRadius = 250
	s.ScreenWidth = screenWidth
	s.ScreenHeight = screenHeight
	s.Log = log

	s.Image = aiImage
	s.bulletImage = bulletImage

	s.Debug = debug
	return &s
}

func intersect(pa vector.Vector2D, sa physics.Size, pb vector.Vector2D, sb physics.Size) bool {
	ax, ay := pa.X, pa.Y
	aw, ah := sa.W, sa.H

	bx, by := pb.X, pb.Y
	bw, bh := sb.W, sb.H

	return ((bx-ax)*(bx-ax) < (bw-aw)*(bw-aw) && (by-ay)*(by-ay) < (bh-ah)*(bh-ah))
}

func futureIntersect(agent physics.Physic, agents []physics.Physic) bool {
	// for each enemy agent in visibility range
	for _, a := range agents {
		if a.Type() != physics.AsteroidAgent &&
			a.Type() != physics.RubbleAgent {
			continue
		}
		// check if there is an intersect in the next 10 time unit
		for i := 0; i < 10; i++ {
			if intersect(agent.FuturePosition(float64(i)), agent.Dimension(),
				a.FuturePosition(float64(i)), a.Dimension()) {
				return true
			}
		}
	}
	return false
}

func (s *AI) targetLocked(agents []physics.Physic) bool {
	// create a virtual bullet to simulate a shot
	bullet := NewBullet(s.Log,
		s.Position().X, s.Position().Y,
		s.Orientation,
		s.ScreenWidth,
		s.ScreenHeight,
		s.Unregister,
		s.bulletImage)
	return futureIntersect(bullet, agents)
}

func (s *AI) avoidCollision(agents []physics.Physic) vector.Vector2D {
	result := vector.Vector2D{
		X: 0,
		Y: 0,
	}
	return result
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (s *AI) Update() {
	acceleration := vector.Vector2D{}
	nearestAgent := s.Vision(s.Position().X, s.Position().Y, s.VisionRadius)

	if s.targetLocked(nearestAgent) {
		s.Shot()
	}

	avoidCollision := s.avoidCollision(nearestAgent)
	avoidCollision.Multiply(avoidCollisionFactor)
	acceleration.Add(avoidCollision)

	s.Accelerate(acceleration)

	s.UpdateVelocity()
	s.UpdateOrientation()
	s.UpdatePosition()
	s.nearestAgents = s.Vision(s.Position().X, s.Position().Y, s.VisionRadius)
}

// Shot adds a new bullet to the game.
func (s *AI) Shot() {
	// throtlle call to avoid continuous shooting
	if time.Since(s.lastBulletTime) < bulletThrottle {
		return
	}
	s.lastBulletTime = time.Now()

	bullet := NewBullet(s.Log,
		s.Position().X, s.Position().Y,
		s.Orientation,
		s.ScreenWidth,
		s.ScreenHeight,
		s.Unregister,
		s.bulletImage)
	s.Register(bullet)
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (s *AI) Draw(screen *ebiten.Image) {
	defer s.Body.Draw(screen)
	s.LinkAgents(screen, s.nearestAgents, []string{physics.AsteroidAgent, physics.RubbleAgent})
}

// Explode proceeds the rubble termination.
func (s *AI) Explode() {
	s.Unregister(s.ID(), s.Type())
	go func() {
		_ = sounds.BangLargePlayer.Rewind()
		sounds.BangLargePlayer.Play()
	}()
}
