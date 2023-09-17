package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	enabledFlag bool
)

func main() {
	log.Printf("ucd-v0.1\n")

	homeDir, _ := os.UserHomeDir()
	log.Printf("Home directory = %v\n", homeDir)

	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("Current working directory = %v\n", curDir)

	args := os.Args
	log.Printf("Arguments: %v\n", args)

	// fmt.Print sends output to stdout, this will be consumed by builtin `cd` command
	if len(args) == 1 {
		fmt.Print(homeDir)
	} else {
		fmt.Print(args[1])
	}

	flag.Parse()
}
