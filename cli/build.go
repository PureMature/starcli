package cli

import (
	"os"

	"github.com/1set/gut/ystring"
	"github.com/1set/starbox"
	"github.com/PureMature/starcli/module/sys"
)

// BuildBox creates a new Starbox with the given name, include path, and modules.
func BuildBox(name string, args *Args) *starbox.Starbox {
	// create a new Starbox instance
	box := starbox.New(name)
	box.AddNamedModules(args.LoadModules...)
	if p := args.IncludePath; ystring.IsNotBlank(p) {
		box.SetFS(os.DirFS(p))
	}
	// add default modules
	box.AddModuleLoader(sys.ModuleName, sys.NewModule(args.Arguments))
	return box
}
