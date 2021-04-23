package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jtbonhomme/asteboids/internal/fonts"
)

// DrawAgents loops over all game agents to update them
func (g *Game) DrawAgents(screen *ebiten.Image) {
	for _, s := range g.starships {
		s.Draw(screen)
	}
	for _, a := range g.asteroids {
		a.Draw(screen)
	}
	for _, b := range g.bullets {
		b.Draw(screen)
	}
	for _, b := range g.boids {
		b.Draw(screen)
	}
}

func (g *Game) Score() int {
	return int(g.gameDuration.Seconds()/g.conf.ScoreTimeUnit) + g.kills*2
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Erase the image.
	screen.Fill(g.backgroundColor)

	// Draw the agents
	g.DrawAgents(screen)

	if g.debug {
		nAgents := len(g.asteroids) + len(g.starships) + len(g.bullets) + len(g.boids)
		msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nAgents: %d", ebiten.CurrentTPS(), ebiten.CurrentFPS(), nAgents)
		ebitenutil.DebugPrint(screen, msg)
	}

	g.drawTimeElapsed(screen)

	g.drawScore(screen)
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.RestartGame()
		}
		// Title
		title := "Asteboids"
		titleTextDim := text.BoundString(fonts.FurturisticRegularFontTitle, title)
		titleTextWidth := titleTextDim.Max.X - titleTextDim.Min.X
		text.Draw(
			screen,
			title,
			fonts.FurturisticRegularFontTitle,
			int(g.conf.ScreenWidth/2)-titleTextWidth/2,
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
			int(g.conf.ScreenWidth/2)-gameOverTextWidth/2,
			int(g.conf.ScreenHeight/2)-gameOverTextHeight/2,
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
			int(g.conf.ScreenWidth/2)-replayTextWidth/2,
			int(g.conf.ScreenHeight/2)+gameOverTextHeight/2+replayTextHeight/2,
			color.Gray16{0xbbbf},
		)
	}
}

func (g *Game) drawScore(screen *ebiten.Image) {
	// Score
	score := fmt.Sprintf("Score %d", g.Score())
	scoreTextDim := text.BoundString(fonts.FurturisticRegularFontMenu, score)
	scoreTextHeight := scoreTextDim.Max.Y - scoreTextDim.Min.Y
	text.Draw(
		screen,
		score,
		fonts.FurturisticRegularFontMenu,
		900,
		scoreTextHeight+10,
		color.Gray16{0xffff},
	)
}

func (g *Game) drawTimeElapsed(screen *ebiten.Image) {
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
}
