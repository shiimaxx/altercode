altercode
=========

[![Build Status](https://travis-ci.org/shiimaxx/altercode.svg?branch=master)](https://travis-ci.org/shiimaxx/altercode)
[![Go Report Card](https://goreportcard.com/badge/github.com/shiimaxx/altercode)](https://goreportcard.com/report/github.com/shiimaxx/altercode)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Wrap command and can alter exit code.

## Description

A command line tool that handles exit code for command specified by an argument. It’s useful in case you want to force failed for CI and so on.


## Usage

Execute `command` and exit with 1 status code when command output contains “Error” otherwise, exit with an original status code.
This is most simple case.

```
altercode -contain "Error" -exit-code 1 -- command
```

you can write rules in a configuration file and run a command with specifying that file in `-c` option.
In the following case, created a configuration file named "altercode.toml". This result will same to above example.

```
altercode -c altercode.toml -- command
```

```toml
[[rule]]
type = "contain"
condition = "Error"
exit_code = 1
```

When multiple rules, should define multiple rules in a configuration file. Command line options do not support multiple rules.
Rules priority is same to definition order.   

```toml
[[rule]]
type = "contain"
condition = "Error"
exit_code = 1

[[rule]]
type = "contain"
condition = "Warning"
exit_code = 1
```

## License

[MIT](https://github.com/shiimaxx/altercode/blob/master/LICENCE)


## Author

[shiimaxx](https://github.com/shiimaxx)
