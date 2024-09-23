package database

import (
	"fmt"
	"sync"
)

type Table struct {
	columns []string
	records []Record
	mu      sync.RWMutex
}

func NewTable(columns []string) *Table {
	return &Table{
		columns: columns,
		records: make([]Record, 0),
	}
}

func (t *Table) Insert(values []string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(values) != len(t.columns) {
		return fmt.Errorf("invalid number of values: expected %d, got %d", len(t.columns), len(values))
	}

	record := make(Record)
	for i, col := range t.columns {
		record[col] = values[i]
	}

	t.records = append(t.records, record)
	return nil
}

func (t *Table) Select(condition func(Record) bool) []Record {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]Record, 0)
	for _, record := range t.records {
		if condition == nil || condition(record) {
			result = append(result, record.Copy())
		}
	}
	return result
}

func (t *Table) Update(updates map[string]string, condition func(Record) bool) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	count := 0
	for i, record := range t.records {
		if condition == nil || condition(record) {
			for col, val := range updates {
				if _, exists := record[col]; exists {
					t.records[i][col] = val
				}
			}
			count++
		}
	}
	return count
}

func (t *Table) Delete(condition func(Record) bool) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	count := 0
	newRecords := make([]Record, 0)
	for _, record := range t.records {
		if condition == nil || condition(record) {
			count++
		} else {
			newRecords = append(newRecords, record)
		}
	}
	t.records = newRecords
	return count
}

func (t *Table) Columns() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return append([]string{}, t.columns...)
}

func (t *Table) RecordCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.records)
}