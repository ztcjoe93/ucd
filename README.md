# Unnecessary chdir (UCD)
A wrapper for the common `cd` shell utility, with totally unnecessary features.

## Setup

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
| -c | bool | false | clear history and stash list |
| -d | int | 0 | swap directory at -d parent directories |
| -l | - | - | display Most Recently Used (MRU) list of paths chdir-ed into |
| -ls | - | - | display list of stashed cd commands |
| -p | int | 0 | chdir to the indicated # from MRU list |
| -n | int | 1 | no. of times to execute chdir |
| -s | bool | false | stash cd path into a separate list |

