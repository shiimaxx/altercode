altercode
=========

Wrap command and can alter exit code.

## Description

A command line tool that handles exit code for command specified by an argument. It’s useful in case you want to force failed for CI and so on.


## Usage

Execute `command` and exit with 1 status code when command output contains “Error”. This is most simple case.

```
altercode -contain "Error" -exit-code 1 -- command
```

When multiple rules, you can write rules in a configuration file and run a command with specifying that file in `-c` option.
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

## License

[MIT](https://github.com/shiimaxx/altercode/blob/master/LICENCE)


## Author

[shiimaxx](https://github.com/shiimaxx)

