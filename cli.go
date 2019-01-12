package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

var typeContain = "contain"

type rule struct {
	CondType  string `toml:"type"`
	Condition string
	ExitCode  int `toml:"exit_code"`
}

type rules struct {
	Rule []rule
}

// Run invokes the CLI with the given arguments.
func (c *CLI) Run(args []string) int {
	var (
		contain  string
		exitCode int
		config   string
		version  bool
	)

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(c.outStream)
	flags.StringVar(&contain, "contain", "", "")
	flags.IntVar(&exitCode, "exit-code", ExitCodeError, "")
	flags.StringVar(&config, "c", "", "")
	flags.BoolVar(&version, "version", false, "")

	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprint(c.errStream, "flag parse failed: ", err.Error())
		return ExitCodeError
	}

	if version {
		fmt.Fprintf(c.outStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	if contain != "" && config != "" {
		fmt.Fprintln(c.errStream, "invalid err")
		return ExitCodeError
	}

	if config != "" {
		var r rules
		if _, err := toml.DecodeFile(config, &r); err != nil {
			fmt.Fprintf(c.errStream, "decode %s failed: %s\n", config, err.Error())
			return ExitCodeError
		}
		contain = r.Rule[0].Condition
		exitCode = r.Rule[0].ExitCode
	}

	if len(flags.Args()) < 1 {
		fmt.Fprintln(c.errStream, "missing argument")
		return ExitCodeError
	}

	cmdName, cmdArgs := flags.Args()[0], flags.Args()[1:]
	cmd := exec.CommandContext(context.TODO(), cmdName, cmdArgs...)

	out, _ := cmd.Output()

	fmt.Fprint(c.outStream, string(out))

	if strings.Contains(string(out), contain) {
		return exitCode
	}

	return ExitCodeOK
}
