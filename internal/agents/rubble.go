package agents

import (
	"crypto/rand"
	"fmt"
	"image/color"
	"math"
	"math/big"

	// anonymous import for png decoder
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jtbonhomme/asteboids/internal/fonts"
	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/sirupsen/logrus"
)

const (
	rubbleVelocity      float64 = 3.5
	rubbleRotationSpeed float64 = 0.07
)

// Rubble is a PhysicalBody agent
// It represents a bullet shot by a starship agent.
type Rubble struct {
	physics.PhysicBody
}

// NewRubble creates a new Rubble (PhysicalBody agent)
func NewRubble(log *logrus.Logger,
	x, y,
	screenWidth, screenHeight float64,
	cbu physics.AgentUnregister,
	debug bool) *Rubble {
	r := Rubble{}
	r.AgentType = physics.RubbleAgent
	r.Unregister = cbu

	r.Init()
	r.Log = log

	nBig, err := rand.Int(rand.Reader, big.NewInt(32))
	if err != nil {
		r.Log.Fatal(err)
	}
	r.Orientation = math.Pi / 16 * float64(nBig.Int64())
	r.Velocity = physics.Vector{
		X: rubbleVelocity * math.Cos(r.Orientation),
		Y: rubbleVelocity * math.Sin(r.Orientation),
	}
	r.Size = 3
	r.PhysicWidth = 50
	r.PhysicHeight = 50
	r.ScreenWidth = screenWidth
	r.ScreenHeight = screenHeight
	r.X = x
	r.Y = y

	nLittle, err := rand.Int(rand.Reader, big.NewInt(5))
	if err != nil {
		r.Log.Fatal(err)
	}
	filename := fmt.Sprintf("./resources/images/rubble%d.png", int(nLittle.Int64()))
	err = r.LoadImage(filename)
	if err != nil {
		log.Errorf("error when loading image from file: %s", err.Error())
	}
	r.Debug = debug
	return &r
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
// Update maintains a TTL counter to limit live of bullets.
func (r *Rubble) Update() {
	r.Rotate(rubbleRotationSpeed)
	// update position
	r.X += r.Velocity.X
	r.Y += r.Velocity.Y

	if r.X > r.ScreenWidth {
		r.X = 0
	} else if r.X < 0 {
		r.X = r.ScreenWidth
	}
	if r.Y > r.ScreenHeight {
		r.Y = 0
	} else if r.Y < 0 {
		r.Y = r.ScreenHeight
	}
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (r *Rubble) Draw(screen *ebiten.Image) {
	defer r.PhysicBody.Draw(screen)

	if r.Debug {
		msg := r.String()
		textDim := text.BoundString(fonts.MonoSansRegularFont, msg)
		textWidth := textDim.Max.X - textDim.Min.X
		text.Draw(screen, msg, fonts.MonoSansRegularFont, int(r.X)-textWidth/2, int(r.Y+r.PhysicHeight/2+5), color.Gray16{0x999f})
	}
}

// Explode proceeds the rubble termination.
func (r *Rubble) Explode() {
	r.Unregister(r.ID(), r.Type())
}
