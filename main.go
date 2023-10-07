package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ztcjoe93/ucd/configurations"
	"github.com/ztcjoe93/ucd/records"
)

var (
	configs         configurations.Configuration
	aliasFlag       string
	helpFlag        bool
	clearFlag       bool
	clearLimitFlag  bool
	clearStashFlag  bool
	dynamicSwapFlag int
	numRepeatFlag   int
	listFlag        bool
	listStashFlag   bool
	historyPathFlag int
	aliasPathFlag   string
	modifyAliasFlag int
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
	flag.BoolVar(&clearLimitFlag, "cl", false, "clear history list till MaxMRUDisplay limit")
	flag.BoolVar(&clearStashFlag, "cs", false, "clear stash list")
	flag.IntVar(&dynamicSwapFlag, "d", 0, "swap out directory to arg after -d parent directories")
	flag.BoolVar(&listFlag, "l", false, "display Most Recently Used (MRU) list of paths chdir-ed into")
	flag.BoolVar(&listStashFlag, "ls", false, "display list of stashed cd commands")
	flag.IntVar(&modifyAliasFlag, "ma", 0, "modify alias of indicated # from the stash list")
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
		returnCwd()
	}

	if versionFlag {
		log.Printf("%v\n", version)
		returnCwd()
	}

	configs = configs.GetConfigurations()

	cachePath = homeDir + "/.ucd-cache"
	cacheFile, _ := os.Open(cachePath)
	defer cacheFile.Close()
	byteValue, _ := io.ReadAll(cacheFile)

	var r records.Records
	err := json.Unmarshal(byteValue, &r)
	if err != nil {
		r = records.Records{
			PathRecords:  map[string]records.PathRecord{},
			StashRecords: map[string]records.StashRecord{},
		}
	}

	if configs.MaxMRUDisplay < 0 {
		configs.MaxMRUDisplay = len(r.PathRecords)
	}

	if clearFlag {
		r = records.Records{
			PathRecords:  map[string]records.PathRecord{},
			StashRecords: r.StashRecords,
		}
		output, _ := json.Marshal(r)
		os.WriteFile(cachePath, output, 0644)
		returnCwd()
	}

	if clearLimitFlag {
		rk := records.SortRecords(r.PathRecords)

		if len(rk) > configs.MaxMRUDisplay {
			for i := configs.MaxMRUDisplay; i < len(rk); i++ {
				delete(r.PathRecords, rk[i])
			}

		}
		output, _ := json.Marshal(r)
		os.WriteFile(cachePath, output, 0644)
		returnCwd()
	}

	if clearStashFlag {
		r = records.Records{
			PathRecords:  r.PathRecords,
			StashRecords: map[string]records.StashRecord{},
		}
		output, _ := json.Marshal(r)
		os.WriteFile(cachePath, output, 0644)
		returnCwd()
	}

	// exit earlier depending on flag passed in
	if listFlag {
		r.ListRecords("path", configs.MaxMRUDisplay)
		returnCwd()
	}

	if listStashFlag {
		r.ListRecords("stash", configs.MaxMRUDisplay)
		returnCwd()
	}

	if len(args) > 1 {
		log.Fatalln("Only < 1 arguments can be passed to ucd")
		returnCwd()
	}

	if modifyAliasFlag > 0 {
		srk := records.SortRecords(r.StashRecords)
		sr := r.StashRecords[srk[modifyAliasFlag-1]]
		sr.Alias = args[0]
		r.StashRecords[srk[modifyAliasFlag-1]] = sr

		output, _ := json.Marshal(r)
		os.WriteFile(cachePath, output, 0644)

		r.ListRecords("stash", configs.MaxMRUDisplay)
		returnCwd()
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
			returnCwd()
		}
	} else if historyPathFlag > 0 {
		mruRecords := records.SortRecords(r.PathRecords)
		if historyPathFlag-1 > len(mruRecords)-1 {
			log.Printf("invalid #, there are only %v records\n", len(mruRecords))
			returnCwd()
		}
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
		returnCwd()
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

	if stashFlag {
		if r.AliasExists(aliasFlag) {
			log.Printf("Alias `%v` already exists\n", aliasFlag)
			returnCwd()
		}
		r.StashRecords[targetPath] = records.StashRecord{Alias: aliasFlag, Timestamp: timeNow()}
	}
	fmt.Print(targetPath)

	output, _ := json.Marshal(r)
	os.WriteFile(cachePath, output, 0644)
}

func returnCwd() {
	fmt.Print(".")
	os.Exit(1)
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
		log.Printf("path `%v` is not a valid path\n", targetPath)
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
