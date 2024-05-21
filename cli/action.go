package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/1set/gut/ystring"
	"github.com/1set/starbox"
	"github.com/PureMature/starcli/config"
	"github.com/PureMature/starcli/web"
	flag "github.com/spf13/pflag"
	"golang.org/x/term"
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
		b := BuildBox("web",
			args.IncludePath,
			args.ModulesToLoad,
			args.Arguments,
		)
		return runner.Starbox(b)
	}
	return web.Start(webPort, build)
}

func runDirectCode(args *Args) error {
	box := BuildBox("direct",
		args.IncludePath,
		args.ModulesToLoad,
		append([]string{`-c`}, args.Arguments...),
	)
	_, err := box.Run(args.CodeContent)
	return err
}

func runREPL(args *Args) error {
	stdinIsTerminal := term.IsTerminal(int(os.Stdin.Fd()))
	if stdinIsTerminal {
		config.DisplayBuildInfo()
	}
	box := BuildBox("repl",
		args.IncludePath,
		args.ModulesToLoad,
		[]string{``},
	)
	err := box.REPL()
	if stdinIsTerminal {
		fmt.Println()
	}
	return err
}

func runScriptFile(args *Args) error {
	// load file
	fileName := args.Arguments[0]
	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	// build and run
	name := filepath.Base(fileName)
	box := BuildBox(name,
		args.IncludePath,
		args.ModulesToLoad,
		args.Arguments,
	)
	_, err = box.CreateRunConfig().
		FileName(name).
		Script(string(bs)).
		Execute()
	return err
}

func showVersion(args *Args) error {
	config.DisplayBuildInfo()
	return nil
}

func showHelp(args *Args) error {
	flag.Usage()
	return nil
}
