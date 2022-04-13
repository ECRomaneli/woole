package dashboard

import (
	"time"
	"woole/recorder"
)

type Item struct {
	ID      string        `json:"id"`
	Path    string        `json:"path"`
	Method  string        `json:"method"`
	Status  int           `json:"status"`
	Elapsed time.Duration `json:"elapsed"`
}

type Items []Item

func (items *Items) FromRecords(records *recorder.Records) *Items {
	var slice Items

	records.Each(func(rec *recorder.Record) {
		slice = append(slice, *(&Item{}).FromRecord(rec))
	})

	return &slice
}

func (item *Item) FromRecord(rec *recorder.Record) *Item {
	if rec == nil {
		return nil
	}

	item.ID = rec.ID
	item.Path = rec.Request.Path
	item.Method = rec.Request.Method
	item.Status = rec.Response.Code
	item.Elapsed = rec.Elapsed

	return item
}
