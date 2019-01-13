package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"

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

	if (contain != "" && config != "") || (contain == "" && config == "") {
		fmt.Fprintln(c.errStream, "invalid err")
		return ExitCodeError
	}

	var r rules
	if contain != "" {
		r.Rule = []rule{{}}
		r.Rule[0].CondType = typeContain
		r.Rule[0].Condition = contain
		r.Rule[0].ExitCode = exitCode
	}

	if config != "" {
		if _, err := toml.DecodeFile(config, &r); err != nil {
			fmt.Fprintf(c.errStream, "decode %s failed: %s\n", config, err.Error())
			return ExitCodeError
		}
	}

	if len(flags.Args()) < 1 {
		fmt.Fprintln(c.errStream, "missing argument")
		return ExitCodeError
	}

	cmdName, cmdArgs := flags.Args()[0], flags.Args()[1:]
	cmd := exec.CommandContext(context.TODO(), cmdName, cmdArgs...)

	out, err := cmd.Output()
	if err != nil {
		if execErr, ok := err.(*exec.Error); ok {
			if execErr.Err == exec.ErrNotFound {
				fmt.Fprintln(c.errStream, err.Error())
				return ExitCodeError
			}
		}

		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Fprintf(c.errStream, string(exitErr.Stderr))
			if ws, ok := exitErr.ProcessState.Sys().(syscall.WaitStatus); ok {
				return ws.ExitStatus()
			}
			return ExitCodeError
		}
	}

	fmt.Fprint(c.outStream, string(out))

	for _, rr := range r.Rule {
		if strings.Contains(string(out), rr.Condition) {
			return rr.ExitCode
		}
	}

	return ExitCodeOK
}
