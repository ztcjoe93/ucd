package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/table"
)

var (
	clearFlag       bool
	repeatFlag      int
	listFlag        bool
	listStashFlag   bool
	historyPathFlag int
	stashFlag       bool
	cachePath       string
	cacheFile       *os.File

	invalidPath bool = false
)

type PathRecords struct {
	Records      map[string]PathRecord  `json:"paths"`
	StashRecords map[string]StashRecord `json:"stash"`
}

type PathRecord struct {
	Count     int    `json:"count"`
	Timestamp string `json:"ts"`
}

type StashRecord struct {
	Timestamp string `json:"ts"`
}

func main() {
	log.Printf("ucd-v0.1\n")

	// flags
	flag.BoolVar(&clearFlag, "c", false, "clear history list")
	flag.BoolVar(&listFlag, "l", false, "MRU list for recently used cd commands")
	flag.BoolVar(&listStashFlag, "ls", false, "list stashed cd commands")
	flag.BoolVar(&stashFlag, "s", false, "stash cd path into a separate list")
	flag.IntVar(&repeatFlag, "r", 1, "repeat dynamic cd path (for ..)")
	flag.IntVar(&historyPathFlag, "p", 0, "execute the # path listed from MRU list")
	flag.Parse()

	args := flag.Args()
	log.Printf("args: %v\n", args)
	log.Printf("listFlag: %v\n", listFlag)

	homeDir, _ := os.UserHomeDir()
	curDir, _ := os.Getwd()

	log.Printf("~: %v\n", homeDir)
	log.Printf("cwd: %v\n", curDir)

	cachePath = homeDir + "/.ucd-cache"

	// cacheFile, _ = os.OpenFile(cachePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	log.Printf("Configuring cachePath to %v\n", cachePath)

	cacheFile, _ := os.Open(cachePath)
	byteValue, _ := ioutil.ReadAll(cacheFile)

	var pr PathRecords
	err := json.Unmarshal(byteValue, &pr)
	if err != nil {
		pr = PathRecords{
			Records:      map[string]PathRecord{},
			StashRecords: map[string]StashRecord{},
		}
	}

	if clearFlag {
		pr = PathRecords{
			Records:      map[string]PathRecord{},
			StashRecords: map[string]StashRecord{},
		}
		output, _ := json.Marshal(pr)
		ioutil.WriteFile(cachePath, output, 0644)
		fmt.Print(".")
		os.Exit(1)
	}

	// exit earlier depending on flag passed in
	if listFlag {
		listRecords(pr)
		fmt.Print(".")

		os.Exit(1)
	}

	if listStashFlag {
		listStashRecords(pr)
		fmt.Print(".")

		os.Exit(1)
	}

	if len(args) > 1 {
		log.Fatalln("Only < 1 arguments can be passed to ucd")
	}

	// fmt.Print sends output to stdout, this will be consumed by builtin `cd` command

	var targetPath string
	if historyPathFlag > 0 {
		mruRecords := sortedRecordKeys(pr)
		targetPath = mruRecords[historyPathFlag-1]
	} else {
		if len(args) > 0 {
			targetPath = repeat(args[0], repeatFlag)
		} else {
			targetPath = homeDir
		}
	}
	log.Printf("targetPath: %v\n", targetPath)

	// attempt to chdir into target path
	err = os.Chdir(targetPath)
	if err != nil {
		invalidPath = true
	}
	targetPath, _ = os.Getwd()

	if invalidPath {
		fmt.Print(targetPath)
		os.Exit(1)
	}

	rec, ok := pr.Records[targetPath]
	if ok {
		rec.Count++
		rec.Timestamp = timeNow()
		pr.Records[targetPath] = rec
	} else {
		pr.Records[targetPath] = PathRecord{Count: 1, Timestamp: timeNow()}
	}

	fmt.Print(targetPath)
	if stashFlag {
		pr.StashRecords[targetPath] = StashRecord{Timestamp: timeNow()}
	}

	output, _ := json.Marshal(pr)
	ioutil.WriteFile(cachePath, output, 0644)
	cacheFile.Close()
}

func sortedRecordKeys(pr PathRecords) []string {
	keys := make([]string, 0, len(pr.Records))
	for key := range pr.Records {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return pr.Records[keys[i]].Timestamp > pr.Records[keys[j]].Timestamp
	})

	return keys
}

func sortedStashRecordKeys(pr PathRecords) []string {
	keys := make([]string, 0, len(pr.StashRecords))
	for key := range pr.StashRecords {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return pr.StashRecords[keys[i]].Timestamp > pr.StashRecords[keys[j]].Timestamp
	})

	return keys
}

func listRecords(pr PathRecords) {
	t := table.NewWriter()
	t.SetOutputMirror(log.Writer())
	t.AppendHeader(table.Row{"#", "path", "count", "timestamp"})

	keys := sortedRecordKeys(pr)

	index := 1

	for _, key := range keys {
		t.AppendRow([]interface{}{
			index,
			key,
			pr.Records[key].Count,
			pr.Records[key].Timestamp,
		})
		index++
	}

	t.Render()
}

func listStashRecords(pr PathRecords) {
	t := table.NewWriter()
	t.SetOutputMirror(log.Writer())
	t.AppendHeader(table.Row{"#", "path", "timestamp"})

	keys := sortedStashRecordKeys(pr)

	index := 1

	for _, key := range keys {
		t.AppendRow([]interface{}{
			index,
			key,
			pr.StashRecords[key].Timestamp,
		})
		index++
	}

	t.Render()
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
