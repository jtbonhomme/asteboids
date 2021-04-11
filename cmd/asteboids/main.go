package main

import (
	"github.com/dimiro1/banner"
	"github.com/jtbonhomme/asteboids"
	"github.com/jtbonhomme/asteboids/internal/version"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

func main() {
	isEnabled := true
	isColorEnabled := true
	templ := `{{ .Title "Asteboids" "" 4 }}
{{ .AnsiColor.BrightCyan }}Classical Asteroids game with Boids inside !{{ .AnsiColor.Default }}
GoVersion: {{ .GoVersion }}
Now: {{ .Now "Monday, 2 Jan 2006" }}
Asteboids Version: ` + version.Read().Tag + `

`
	banner.InitString(colorable.NewColorableStdout(), isEnabled, isColorEnabled, templ)
	log := logrus.New()
	// Only log the warning severity or above.
	log.SetLevel(logrus.DebugLevel)

	err := asteboids.Run(log)
	if err != nil {
		log.Fatalf("error while running asteboids: %s", err.Error())
	}
}
