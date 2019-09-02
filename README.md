# tfpolicy
tfpolicy is terraform hcl2 file linter that enable to policy as as code 

## Status

**🚧Work In Progress🚧**

## Installation

Install or upgrade tfpolicy with this command.

```bash
go get -u github.com/laqiiz/tfpolicy
```

## Usage

```console
$ tfpolicy --help
Usage of tfpolicy:
  -d string
        d is target directory path (default ".")
  -dir string
        dir is target directory path (default ".")
```

### For Example

```
$ git clone https://github.com/laqiiz/tfpolicy
$ cd tfpolicy
$ tfpolicy examples
2019/04/28 22:57:32 2 errors occurred:
	* [ERROR] \testdata\example.tf: resource label must match pattern [a-zA-Z]+-[a-zA-Z]+: label='sql',type='mst_db'
	* [ERROR] \testdata\example.tf: resource label must match pattern [a-zA-Z]+-[a-zA-Z]+: label='sql',type='ope_db'
```

## License

This project is licensed under the Apache License 2.0 License - see the [LICENSE](LICENSE) file for details
