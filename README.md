# gowatch

The `gowatch` tool is a simple command-line stopwatch. It supports multiple stopwatches
separated by name.


## Quickstart

1. Install from the github repository:

```bash
go get -u github.com/dusktreader/gowatch
```


## Example usage

```bash
$ gowatch start
```

Some time later...
```bash
$ gowatch stop
1m59.74s
```


## Getting help

Simply run `gowatch --help`:

```
$ gowatch --help
A command line stopwatch written in Go

Usage:
  gowatch [flags]
  gowatch [command]

Available Commands:
  clear       Clear timers
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        List all timers
  reset       Reset a timer
  show        Show a timer
  start       Start a timer
  stop        Stop a timer
  toggle      Toggle a timer

Flags:
  -h, --help      help for gowatch
  -v, --verbose   Show verbose logging output

Use "gowatch [command] --help" for more information about a command.
```


## License

Distributed under the MIT License. See `LICENSE` for more information.
