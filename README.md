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

