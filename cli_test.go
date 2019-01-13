package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var cmdOutputPrefix = "test_command"
var cmdErrOutputPrefix = "test_error_command"

func TestRun_argsValidation(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	var cases = []struct {
		cmd  string
		want int
	}{
		{cmd: "altercode -contain warning -exit-code 254 -c testdata/test.toml -- testdata/bin/test_command warning", want: 2},
		{cmd: "altercode -- testdata/bin/test_command warning", want: 2},
		{cmd: "altercode -contain warning -exit-code 254", want: 2},
		{cmd: "altercode -exit-code 254 -- testdata/bin/test_command warning", want: 2},
	}

	for _, c := range cases {
		t.Run(c.cmd, func(t *testing.T) {
			args := strings.Split(c.cmd, " ")
			if got, want := cli.Run(args), c.want; got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	args := strings.Split("altercode -contain warning -exit-code 254 -- testdata/bin/test_command warning", " ")
	if got, want := cli.Run(args), 254; got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	if got, want := outStream.String(), fmt.Sprintf("%s: warning\n", cmdOutputPrefix); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestRun_error(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	args := strings.Split("altercode -contain warning -exit-code 254 -- testdata/bin/test_error_command warning", " ")
	if got, want := cli.Run(args), 1; got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	if got, want := errStream.String(), fmt.Sprintf("%s: warning\n", cmdErrOutputPrefix); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestRun_withConfigFile(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	args := strings.Split("altercode -c testdata/test.toml -- testdata/bin/test_command warning", " ")
	if got, want := cli.Run(args), 254; got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	if got, want := outStream.String(), fmt.Sprintf("%s: warning\n", cmdOutputPrefix); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestRun_withConfigFileError(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	args := strings.Split("altercode -c testdata/test.toml -- testdata/bin/test_error_command warning", " ")
	if got, want := cli.Run(args), 1; got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	if got, want := errStream.String(), fmt.Sprintf("%s: warning\n", cmdErrOutputPrefix); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestRun_withConfigFile_multiRules(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	var cases = []struct {
		cmd  string
		want int
	}{
		{cmd: "altercode -c testdata/test_multi.toml -- testdata/bin/test_command warning", want: 254},
		{cmd: "altercode -c testdata/test_multi.toml -- testdata/bin/test_command deprecated", want: 253},
		{cmd: "altercode -c testdata/test_multi.toml -- testdata/bin/test_command warning_deprecated", want: 254},
		{cmd: "altercode -c testdata/test_multi.toml -- testdata/bin/test_command ok", want: 0},
		{cmd: "altercode -c testdata/test_multi.toml -- testdata/bin/test_error_command ok", want: 1},
	}

	for _, c := range cases {
		t.Run(c.cmd, func(t *testing.T) {
			args := strings.Split(c.cmd, " ")
			if got, want := cli.Run(args), c.want; got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})
	}
}
