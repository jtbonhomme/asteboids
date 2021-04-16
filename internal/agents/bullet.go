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
	bulletVelocity float64 = 5.5
	bulletTTL      int     = 100
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

	b.Orientation = orientation
	b.Velocity = game.Vector{
		X: bulletVelocity * math.Cos(b.Orientation),
		Y: bulletVelocity * math.Sin(b.Orientation),
	}
	b.Size = 2
	b.PhysicWidth = 5
	b.PhysicHeight = 5
	b.ScreenWidth = screenWidth
	b.ScreenHeight = screenHeight
	b.X = x
	b.Y = y
	b.Log = log

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
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		})
	}
	centerX := b.X
	centerY := b.Y
	// bullet vertices
	vs[0].DstX = float32(centerX + int(b.Size*math.Cos(b.Orientation+math.Pi/16)))
	vs[0].DstY = float32(centerY + int(b.Size*math.Sin(b.Orientation+math.Pi/16)))
	vs[1].DstX = float32(centerX + int(b.Size*math.Cos(b.Orientation-math.Pi/16)))
	vs[1].DstY = float32(centerY + int(b.Size*math.Sin(b.Orientation-math.Pi/16)))
	vs[2].DstX = float32(centerX + int(b.Size*math.Cos(b.Orientation+math.Pi+math.Pi/16)))
	vs[2].DstY = float32(centerY + int(b.Size*math.Sin(b.Orientation+math.Pi+math.Pi/16)))
	vs[3].DstX = float32(centerX + int(b.Size*math.Cos(b.Orientation+math.Pi-math.Pi/16)))
	vs[3].DstY = float32(centerY + int(b.Size*math.Sin(b.Orientation+math.Pi-math.Pi/16)))

	b.vertices = vs
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
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
