[![go-tests](https://github.com/ztcjoe93/ucd/actions/workflows/test.yml/badge.svg?branch=main&event=workflow_dispatch)](https://github.com/ztcjoe93/ucd/actions/workflows/test.yml)

# Unnecessary chdir (UCD)
A wrapper for the common `cd` shell utility, with totally unnecessary features.

## Setup

### Compiling ucd binary

If you have golang installed in your environment, you can build the golang binary and shift it into your usr/bin directory.  
```shell
go build . && sudo chmod +x ucd && sudo mv ucd /usr/local/bin/ucd
```

Otherwise, you can download the binary and shift it into your usr/bin directory.


### Redirecting stdout to builtin shell

Append the following to your specific shell [runcom](https://en.wikipedia.org/wiki/RUNCOM) file  to forward stdout from `ucd` to shell's builtin `cd` command.

Example for .zshrc  
```shell
function cd() { builtin cd $(ucd $@) }
```

## Usage

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| -h | - | - | display help |
| -v | - | - | display ucd version | 
| -a | string |  | alias for stashed path, used in conjunction with -s |
| -c | bool | false | clear history list |
| -cs | bool | false | clear stash list |
| -d | int | 0 | swap directory at -d parent directories |
| -l | - | - | display Most Recently Used (MRU) list of paths chdir-ed into |
| -ls | - | - | display list of stashed cd commands |
| -ma | int | 0 | modify alias of indicated # from the stash list |
| -p | int | 0 | chdir to the indicated # from MRU list |
| -ps | int | 0 | chdir to the indicated # from stash list |
| -pa | string |  | chdir to path with matching alias from stash list |
| -n | int | 1 | no. of times to execute chdir |
| -s | bool | false | stash cd path into a separate list |


### -d usage

Dynamic swapping of a sub-directory path when sub-directory trees are similar.  

For example, you have the following directory tree:
```shell
my
└── path
    ├── ci
    │   └── to
    │       └── a
    │           └── particular
    │               └── directory
    └── uat
        └── to
            └── a
                └── particular
                    └── directory
```

To shift from `/my/path/ci/to/a/particular/directory` to `my/path/uat/to/a/particular/directory`, swap the directory to the argument after traversing to the parent directory `4` times.  

```shell
$ pwd
# /home/zt/my/path/ci/to/a/particular/directory
$ cd -d 4 uat
$ pwd
# /home/zt/my/path/uat/to/a/particular/directory
```

### -s / -a usage

Stashes the cd-ed path when `-s` is provided. An alias can be provided with the `-a` parameter.

```shell
$ pwd
# /
$ cd -s -a usr-bin usr/bin
$ cd -ls
+---+---------+----------+-------------------------+
| # | ALIAS   | PATH     | TIMESTAMP               |
+---+---------+----------+-------------------------+
| 1 | usr-bin | /usr/bin | 2023-09-24 22:27:25 +08 |
+---+---------+----------+-------------------------+
```

### -ma usage

Modifies the alias for the indicated # path from the stash list.  

```shell
$ cd -ls
+---+-------+------------------+-------------------------+
| # | ALIAS | PATH             | TIMESTAMP               |
+---+-------+------------------+-------------------------+
| 1 | ucd   | /home/zt/dev/ucd | 2023-09-24 22:28:11 +08 |
| 2 | mybin | /usr/bin         | 2023-09-24 22:27:25 +08 |
+---+-------+------------------+-------------------------+
$ cd -ma 2 usr/bin
+---+---------+------------------+-------------------------+
| # | ALIAS   | PATH             | TIMESTAMP               |
+---+---------+------------------+-------------------------+
| 1 | ucd     | /home/zt/dev/ucd | 2023-09-24 22:28:11 +08 |
| 2 | usr/bin | /usr/bin         | 2023-09-24 22:27:25 +08 |
+---+---------+------------------+-------------------------+
$ cd -ls
+---+---------+------------------+-------------------------+
| # | ALIAS   | PATH             | TIMESTAMP               |
+---+---------+------------------+-------------------------+
| 1 | ucd     | /home/zt/dev/ucd | 2023-09-24 22:28:11 +08 |
| 2 | usr/bin | /usr/bin         | 2023-09-24 22:27:25 +08 |
+---+---------+------------------+-------------------------+
```

### -p / -ps usage

Does a `chdir` into the indicated # path from either the history/stash list.  

```shell
$ cd -ls
+---+-------------------------+-------------------------+
| # | PATH                    | TIMESTAMP               |
+---+-------------------------+-------------------------+
| 1 | /home/zt/.config/waybar | 2023-09-23 12:53:39 +08 |
| 2 | /home/zt/.config/hypr   | 2023-09-23 12:53:02 +08 |
| 3 | /home/zt                | 2023-09-23 12:52:51 +08 |
+---+-------------------------+-------------------------+
$ pwd
# /home/zt
$ cd -ps 1
$ pwd
# /home/zt/.config/waybar
```

### -pa usage

Does a `chdir` into the matching path from the stash list if the provided alias exist.

```shell
$ cd -ls
+---+-----------+-------------------------+-------------------------+
| # | ALIAS     | PATH                    | TIMESTAMP               |
+---+-----------+-------------------------+-------------------------+
| 1 | ucd       | /home/zt/dev/ucd        | 2023-09-24 22:02:18 +08 |
| 2 | hyperland | /home/zt/.config/hypr   | 2023-09-24 21:56:12 +08 |
| 3 |           | /home/zt/.config/waybar | 2023-09-23 12:53:39 +08 |
| 4 |           | /home/zt                | 2023-09-23 12:52:51 +08 |
+---+-----------+-------------------------+-------------------------+
$ cd -pa hyperland
$ pwd
# /home/zt/.config/hypr
$ cd -pa ucd
$ pwd
# /home/zt/dev/ucd
```

### -n usage

Repeats the `chdir` command a number of `-n` times. Solely for parent directory jumping.  

```shell
$ pwd
# /home/zt/my/path/uat/to/a/particular/directory
$ cd -n 3 ..
$ pwd
# /home/zt/my/path/uat/to
```

## Configuration

On `ucd`'s first run, a `ucd.conf` JSON file is generated at `$HOME/.config/ucd`, where k-v pairs can be passed in to tweak certain features.  

| Parameter | Type | Default | Description |
| --- | --- | --- | --- |
| MaxMRUDisplay | int | -1 | Limits the total number of paths displayed when using `-l` or `-ls`. Set this to `-1` to show all records. |
| FileFallbackBehavior | bool | true | Toggle to set if the default behavior of cd-ing to a file is to fallback to its parent directory |

## Testing

To run all test suites, run the following in the repo root directory: 
```shell
go test -v ./...
```

## License

`ucd` is released under the [MIT](LICENSE.md) license.
