package main

import (
	"flag"
	"os"
	"runtime/pprof"

	"github.com/dimiro1/banner"
	"github.com/jtbonhomme/asteboids"
	"github.com/jtbonhomme/asteboids/internal/version"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	debug := flag.Bool("debug", false, "Debug log level")
	optim := flag.Bool("optim", false, "Optimized mode")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
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

	if *debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	err := asteboids.Run(log, *optim, *debug)
	if err != nil {
		log.Panic("error while running asteboids: ", err)
	}
}
