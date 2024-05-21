package main

import (
	"os"

	"bitbucket.org/neiku/hlog"
	"bitbucket.org/neiku/winornot"
	"github.com/PureMature/starcli/cli"
	"github.com/PureMature/starcli/module/sys"
	"github.com/PureMature/starcli/web"
	"go.uber.org/zap"
)

func init() {
	// fix for Windows terminal output
	winornot.EnableANSIControl()
}

func main() {
	// parse args
	args := cli.ParseArgs()
	// set log level
	initLogs(args.LogLevel)
	os.Exit(cli.Process(args))
}

func initLogs(level string) {
	lg := hlog.NewSimpleLogger()
	if err := lg.SetLevelString(level); err != nil {
		lg.Error(err)
	}
	log := lg.SugaredLogger.With(zap.Int("pid", os.Getpid()))
	// set log for sub-packages
	web.SetLog(log)
	sys.SetLog(log)
}
