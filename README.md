# Unnecessary chdir (UCD)
A wrapper for the common `cd` shell utility, with totally unnecessary features.

## Setup

### Compiling ucd binary

If you have golang installed in your environment, you can build the golang binary and shift it into your usr/bin directory.  
```shell
go build . && sudo chmod +x ucd && sudo mv ucd/usr/local/bin/ucd
```

Otherwise, you can download the binary and shift it into your usr/bin directory.


### Redirecting stdout to builtin shell

Append the following to your specific shell [RCfile](https://en.wikipedia.org/wiki/RCFile) to forward stdout from `ucd` to shell's builtin `cd` command.

Example for .zshrc  
```shell
function cd() { builtin cd $(ucd $@) }
```

## Usage

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| -h | - | - | display help |
| -v | - | - | display ucd version | 
| -c | bool | false | clear history list |
| -cs | bool | false | clear stash list |
| -d | int | 0 | swap directory at -d parent directories |
| -l | - | - | display Most Recently Used (MRU) list of paths chdir-ed into |
| -ls | - | - | display list of stashed cd commands |
| -p | int | 0 | chdir to the indicated # from MRU list |
| -ps | int | 0 | chdir to the indicated # from stash list |
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
zt@ragnarok-arch directory$ pwd
# /home/zt/my/path/ci/to/a/particular/directory
zt@ragnarok-arch directory$ cd -d 4 uat
zt@ragnarok-arch directory$ pwd
# /home/zt/my/path/uat/to/a/particular/directory
```

### -p / -ps usage

Does a `chdir` into the indicated # path from either the history/stash list.  

```shell
zt@ragnarok-arch zt$ cd -ls
+---+-------------------------+-------------------------+
| # | PATH                    | TIMESTAMP               |
+---+-------------------------+-------------------------+
| 1 | /home/zt/.config/waybar | 2023-09-23 12:53:39 +08 |
| 2 | /home/zt/.config/hypr   | 2023-09-23 12:53:02 +08 |
| 3 | /home/zt                | 2023-09-23 12:52:51 +08 |
+---+-------------------------+-------------------------+
zt@ragnarok-arch zt$ pwd
/home/zt
zt@ragnarok-arch zt$ cd -ps 1
zt@ragnarok-arch waybar$ pwd
/home/zt/.config/waybar
```

