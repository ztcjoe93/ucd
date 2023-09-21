# Unnecessary chdir (UCD)
A wrapper for the common `cd` shell utility, with totally unnecessary features.

## Setup

Append the following to your specific shell [RCfile](https://en.wikipedia.org/wiki/RCFile) to forward stdout from `ucd` to shell's builtin `cd` command.

Example for .zshrc  
```shell
function cd() {
    builtin cd $(ucd $@)
}
```

## Usage

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| -h | - | - | Display help |
| -v | - | - | Display version | 
| -c | bool | false | Clear history and stash list |
| -l | - | - | Display a Most Recently Used (MRU) list of paths cd-ed |
| -ls | - | - | Display a list of stashed cd commands |
| -s | bool | false | stash the cd path to a separately tracked list |
| -p | int | 0 | chdir to the indicated # from MRU list |
| -r | int | 1 | Number of repeats for a given path (for dynamic path i.e ..) |

