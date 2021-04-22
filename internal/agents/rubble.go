package agents

import (
	"image/color"
	"math"
	"math/rand"

	// anonymous import for png decoder
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jtbonhomme/asteboids/internal/fonts"
	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/jtbonhomme/asteboids/internal/vector"
	"github.com/sirupsen/logrus"
)

const (
	rubbleMaxVelocity   float64 = 2.0
	rubbleRotationSpeed float64 = 0.07
)

// Rubble is a PhysicalBody agent
// It represents a bullet shot by a starship agent.
type Rubble struct {
	physics.Body
}

// NewRubble creates a new Rubble (PhysicalBody agent)
func NewRubble(log *logrus.Logger,
	x, y,
	screenWidth, screenHeight float64,
	cbu physics.AgentUnregister,
	rubbleImage *ebiten.Image,
	debug bool) *Rubble {
	r := Rubble{}
	r.AgentType = physics.RubbleAgent
	r.Unregister = cbu

	r.Orientation = math.Pi / 16 * float64(rand.Intn(32))

	r.Init(vector.Vector2D{
		X: rubbleMaxVelocity * math.Cos(r.Orientation),
		Y: rubbleMaxVelocity * math.Sin(r.Orientation),
	})
	r.Log = log
	r.LimitVelocity(rubbleMaxVelocity)

	r.Move(vector.Vector2D{
		X: x,
		Y: y,
	})
	r.PhysicWidth = 50
	r.PhysicHeight = 50
	r.ScreenWidth = screenWidth
	r.ScreenHeight = screenHeight

	/*filename := fmt.Sprintf("./resources/images/rubble%d.png", rand.Intn(5))*/
	r.Image = rubbleImage
	r.Debug = debug
	return &r
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (r *Rubble) Update() {
	defer r.Body.UpdatePosition()
	r.Rotate(rubbleRotationSpeed)
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (r *Rubble) Draw(screen *ebiten.Image) {
	defer r.Body.Draw(screen)

	if r.Debug {
		msg := r.String()
		textDim := text.BoundString(fonts.MonoSansRegularFont, msg)
		textWidth := textDim.Max.X - textDim.Min.X
		text.Draw(screen, msg, fonts.MonoSansRegularFont, int(r.Position().X)-textWidth/2, int(r.Position().Y+r.PhysicHeight/2+5), color.Gray16{0x999f})
	}
}

// Explode proceeds the rubble termination.
func (r *Rubble) Explode() {
	r.Unregister(r.ID(), r.Type())
}
