package game

import (
	"math/rand"
	"time"

	"github.com/jtbonhomme/asteboids/internal/physics"
)

// UpdateAgents loops over all game agents to update them
func (g *Game) UpdateAgents() {
	for _, b := range g.bullets {
		b.Update()
	}
	for _, a := range g.asteroids {
		a.Update()
	}
	for _, s := range g.starships {
		s.Update()
	}
	for _, b := range g.boids {
		b.Update()
	}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Update the agents
	g.UpdateAgents()

	// detect starship collision with asteroids
	for _, starship := range g.starships {
		_, ok := starship.IntersectMultiple(g.asteroids)
		if ok {
			starship.Explode()
		}
	}

	// detect asteroid collision with bullet
	for _, asteroid := range g.asteroids {
		bID, ok := asteroid.IntersectMultiple(g.bullets)
		if ok {
			asteroid.Explode()
			asteroidType := asteroid.Type()
			delete(g.bullets, bID)
			g.kills++
			// Only add a new asteroids if the destroyed agent is also an asteroid (not a rubble)
			if asteroidType == physics.AsteroidAgent {
				g.AddAsteroid(g.asteroidImages[rand.Intn(5)])
			}
		}
	}

	// game ends when there is no starship left
	if len(g.starships) == 0 {
		g.gameOver = true
		g.gameWon = false
	}

	// update time until game ends
	if !g.gameOver {
		g.gameDuration = time.Since(g.startTime).Round(time.Second)
	}
	// periodically add new asteroids
	if int(g.gameDuration.Seconds()/autoGenerateAsteroidsRatio) > len(g.asteroids) {
		g.AddAsteroid(g.asteroidImages[rand.Intn(5)])
	}
	return nil
}
