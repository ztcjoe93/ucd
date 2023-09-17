package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	repeatFlag int
)

func main() {
	log.Printf("ucd-v0.1\n")

	homeDir, _ := os.UserHomeDir()
	curDir, _ := os.Getwd()

	log.Printf("~: %v\n", homeDir)
	log.Printf("cwd: %v\n", curDir)

	// flags
	flag.IntVar(&repeatFlag, "r", 1, "repeat dynamic cd path (for ..)")
	flag.Parse()

	args := flag.Args()
	log.Printf("args: %v\n", args)

	if len(args) > 1 {
		log.Fatalln("Only < 1 arguments can be passed to ucd")
	}

	// fmt.Print sends output to stdout, this will be consumed by builtin `cd` command
	if len(args) > 0 {
		fmt.Print(repeat(args[0], repeatFlag))
	} else {
		fmt.Print(homeDir)
	}
}

func repeat(str string, times int) string {
	s := make([]string, times)
	for i := range s {
		s[i] = str
	}

	return strings.Join(s, "/")
}
