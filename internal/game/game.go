package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jtbonhomme/asteboids/internal/fonts"
	"github.com/sirupsen/logrus"
)

const (
	defaultScreenWidth  = 1080
	defaultScreenHeight = 720
)

type Game struct {
	log             *logrus.Logger
	ScreenWidth     int
	ScreenHeight    int
	backgroundColor color.RGBA
	agents          map[string]Physic
	StarshipID      string
}

func New(log *logrus.Logger) *Game {
	log.Infof("New Game")
	return &Game{
		log:             log,
		ScreenWidth:     defaultScreenWidth,
		ScreenHeight:    defaultScreenHeight,
		backgroundColor: color.RGBA{0x00, 0x00, 0x00, 0xff},
		agents:          make(map[string]Physic),
	}
}

// Register adds a new agent (player or ai) to the game.
func (g *Game) Register(agent Physic) {
	g.agents[agent.ID()] = agent
}

// Unregister deletes an agent (player or ai) from the game.
func (g *Game) Unregister(id string) {
	delete(g.agents, id)
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	// Update the agents
	for _, a := range g.agents {
		a.Update()
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	// Draw the ground image.
	screen.Fill(g.backgroundColor)

	// Title
	title := "Asteboids"
	textDim := text.BoundString(fonts.FurturisticRegularFontTitle, title)
	textWidth := textDim.Max.X - textDim.Min.X
	textHeight := textDim.Max.Y - textDim.Min.Y
	text.Draw(screen, title, fonts.FurturisticRegularFontTitle, g.ScreenWidth/2-textWidth/2, textHeight/2+20, color.Gray16{0xffff})

	// Draw the agents
	for _, a := range g.agents {
		a.Draw(screen)
	}
	// Draw the message.
	usage := "s: take a screenshot\nCmd + q: exit"
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	msg += fmt.Sprintf("%s\n", usage)
	msg += fmt.Sprintf("%s\n", g.agents[g.StarshipID])
	ebitenutil.DebugPrint(screen, msg)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenWidth, g.ScreenHeight
}

func (g *Game) String() string {
	return fmt.Sprintf(`Asteboids
	- screen size: %d x %d
`, g.ScreenWidth, g.ScreenHeight)
}
