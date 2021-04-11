package version

var (
	// GitCommit will be overwritten automatically by the build system
	GitCommit = "HEAD"
	// Tag will be overwritten automatically by the build system
	Tag = "0.0.0"
	// BuildTime will be overwritten automatically by the build system
	BuildTime = "unknown"
)

// Version is a strucuture to describe current version
type Version struct {
	// The current git commit (short format)
	GitCommit string `json:"gitCommit"`
	// The current git tag (i exist).
	Tag string `json:"tag"`
	// The build time.
	Buildtime string `json:"buildTime"`
}

func Read() *Version {
	return &Version{
		Buildtime: BuildTime,
		GitCommit: GitCommit,
		Tag:       Tag,
	}
}
