package config

import "github.com/jtbonhomme/conf"

type Config struct {
	Debug                      bool    `conf:"debug" help:"Debug log level activated."`
	Optim                      bool    `conf:"optim" help:"Optimized mode activated."`
	CPUProfile                 string  `conf:"cpuprofile" help:"Write CPU profile to file."`
	Asteroids                  int     `conf:"asteroids" help:"Number of asteroids at the start of the game."`
	Boids                      int     `conf:"boids" help:"Number of boids at the start of the game."`
	ScreenWidth                float64 `conf:"screenWidth" help:"Screen width (in pixels)."`
	ScreenHeight               float64 `conf:"screenHeight" help:"Screen height (in pixels)."`
	ScoreTimeUnit              float64 `conf:"scoreTimeUnit" help:"Time delay (in second) to win one point, just being still alive."`
	AutoGenerateAsteroidsRatio float64 `conf:"autoGenerateAsteroidsRatio" help:"Time delay (in second) before a new asteroids spawn."`
}

func New() *Config {
	config := &Config{}
	conf.Load(config)
	return config
}
