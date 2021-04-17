package game

import (
	"crypto/rand"
	"fmt"
	"image/color"
	"math/big"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jtbonhomme/asteboids/internal/agents"
	"github.com/jtbonhomme/asteboids/internal/fonts"
	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/sirupsen/logrus"
)

const (
	defaultScreenWidth  = 1080
	defaultScreenHeight = 720
	scoreTimeUnit       = 5
)

type Game struct {
	log             *logrus.Logger
	gameOver        bool
	gameWon         bool
	nAsteroids      int
	startTime       time.Time
	gameDuration    time.Duration
	highestDuration time.Duration
	highScore       int
	kills           int
	debug           bool
	ScreenWidth     int
	ScreenHeight    int
	backgroundColor color.RGBA
	starships       map[string]physics.Physic
	asteroids       map[string]physics.Physic
	bullets         map[string]physics.Physic
}

func New(log *logrus.Logger, nAsteroids int, debug bool) *Game {
	log.Infof("New Game")
	return &Game{
		log:             log,
		gameOver:        false,
		gameWon:         false,
		nAsteroids:      nAsteroids,
		startTime:       time.Now(),
		gameDuration:    0,
		kills:           0,
		highScore:       0,
		highestDuration: 0,
		debug:           debug,
		ScreenWidth:     defaultScreenWidth,
		ScreenHeight:    defaultScreenHeight,
		backgroundColor: color.RGBA{0x00, 0x00, 0x00, 0xff},
		starships:       make(map[string]physics.Physic),
		asteroids:       make(map[string]physics.Physic),
		bullets:         make(map[string]physics.Physic),
	}
}

// Start initializes a new game.
func (g *Game) Start() {
	// add starship
	p := agents.NewStarship(g.log, g.ScreenWidth/2, g.ScreenHeight/2, g.ScreenWidth, g.ScreenHeight, g.Register, g.Unregister, g.debug)
	g.log.Infof("added starship: %s", p.ID())
	g.Register(p)

	// add asteroids
	for i := 0; i < g.nAsteroids; i++ {
		g.AddAsteroid()
	}
	g.startTime = time.Now()
	g.gameDuration = 0
	g.gameOver = false
	g.gameWon = false
	g.kills = 0
}

// AddAsteroid insert a new asteroid in the game.
func (g *Game) AddAsteroid() {
	nWidth, err := rand.Int(rand.Reader, big.NewInt(int64(g.ScreenWidth/2)))
	if err != nil {
		g.log.Fatal(err)
	}
	nHeight, err := rand.Int(rand.Reader, big.NewInt(int64(g.ScreenHeight/2)))
	if err != nil {
		g.log.Fatal(err)
	}
	a := agents.NewAsteroid(g.log,
		int(nWidth.Int64()),
		int(nHeight.Int64()),
		g.ScreenWidth, g.ScreenHeight,
		g.Register, g.Unregister,
		g.debug)
	g.Register(a)
}

// Restart cleans current game and a start a new game.
func (g *Game) Restart() {
	for k := range g.starships {
		delete(g.starships, k)
	}
	for k := range g.asteroids {
		delete(g.asteroids, k)
	}
	for k := range g.bullets {
		delete(g.bullets, k)
	}

	g.Start()
}

// Register adds a new agent (player or ai) to the game.
func (g *Game) Register(agent physics.Physic) {
	switch agent.Type() {
	case physics.StarshipAgent:
		g.starships[agent.ID()] = agent
	case physics.AsteroidAgent:
		g.asteroids[agent.ID()] = agent
	case physics.RubbleAgent:
		g.asteroids[agent.ID()] = agent
	case physics.BulletAgent:
		g.bullets[agent.ID()] = agent
	default:
	}
}

// Unregister deletes an agent (player or ai) from the game.
func (g *Game) Unregister(id, agentType string) {
	switch agentType {
	case physics.StarshipAgent:
		delete(g.starships, id)
	case physics.AsteroidAgent:
		delete(g.asteroids, id)
	case physics.RubbleAgent:
		delete(g.asteroids, id)
	case physics.BulletAgent:
		delete(g.bullets, id)
	default:
	}
}

// Agents returns a map that combine all registered agents
func (g *Game) Agents() map[string]physics.Physic {
	res := make(map[string]physics.Physic)
	for k, v := range g.starships {
		res[k] = v
	}
	for k, v := range g.asteroids {
		res[k] = v
	}
	for k, v := range g.bullets {
		res[k] = v
	}

	return res
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	// Update the agents
	for _, a := range g.Agents() {
		a.Update()
	}
	// Collision detection
	// Convert map to slice of values.
	asteroidsList := []physics.Physic{}
	for _, asteroid := range g.asteroids {
		asteroidsList = append(asteroidsList, asteroid)
	}
	bulletList := []physics.Physic{}
	for _, bullet := range g.bullets {
		bulletList = append(bulletList, bullet)
	}

	// detect starship collision with asteroids
	for _, starship := range g.starships {
		_, ok := starship.IntersectMultiple(asteroidsList)
		if ok {
			starship.Explode()
		}
	}

	// detect asteroid collision with bullet
	for _, asteroid := range g.asteroids {
		bID, ok := asteroid.IntersectMultiple(bulletList)
		if ok {
			asteroid.Explode()
			asteroidType := asteroid.Type()
			delete(g.bullets, bID)
			g.kills++
			// Only add a new asteroids if the destroyed agent is also an asteroid (not a rubble)
			if asteroidType == physics.AsteroidAgent {
				g.AddAsteroid()
			}
		}
	}
	if len(g.starships) == 0 {
		g.gameOver = true
		g.gameWon = false
	}

	if len(g.asteroids) == 0 {
		g.gameOver = true
		g.gameWon = true
	}

	if !g.gameOver {
		g.gameDuration = time.Since(g.startTime).Round(time.Second)
	}

	return nil
}

func (g *Game) Score() int {
	return int(g.gameDuration.Seconds()/scoreTimeUnit) + g.kills*2
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	// Draw the ground image.
	screen.Fill(g.backgroundColor)

	// Draw the agents
	for _, a := range g.Agents() {
		a.Draw(screen)
	}

	if g.debug {
		msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n", ebiten.CurrentTPS(), ebiten.CurrentFPS())
		ebitenutil.DebugPrint(screen, msg)
	}

	// Time elapsed
	elapsed := "Time elapsed " + g.gameDuration.String()
	elapsedTextDim := text.BoundString(fonts.FurturisticRegularFontMenu, elapsed)
	elapsedTextHeight := elapsedTextDim.Max.Y - elapsedTextDim.Min.Y
	text.Draw(
		screen,
		elapsed,
		fonts.FurturisticRegularFontMenu,
		100,
		elapsedTextHeight+10,
		color.Gray16{0xffff},
	)

	// Score
	score := fmt.Sprintf("Score %d", g.Score())
	scoreTextDim := text.BoundString(fonts.FurturisticRegularFontMenu, elapsed)
	scoreTextHeight := scoreTextDim.Max.Y - scoreTextDim.Min.Y
	text.Draw(
		screen,
		score,
		fonts.FurturisticRegularFontMenu,
		900,
		scoreTextHeight+10,
		color.Gray16{0xffff},
	)

	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.Restart()
		}
		// Title
		title := "Asteboids"
		titleTextDim := text.BoundString(fonts.FurturisticRegularFontTitle, title)
		titleTextWidth := titleTextDim.Max.X - titleTextDim.Min.X
		text.Draw(
			screen,
			title,
			fonts.FurturisticRegularFontTitle,
			g.ScreenWidth/2-titleTextWidth/2,
			100,
			color.Gray16{0xffff},
		)

		if g.gameDuration > g.highestDuration {
			g.highestDuration = g.gameDuration
		}
		if g.Score() > g.highScore {
			g.highScore = g.Score()
		}

		var gameOver string
		if g.gameWon {
			gameOver = "YOU WIN !"
		} else {
			gameOver = "GAME OVER"
		}

		gameOverTextDim := text.BoundString(fonts.KarmaticArcadeFont, gameOver)
		gameOverTextWidth := gameOverTextDim.Max.X - gameOverTextDim.Min.X
		gameOverTextHeight := gameOverTextDim.Max.Y - gameOverTextDim.Min.Y
		text.Draw(
			screen,
			gameOver,
			fonts.KarmaticArcadeFont,
			g.ScreenWidth/2-gameOverTextWidth/2,
			g.ScreenHeight/2-gameOverTextHeight/2,
			color.Gray16{0xffff},
		)

		replay := "press   enter   to   play  again"
		replayTextDim := text.BoundString(fonts.ArcadeClassicFont, replay)
		replayTextWidth := replayTextDim.Max.X - replayTextDim.Min.X
		replayTextHeight := replayTextDim.Max.Y - replayTextDim.Min.Y
		text.Draw(
			screen,
			replay,
			fonts.ArcadeClassicFont,
			g.ScreenWidth/2-replayTextWidth/2,
			g.ScreenHeight/2+gameOverTextHeight/2+replayTextHeight/2,
			color.Gray16{0xbbbf},
		)
	}
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
