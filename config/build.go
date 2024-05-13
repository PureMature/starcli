// Package config provides the build information, configuration for the application.
package config

// This file contains the build information for the application.

var (
	// AppName is the name of the application.
	AppName = "starcli"
	// CIBuildNum is the build number of the application on Continuous Integration system.
	CIBuildNum string
	// BuildDate is the date of the application build.
	BuildDate string
	// BuildHost is the host/machine name of the application build.
	BuildHost string
	// GoVersion is the version of Go used to build the application. i.e. go version
	GoVersion string
	// GitBranch is the git branch name of the application source code. i.e. git symbolic-ref -q --short HEAD
	GitBranch string
	// GitCommit is the git commit hash of the application source code. i.e. git rev-parse --short HEAD
	GitCommit string
	// GitSummary is the git summary of the application source code. i.e. git describe --tags --dirty --always
	GitSummary string
)

var (
	// AppLogoArt is the ASCII art of the application logo. Use `kfiglet starcli -f ansishadow -p 2 | pbcopy` to generate.
	AppLogoArt = `
  ███████╗████████╗ █████╗ ██████╗  ██████╗██╗     ██╗
  ██╔════╝╚══██╔══╝██╔══██╗██╔══██╗██╔════╝██║     ██║
  ███████╗   ██║   ███████║██████╔╝██║     ██║     ██║
  ╚════██║   ██║   ██╔══██║██╔══██╗██║     ██║     ██║
  ███████║   ██║   ██║  ██║██║  ██║╚██████╗███████╗██║
  ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚══════╝╚═╝
`
)
