package util

import (
	"bufio"
	"encoding/json"
	"os"
	"time"
)

type LogEntry struct {
	Operation string   `json:"operation"`
	TableName string   `json:"table_name"`
	Values    interface{} `json:"values,omitempty"`
	Condition interface{}   `json:"condition,omitempty"`
	Time 	  time.Time 			`json:"created_at"`
}

type WAL struct {
	file *os.File
}

func NewWAL(filename string) (*WAL, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &WAL{file: file}, nil
}

func (wal *WAL) WriteEntry(entry LogEntry) error {
	entry.Time = time.Now()
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	if _, err := wal.file.Write(append(data, '\n')); err != nil {
		return err
	}
	return wal.file.Sync() // Ensure data is flushed to disk
}

func (wal *WAL) ReadEntries() ([]LogEntry, error) {
	file, err := os.Open(wal.file.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func (wal *WAL) Close() error {
	return wal.file.Close()
}
