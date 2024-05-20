package cli

import (
	"errors"
	"fmt"

	"github.com/1set/gut/ystring"
	"github.com/1set/starbox"
	"github.com/PureMature/starcli/box"
	"github.com/PureMature/starcli/module/sys"
	"github.com/PureMature/starcli/web"
	flag "github.com/spf13/pflag"
)

// runWebServer starts a web server that creates a Starbox with given code for each request.
func runWebServer(args *Args) error {
	var (
		runner        = starbox.NewRunConfig()
		webPort       = args.WebPort
		numArg        = args.NumberOfArgs
		useDirectCode = ystring.IsNotBlank(args.CodeContent)
	)

	// prepare runner
	if useDirectCode {
		// if code content is provided in flag, just use it
		runner = runner.FileName("web.star").Script(args.CodeContent)
	} else if numArg >= 1 {
		// or use the first argument as file name
		runner = runner.FileName(args.Arguments[0])
	} else {
		// no repl mode for web server, just quit if no code if provided
		return errors.New("no code to run as web server")
	}

	// start web server
	build := func() *starbox.RunnerConfig {
		b := box.Build("web", args.IncludePath, args.LoadModules)
		b.AddModuleLoader(sys.ModuleName, sys.NewModule(args.Arguments))
		return runner.Starbox(b)
	}
	return web.Start(webPort, build)
}

func runDirectCode(args *Args) error {
	fmt.Println("runDirectCode", args)
	return nil
}

func runREPL(args *Args) error {
	fmt.Println("runREPL", args)
	return nil
}

func runScriptFile(args *Args) error {
	fmt.Println("runScriptFile", args)
	return nil
}

func showHelp(args *Args) error {
	flag.Usage()
	return nil
}
