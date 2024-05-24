# :stars: StarCLI

StarCLI is a command-line interface for executing Starlark scripts with various options for configuration, logging, and module inclusion.

## Installation

### Windows

Download the executable from the [releases page](https://github.com/PureMature/starcli/releases/download/b2/starcli.exe).

### macOS

Run the following commands in your terminal:

```sh
curl -sSL https://github.com/PureMature/starcli/releases/download/b2/starcli.macos -o starcli
install -m 0755 starcli /usr/local/bin
```

### Linux

Run the following commands in your terminal:

```sh
curl -sSL https://github.com/PureMature/starcli/releases/download/b2/starcli.linux -o starcli
install -m 0755 starcli /usr/local/bin
```

## Usage

Run `starcli` with the `-h` option to see a list of available commands and options.

```bash
$ ./starcli -h
Usage of ./starcli:
  -c, --code string      Starlark code to execute
  -C, --config string    Config file to load
  -g, --globalreassign   Allow reassigning global variables in Starlark code
  -I, --include string   Include path for Starlark code to load modules from (default ".")
  -i, --interactive      Enter interactive mode after executing
  -l, --log string       Log level: debug, info, warn, error, dpanic, panic, fatal (default "info")
  -m, --module strings   Modules to load before executing Starlark code (default [atom,base64,csv,file,go_idiomatic,hashlib,http,json,log,math,path,random,re,runtime,string,struct,sys,time])
  -o, --output string    Output printer: none, stdout, stderr, basic, lineno, auto (default "auto")
  -r, --recursion        Allow recursion in Starlark code
  -V, --version          Print version & build information
  -w, --web uint16       Run web server on specified port, it provides request and response structs for Starlark code to use
```

### Examples

- **Start REPL Mode:** `starcli`
- **Execute Starlark Code:** `starcli -c 'print("Hello, World!")'`
- **Enter Interactive Mode After Execution:** `starcli -c 's = "Hello, World!"' -i`
- **Execute Starlark Code with Log Level:** `starcli --log debug test.star`

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any changes.

## Contact

For any questions or support, please open an issue on [GitHub](https://github.com/PureMature/starcli/issues).
