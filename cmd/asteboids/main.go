package main

import (
	"github.com/dimiro1/banner"
	"github.com/jtbonhomme/asteboids"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

func main() {
	isEnabled := true
	isColorEnabled := true
	templ := `{{ .Title "Asteboids" "" 4 }}
{{ .AnsiColor.BrightCyan }}The title will be ascii and indented 4 spaces{{ .AnsiColor.Default }}
GoVersion: {{ .GoVersion }}
Now: {{ .Now "Monday, 2 Jan 2006" }}

`
	banner.InitString(colorable.NewColorableStdout(), isEnabled, isColorEnabled, templ)
	log := logrus.New()
	asteboids.Run(log)
}
