package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	args := strings.Split("altercode -contain 'warning' -exit-code 254 -- testdata/bin/test_command warning", " ")
	if got, want := cli.Run(args), 254; got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestRun_withConfigFile(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	args := strings.Split("altercode -c testdata/test.toml -- testdata/bin/test_command warning", " ")
	if got, want := cli.Run(args), 254; got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
