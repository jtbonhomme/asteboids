package main

import (
	"os"
	"runtime/pprof"

	"github.com/dimiro1/banner"
	"github.com/jtbonhomme/asteboids"
	"github.com/jtbonhomme/asteboids/internal/config"
	"github.com/jtbonhomme/asteboids/internal/version"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	conf := config.New()

	if conf.CPUProfile != "" {
		f, err := os.Create(conf.CPUProfile)
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	isEnabled := true
	isColorEnabled := true
	templ := `{{ .Title "Asteboids" "" 4 }}
{{ .AnsiColor.BrightCyan }}Classical Asteroids game with Boids inside !{{ .AnsiColor.Default }}
GoVersion: {{ .GoVersion }}
Now: {{ .Now "Monday, 2 Jan 2006" }}
Asteboids Version: ` + version.Read().Tag + `

`
	banner.InitString(colorable.NewColorableStdout(), isEnabled, isColorEnabled, templ)

	if conf.Debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	log.Infof("config: %#v", conf)
	err := asteboids.Run(log, conf)
	if err != nil {
		log.Panic("error while running asteboids: ", err)
	}
}
