package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/1set/gut/ystring"
	"github.com/1set/starbox"
	"github.com/1set/starlet"
	"github.com/PureMature/starcli/module/sys"
	"github.com/PureMature/starcli/util"
	"go.starlark.net/starlark"
	"go.uber.org/atomic"
)

type scenarioCode uint

const (
	scenarioREPL scenarioCode = iota + 1
	scenarioDirect
	scenarioFile
	scenarioWeb
)

// BoxOpts defines the options for creating a new Starbox instance.
type BoxOpts struct {
	scenario     scenarioCode
	name         string
	includePath  string
	moduleToLoad []string
	cmdArgs      []string
	printerName  string
}

// BuildBox creates a new Starbox with the given options.
func BuildBox(opts *BoxOpts) (*starbox.Starbox, error) {
	// create a new Starbox instance
	box := starbox.New(opts.name)
	box.AddNamedModules(opts.moduleToLoad...)
	if ystring.IsNotBlank(opts.includePath) {
		box.SetFS(os.DirFS(opts.includePath))
	}
	box.SetPrintFunc(getPrinterFunc(opts.name, opts.printerName))
	// add default modules
	box.AddModuleLoader(sys.ModuleName, sys.NewModule(opts.cmdArgs))
	return box, nil
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

// getPrinterFunc returns a function to print output based on the given printer name.
func getPrinterFunc(name, printer string) func(*starlark.Thread, string) {
	switch strings.ToLower(strings.TrimSpace(printer)) {
	case "none", "nil", "no":
		return func(thread *starlark.Thread, msg string) {}
	case "stdout":
		return func(thread *starlark.Thread, msg string) {
			fmt.Println(msg)
		}
	case "stderr":
		return func(thread *starlark.Thread, msg string) {
			fmt.Fprintln(os.Stderr, msg)
		}
	case "basic":
		return nil
	case "lineno", "linenum":
		cnt := atomic.NewInt64(0)
		return func(thread *starlark.Thread, msg string) {
			//prefix := fmt.Sprintf("%04d [⭐|%s](%s)", cnt.Inc(), name, time.Now().UTC().Format(`15:04:05.000`))
			prefix := fmt.Sprintf("[%04d](%s)", cnt.Inc(), time.Now().UTC().Format(`15:04:05.000`))
			fmt.Fprintln(os.Stderr, prefix, msg)
		}
	default:
		return nil
	}
}
