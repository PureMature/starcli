package cli

import (
	"github.com/1set/gut/ystring"
	"github.com/1set/starlet"
	"github.com/PureMature/starcli/util"
	flag "github.com/spf13/pflag"
)

// Args is the command line arguments for StarCLI.
type Args struct {
	AllowRecursion      bool
	AllowGlobalReassign bool
	LoadModules         []string
	IncludePath         string
	FileName            string
	CodeContent         string
	WebPort             uint16
	NumberOfArgs        int
	Arguments           []string
}

var (
	defaultModules = starlet.GetAllBuiltinModuleNames()
)

// ParseArgs parses command line arguments and returns the Args object.
func ParseArgs() *Args {
	args := &Args{}

	// parse command line arguments
	flag.BoolVarP(&args.AllowRecursion, "recursion", "r", false, "allow recursion in Starlark code")
	flag.BoolVarP(&args.AllowGlobalReassign, "globalreassign", "g", false, "allow reassigning global variables in Starlark code")
	flag.StringSliceVarP(&args.LoadModules, "module", "m", defaultModules, "Modules to load before executing Starlark code")
	flag.StringVarP(&args.IncludePath, "include", "i", ".", "include path for Starlark code to load modules from")
	flag.StringVarP(&args.CodeContent, "code", "c", "", "Starlark code to execute")
	flag.Uint16VarP(&args.WebPort, "web", "w", 0, "run web server on specified port, it provides request and response structs for Starlark code to use")
	flag.Parse()

	// keep the rest of arguments
	args.NumberOfArgs = flag.NArg()
	args.Arguments = flag.Args()
	return args
}

// Process processes the command line arguments and executes desired actions, it returns the exit code.
func Process(args *Args) int {
	// for basic checks
	numArg := args.NumberOfArgs
	useDirectCode := ystring.IsNotBlank(args.CodeContent)

	// determine action
	var action func(*Args) error
	switch {
	case args.WebPort > 0:
		action = runWebServer
	case useDirectCode:
		action = runDirectCode
	case numArg == 0:
		action = runREPL
	case numArg >= 1:
		action = runScriptFile
	default:
		action = showHelp
	}

	// execute action
	err := action(args)
	if err != nil {
		util.PrintError(err)
		return 1
	}
	return 0
}