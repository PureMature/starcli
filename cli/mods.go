package cli

import (
	"fmt"
	"sort"

	"github.com/1set/starbox"
	"github.com/1set/starlet"
	"github.com/PureMature/starcli/module/sys"
	"github.com/samber/lo"
)

var (
	starMods = starlet.GetAllBuiltinModuleNames()
	cliMods  = []string{
		sys.ModuleName,
	}
)

// getDefaultModules returns the default modules for CLI, including builtin modules from Starlet and local modules in CLI.
func getDefaultModules() []string {
	allMods := lo.Union(starMods, cliMods)
	sort.Strings(allMods)
	return allMods
}

// loadModules loads the given modules into the Starbox instance.
func loadModules(box *starbox.Starbox, opts *BoxOpts, usrMods []string) error {
	if len(usrMods) == 0 {
		// no modules to load
		log.Debugw("no modules to load", "user_modules", usrMods)
		return nil
	}

	// inspect user modules
	allMods := getDefaultModules()
	unloadStar, unknown := lo.Difference(allMods, usrMods)
	log.Debugw("inspect user mods", "all_modules", allMods, "user_modules", usrMods, "unload_star_modules", unloadStar, "unknown_modules", unknown)
	if len(unknown) > 0 {
		return fmt.Errorf("unknown module(s): %v", unknown)
	}

	// load star* modules
	selectStarMods := lo.Intersect(starMods, usrMods)
	log.Debugw("selected star modules", "star_modules", selectStarMods)
	if len(selectStarMods) > 0 {
		box.AddNamedModules(selectStarMods...)
	}

	// load cli modules
	selectCLIMods := lo.Intersect(cliMods, usrMods)
	log.Debugw("selected cli modules", "cli_modules", selectCLIMods)
	if len(selectCLIMods) > 0 {
		for _, name := range selectCLIMods {
			ml, err := loadCLIModuleByName(opts, name)
			if err != nil {
				return err
			}
			box.AddModuleLoader(name, ml)
		}
	}

	// all is well
	return nil
}

func loadCLIModuleByName(opts *BoxOpts, name string) (starlet.ModuleLoader, error) {
	switch name {
	case sys.ModuleName:
		return sys.NewModule(opts.cmdArgs), nil
	default:
		return nil, fmt.Errorf("unknown module: %s", name)
	}
}
