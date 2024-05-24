# starcli

Windows

Just download:

https://github.com/PureMature/starcli/releases/download/b2/starcli.exe


macOS

curl -sSL https://github.com/PureMature/starcli/releases/download/b2/starcli.macos -o starcli
install -m 0755 starcli /usr/local/bin

Linux:

curl -sSL https://github.com/PureMature/starcli/releases/download/b2/starcli.linux -o starcli
install -m 0755 starcli /usr/local/bin


Usage:

```bash
$ ./starcli -h
Usage of ./starcli:
  -c, --code string      Starlark code to execute
  -C, --config string    config file to load
  -g, --globalreassign   allow reassigning global variables in Starlark code
  -I, --include string   include path for Starlark code to load modules from (default ".")
  -i, --interactive      enter interactive mode after executing
  -l, --log string       log level: debug, info, warn, error, dpanic, panic, fatal (default "info")
  -m, --module strings   Modules to load before executing Starlark code (default [atom,base64,csv,file,go_idiomatic,hashlib,http,json,log,math,path,random,re,runtime,string,struct,sys,time])
  -o, --output string    output printer: none,stdout,stderr,basic,lineno,auto (default "auto")
  -r, --recursion        allow recursion in Starlark code
  -V, --version          print version & build information
  -w, --web uint16       run web server on specified port, it provides request and response structs for Starlark code to use
```
