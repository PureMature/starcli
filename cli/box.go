package cli

import (
	"os"

	"github.com/1set/gut/ystring"
	"github.com/1set/starbox"
	"github.com/1set/starlet"
	"github.com/PureMature/starcli/module/sys"
	"github.com/PureMature/starcli/util"
)

// BuildBox creates a new Starbox with the given name, include path, and modules.
func BuildBox(name, inclPath string, modules, args []string) *starbox.Starbox {
	// create a new Starbox instance
	box := starbox.New(name)
	box.AddNamedModules(modules...)
	if ystring.IsNotBlank(inclPath) {
		box.SetFS(os.DirFS(inclPath))
	}
	// add default modules
	box.AddModuleLoader(sys.ModuleName, sys.NewModule(args))
	return box
}

// genInspectCond creates a function for Starbox runner to inspect the result.
func genInspectCond(inspect bool) starbox.InspectCondFunc {
	if inspect {
		return func(m starlet.StringAnyMap, err error) bool {
			if err != nil {
				util.PrintError(err)
			}
			return true
		}
	}
	return func(starlet.StringAnyMap, error) bool {
		return false
	}
}
