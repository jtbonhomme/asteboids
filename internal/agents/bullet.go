package agents

import (
	"image"
	"image/color"
	"math"

	// anonymous import for png decoder
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/game"
	"github.com/sirupsen/logrus"
)

const (
	bulletVelocity float64 = 18.0
	bulletTTL      int     = 30
	deltaAngle     float64 = 20.0
)

// Bullet is a PhysicalBody agent
// It represents a bullet shot by a starship agent.
type Bullet struct {
	game.PhysicBody
	vertices []ebiten.Vertex
	ttl      int
}

// NewBullet creates a new Bullet (PhysicalBody agent)
func NewBullet(log *logrus.Logger, x, y int, orientation float64, screenWidth, screenHeight int, cb game.AgentUnregister) *Bullet {
	b := Bullet{
		vertices: []ebiten.Vertex{},
		ttl:      bulletTTL,
	}
	b.Unregister = cb

	b.Init()
	b.Log = log

	b.Orientation = orientation
	b.Velocity = game.Vector{
		X: bulletVelocity * math.Cos(b.Orientation),
		Y: bulletVelocity * math.Sin(b.Orientation),
	}
	b.Size = 3
	b.PhysicWidth = 8
	b.PhysicHeight = 8
	b.ScreenWidth = screenWidth
	b.ScreenHeight = screenHeight
	b.X = x
	b.Y = y

	b.updateVertices()
	b.Image = ebiten.NewImage(3, 3)
	b.Image.Fill(color.White)
	return &b
}

func (b *Bullet) updateVertices() {
	vs := []ebiten.Vertex{}
	for i := 0; i < 4; i++ {
		vs = append(vs, ebiten.Vertex{
			DstX:   0,
			DstY:   0,
			SrcX:   0,
			SrcY:   0,
			ColorR: 255,
			ColorG: 255,
			ColorB: 255,
			ColorA: 255,
		})
	}
	centerX := b.X
	centerY := b.Y
	// bullet vertices
	vs[0].DstX = float32(centerX + int(b.Size*math.Cos(b.Orientation+math.Pi/deltaAngle)))
	vs[0].DstY = float32(centerY + int(b.Size*math.Sin(b.Orientation+math.Pi/deltaAngle)))
	vs[1].DstX = float32(centerX + int(b.Size*math.Cos(b.Orientation-math.Pi/deltaAngle)))
	vs[1].DstY = float32(centerY + int(b.Size*math.Sin(b.Orientation-math.Pi/deltaAngle)))
	vs[2].DstX = float32(centerX + int(b.Size*math.Cos(b.Orientation+math.Pi+math.Pi/deltaAngle)))
	vs[2].DstY = float32(centerY + int(b.Size*math.Sin(b.Orientation+math.Pi+math.Pi/deltaAngle)))
	vs[3].DstX = float32(centerX + int(b.Size*math.Cos(b.Orientation+math.Pi-math.Pi/deltaAngle)))
	vs[3].DstY = float32(centerY + int(b.Size*math.Sin(b.Orientation+math.Pi-math.Pi/deltaAngle)))

	b.vertices = vs
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
// Update maintains a TTL counter to limit live of bullets.
func (b *Bullet) Update() {
	b.ttl--
	if b.ttl == 0 {
		b.SelfDestroy()
	}
	// update position
	b.X += int(b.Velocity.X)
	b.Y += int(b.Velocity.Y)

	if b.X > b.ScreenWidth {
		b.X = 0
	} else if b.X < 0 {
		b.X = b.ScreenWidth
	}
	if b.Y > b.ScreenHeight {
		b.Y = 0
	} else if b.Y < 0 {
		b.Y = b.ScreenHeight
	}
	b.updateVertices()
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (b *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawTrianglesOptions{}
	op.Address = ebiten.AddressUnsafe
	screen.DrawTriangles(b.vertices, []uint16{0, 1, 2, 1, 2, 3}, b.Image.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)
}

// SelfDestroy removes the agent from the game
func (b *Bullet) SelfDestroy() {
	defer b.Unregister(b.ID())
}
