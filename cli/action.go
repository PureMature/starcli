package cli

import (
	"errors"
	"fmt"

	"github.com/1set/gut/ystring"
	"github.com/1set/starbox"
	"github.com/PureMature/starcli/box"
	flag "github.com/spf13/pflag"
)

func runWebServer(args *Args) error {
	var (
		runner        = starbox.NewRunConfig()
		numArg        = args.NumberOfArgs
		useDirectCode = ystring.IsNotBlank(args.CodeContent)
	)

	// prepare runner
	if useDirectCode {
		runner = runner.FileName("web.star").Script(args.CodeContent)
	} else if numArg >= 1 {
		runner = runner.FileName(args.Arguments[0])
	} else {
		return errors.New("no code to run as web server")
	}

	// start web server
	gotRunner := func() *starbox.RunnerConfig {
		return runner.Starbox(box.Build("web", args.IncludePath, args.LoadModules))
	}

	// HACK: test it
	r := gotRunner()
	res, err := r.Execute()
	if err != nil {
		return err
	}
	fmt.Println(args.WebPort, res)

	return nil
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
