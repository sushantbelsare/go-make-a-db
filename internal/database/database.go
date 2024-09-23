package database

import (
	_ "errors"
	"fmt"
	"sync"
)

type Database struct {
	tables map[string]*Table
	mu     sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		tables: make(map[string]*Table),
	}
}

func (db *Database) CreateTable(name string, columns []string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.tables[name]; exists {
		return fmt.Errorf("table '%s' already exists", name)
	}

	db.tables[name] = NewTable(columns)
	return nil
}

func (db *Database) DropTable(name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.tables[name]; !exists {
		return fmt.Errorf("table '%s' does not exist", name)
	}

	delete(db.tables, name)
	return nil
}

func (db *Database) ListTables() []string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	tables := make([]string, 0, len(db.tables))
	for name := range db.tables {
		tables = append(tables, name)
	}
	return tables
}

func (db *Database) GetTable(name string) (*Table, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	table, exists := db.tables[name]
	if !exists {
		return nil, fmt.Errorf("table '%s' does not exist", name)
	}
	return table, nil
}

func (db *Database) InsertRecord(tableName string, values []string) error {
	table, err := db.GetTable(tableName)
	if err != nil {
		return err
	}

	return table.Insert(values)
}

func (db *Database) SelectRecords(tableName string, condition func(Record) bool) ([]Record, error) {
	table, err := db.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	return table.Select(condition), nil
}

func (db *Database) UpdateRecords(tableName string, updates map[string]string, condition func(Record) bool) (int, error) {
	table, err := db.GetTable(tableName)
	if err != nil {
		return 0, err
	}

	return table.Update(updates, condition), nil
}

func (db *Database) DeleteRecords(tableName string, condition func(Record) bool) (int, error) {
	table, err := db.GetTable(tableName)
	if err != nil {
		return 0, err
	}

	return table.Delete(condition), nil
}