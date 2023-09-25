package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ztcjoe93/ucd/records"
	"github.com/ztcjoe93/ucd/configurations"
)

var (
	aliasFlag       string
	helpFlag        bool
	clearFlag       bool
	clearStashFlag  bool
	dynamicSwapFlag int
	numRepeatFlag   int
	listFlag        bool
	listStashFlag   bool
	historyPathFlag int
	aliasPathFlag   string
	stashPathFlag   int
	stashFlag       bool
	versionFlag     bool
	cachePath       string
	cacheFile       *os.File

	invalidPath bool   = false
	version     string = "ucd v0.1.0"
)

func main() {
	log.SetFlags(0)
	// flags
	flag.BoolVar(&helpFlag, "h", false, "display help")
	flag.StringVar(&aliasFlag, "a", "", "alias for stashed path, used in conjunction with -s")
	flag.BoolVar(&versionFlag, "v", false, "display ucd version")
	flag.BoolVar(&clearFlag, "c", false, "clear history list")
	flag.BoolVar(&clearStashFlag, "cs", false, "clear stash list")
	flag.IntVar(&dynamicSwapFlag, "d", 0, "swap out directory to arg after -d parent directories")
	flag.BoolVar(&listFlag, "l", false, "display Most Recently Used (MRU) list of paths chdir-ed into")
	flag.BoolVar(&listStashFlag, "ls", false, "display list of stashed cd commands")
	flag.IntVar(&numRepeatFlag, "n", 1, "no. of times to execute chdir")
	flag.IntVar(&historyPathFlag, "p", 0, "chdir to the indicated # from the MRU list")
	flag.IntVar(&stashPathFlag, "ps", 0, "chdir to the indicated # from the stash list")
	flag.StringVar(&aliasPathFlag, "pa", "", "chdir to path with matching alias from stash list")
	flag.BoolVar(&stashFlag, "s", false, "stash cd path into a separate list")
	flag.Parse()

	args := flag.Args()
	homeDir, _ := os.UserHomeDir()

	if helpFlag {
		flag.PrintDefaults()
		fmt.Print(".")
		os.Exit(1)
	}

	if versionFlag {
		log.Printf("%v\n", version)
		fmt.Print(".")
		os.Exit(1)
	}

	var c configurations.Configuration
	c = c.GetConfigurations()
	//fmt.Println(c.MaxMRU)

	cachePath = homeDir + "/.ucd-cache"
	cacheFile, _ := os.Open(cachePath)
	byteValue, _ := ioutil.ReadAll(cacheFile)

	var r records.Records
	err := json.Unmarshal(byteValue, &r)
	if err != nil {
		r = records.Records{
			PathRecords:  map[string]records.PathRecord{},
			StashRecords: map[string]records.StashRecord{},
		}
	}

	if clearFlag {
		r = records.Records{
			PathRecords:  map[string]records.PathRecord{},
			StashRecords: r.StashRecords,
		}
		output, _ := json.Marshal(r)
		ioutil.WriteFile(cachePath, output, 0644)
		fmt.Print(".")
		os.Exit(1)
	}

	if clearStashFlag {
		r = records.Records{
			PathRecords:  r.PathRecords,
			StashRecords: map[string]records.StashRecord{},
		}
		output, _ := json.Marshal(r)
		ioutil.WriteFile(cachePath, output, 0644)
		fmt.Print(".")
		os.Exit(1)
	}

	// exit earlier depending on flag passed in
	if listFlag {
		r.ListRecords("path")
		fmt.Print(".")

		os.Exit(1)
	}

	if listStashFlag {
		r.ListRecords("stash")
		fmt.Print(".")

		os.Exit(1)
	}

	if len(args) > 1 {
		log.Fatalln("Only < 1 arguments can be passed to ucd")
	}

	// fmt.Print sends output to stdout, this will be consumed by builtin `cd` command

	var targetPath string
	if dynamicSwapFlag > 0 {
		targetPath = dynamicPathSwap(args[0], dynamicSwapFlag)
	} else if aliasPathFlag != "" {
		found := false
		for key, rec := range r.StashRecords {
			if rec.Alias == aliasPathFlag {
				targetPath = key
				found = true
				break
			}
		}

		if !found {
			log.Printf("unable to cd -- alias ``%v` not found\n", aliasPathFlag)
			fmt.Print(".")
			os.Exit(1)
		}
	} else if historyPathFlag > 0 {
		mruRecords := records.SortRecords(r.PathRecords)
		targetPath = mruRecords[historyPathFlag-1]
	} else if stashPathFlag > 0 {
		stashRecords := records.SortRecords(r.StashRecords)
		targetPath = stashRecords[stashPathFlag-1]
	} else {
		if len(args) > 0 {
			targetPath = repeat(args[0], numRepeatFlag)
		} else {
			targetPath = homeDir
		}
	}

	if isInvalidPath(targetPath) {
		fmt.Print(targetPath)
		os.Exit(1)
	}

	targetPath, _ = os.Getwd()

	rec, ok := r.PathRecords[targetPath]
	if ok {
		rec.Count++
		rec.Timestamp = timeNow()
		r.PathRecords[targetPath] = rec
	} else {
		r.PathRecords[targetPath] = records.PathRecord{Count: 1, Timestamp: timeNow()}
	}

	fmt.Print(targetPath)
	if stashFlag {
		r.StashRecords[targetPath] = records.StashRecord{Alias: aliasFlag, Timestamp: timeNow()}
	}

	output, _ := json.Marshal(r)
	ioutil.WriteFile(cachePath, output, 0644)
	cacheFile.Close()
}

func dynamicPathSwap(swapArg string, upCount int) string {
	paths := make([]string, 0)
	for i := 0; i < upCount; i++ {
		wd, _ := os.Getwd()
		wdArr := strings.Split(wd, "/")
		paths = prependStrSlice(paths, wdArr[len(wdArr)-1])
		os.Chdir("..")
	}

	os.Chdir("..")
	paths = prependStrSlice(paths, swapArg)
	targetPath := strings.Join(paths, "/")

	return targetPath
}

func prependStrSlice(x []string, y string) []string {
	x = append(x, "")
	copy(x[1:], x)
	x[0] = y
	return x
}

func isInvalidPath(targetPath string) bool {
	err := os.Chdir(targetPath)
	if err != nil {
		return true
	}

	return false
}

func repeat(str string, times int) string {
	s := make([]string, times)
	for i := range s {
		s[i] = str
	}

	return strings.Join(s, "/")
}

func timeNow() string {
	return time.Now().Format("2006-01-02 15:04:05 MST")
}

