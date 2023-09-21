package main

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
}

func (sr StashRecord) GetTimestamp() string {
	return sr.Timestamp
}

func (sr StashRecord) HasCount() bool {
	return false
}

func (r Records) listRecords(recType string) {

	var isPath = true

	if recType == "stash" {
		isPath = false
	}

	t := table.NewWriter()
	t.SetOutputMirror(log.Writer())

	if isPath {
		t.AppendHeader(table.Row{"#", "path", "count", "timestamp"})
	} else {
		t.AppendHeader(table.Row{"#", "path", "timestamp"})
	}

	var keys []string

	if isPath {
		keys = sortRecords(r.PathRecords)
	} else {
		keys = sortRecords(r.StashRecords)
	}

	index := 1
	for _, key := range keys {
		if isPath {
			t.AppendRow([]interface{}{index, key, r.PathRecords[key].Count, r.PathRecords[key].Timestamp})
		} else {
			t.AppendRow([]interface{}{index, key, r.StashRecords[key].Timestamp})
		}
		index++
	}

	t.Render()
}

func sortRecords[v Record](r map[string]v) []string {
	keys := make([]string, 0, len(r))
	for key := range r {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return r[keys[i]].GetTimestamp() > r[keys[j]].GetTimestamp()
	})

	return keys
}
