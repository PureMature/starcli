package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/1set/starlet"
	"github.com/1set/starlet/dataconv"
	"go.starlark.net/starlark"
)

func setMachineExtras(m *starlet.Machine, args []string) {
	sysLoader := loadSysModule(args)
	m.AddPreloadModules(starlet.ModuleLoaderList{sysLoader})
	m.AddLazyloadModules(starlet.ModuleLoaderMap{"sys": sysLoader})
}

func loadSysModule(args []string) func() (starlark.StringDict, error) {
	// get sa
	sa := make([]starlark.Value, 0, len(args))
	for _, arg := range args {
		sa = append(sa, starlark.String(arg))
	}
	// build module
	sd := starlark.StringDict{
		"platform": starlark.String(runtime.GOOS),
		"arch":     starlark.String(runtime.GOARCH),
		"version":  starlark.MakeUint(starlark.CompilerVersion),
		"argv":     starlark.NewList(sa),
		"input":    starlark.NewBuiltin("sys.input", rawStdInput),
	}
	return dataconv.WrapModuleData("sys", sd)
}

func rawStdInput(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	// unpack arguments
	var prompt string
	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "prompt?", &prompt); err != nil {
		return starlark.None, err
	}
	// display prompt
	if prompt != "" {
		fmt.Print(prompt)
	}
	// read input from stdin
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	// trim newline characters
	input = strings.TrimRight(input, "\r\n")
	return starlark.String(input), nil
}
