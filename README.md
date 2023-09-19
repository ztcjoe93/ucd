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
| --- | --- | --- |
| -h | - | - | Display help |
| -c | bool | false | Clear history list |
| -l | - | - | Display a Most Recently Used (MRU) list of paths cd-ed |
| -p | int | 0 | chdir to the indicated # from MRU list |
| -r | int | 0 | Repeat given path (solely for dynamic path i.e ..) |

