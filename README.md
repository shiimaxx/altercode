altercode
=========

Wrap command and can alter exit code.

## Description

A command line tool that handles exit code for command specified by an argument. It’s useful in case you want to force failed for CI and so on.


## Usage

Execute `command` and exit with 1 status code when command output contains “Error”.

```
altercode -contain "Error" -exit-code 1 -- command
```


## License

[MIT](https://github.com/shiimaxx/altercode/blob/master/LICENCE)


## Author

[shiimaxx](https://github.com/shiimaxx)

