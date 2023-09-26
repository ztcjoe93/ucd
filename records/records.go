package records

import (
	"log"
	"sort"

	"github.com/jedib0t/go-pretty/table"
)

type Record interface {
	PathRecord | StashRecord
	GetTimestamp() string
	HasCount() bool
}

type Records struct {
	PathRecords  map[string]PathRecord  `json:"paths"`
	StashRecords map[string]StashRecord `json:"stash"`
}

type PathRecord struct {
	Count     int    `json:"count"`
	Timestamp string `json:"ts"`
}

func (pr PathRecord) GetTimestamp() string {
	return pr.Timestamp
}

func (pr PathRecord) HasCount() bool {
	return true
}

type StashRecord struct {
	Timestamp string `json:"ts"`
	Alias     string `json:"alias"`
}

func (sr StashRecord) GetTimestamp() string {
	return sr.Timestamp
}

func (sr StashRecord) HasCount() bool {
	return false
}

func (r Records) ListRecords(recType string, maxLimit int) {

	var isPath = true

	if recType == "stash" {
		isPath = false
	}

	t := table.NewWriter()
	t.SetOutputMirror(log.Writer())

	if isPath {
		t.AppendHeader(table.Row{"#", "path", "count", "timestamp"})
	} else {
		t.AppendHeader(table.Row{"#", "alias", "path", "timestamp"})
	}

	var keys []string

	if isPath {
		keys = SortRecords(r.PathRecords)
	} else {
		keys = SortRecords(r.StashRecords)
	}

	displayLimit := len(keys)
	if maxLimit < displayLimit {
		displayLimit = maxLimit
	}

	for i := 0; i < displayLimit; i++ {
		if isPath {
			t.AppendRow([]interface{}{
				i + 1, keys[i], r.PathRecords[keys[i]].Count, r.PathRecords[keys[i]].Timestamp,
			})
		} else {
			t.AppendRow([]interface{}{
				i + 1, r.StashRecords[keys[i]].Alias, keys[i], r.StashRecords[keys[i]].Timestamp,
			})
		}
	}
	/**
	for _, key := range keys {
		if isPath {
			t.AppendRow([]interface{}{index, key, r.PathRecords[key].Count, r.PathRecords[key].Timestamp})
		} else {
			t.AppendRow([]interface{}{index, r.StashRecords[key].Alias, key, r.StashRecords[key].Timestamp})
		}
		index++
	}
	**/

	t.Render()
}

func SortRecords[v Record](r map[string]v) []string {
	keys := make([]string, 0, len(r))
	for key := range r {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return r[keys[i]].GetTimestamp() > r[keys[j]].GetTimestamp()
	})

	return keys
}
