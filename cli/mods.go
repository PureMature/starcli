package cli

import (
	"sort"

	"github.com/1set/starlet"
	"github.com/PureMature/starcli/module/sys"
	"github.com/samber/lo"
)

var (
	cliMods = []string{
		sys.ModuleName,
	}
)

// getDefaultModules returns the default modules for CLI, including builtin modules from Starlet and local modules in CLI.
func getDefaultModules() []string {
	starMods := starlet.GetAllBuiltinModuleNames()
	allMods := lo.Union(starMods, cliMods)
	sort.Strings(allMods)
	return allMods
}
