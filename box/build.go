package box

import (
	"os"

	"github.com/1set/gut/ystring"
	"github.com/1set/starbox"
)

// Build creates a new Starbox with the given name, include path, and modules.
func Build(name, includePath string, modules []string) *starbox.Starbox {
	box := starbox.New(name)
	box.AddNamedModules(modules...)
	if ystring.IsNotBlank(includePath) {
		box.SetFS(os.DirFS(includePath))
	}
	return box
}
