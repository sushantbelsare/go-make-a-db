package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/sushantbelsare/go-make-a-db/internal/config"
	"github.com/sushantbelsare/go-make-a-db/internal/util"
)

type Database struct {
	tables map[string]*Table
	mu     sync.RWMutex
	wal    *util.WAL
}

func NewDatabase() *Database {
	return &Database{
		tables: make(map[string]*Table),
	}
}

// SaveToFile saves the database state to a JSON file
func (db *Database) SaveToFile(filename string) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	data, err := json.MarshalIndent(db.tables, "", "  ")
	if err != nil {
		return err
	}

	key := []byte(config.GetEnv("SEC_KEY", "default")) // Ensure this key is 32 bytes long for AES-256
	encryptedData, err := util.Encrypt(data, key)

	if err != nil {
		return err
	}

	db.wal.Close()

	return os.WriteFile(filename, []byte(encryptedData), 0644)
}

// LoadFromFile loads the database state from a JSON file
func (db *Database) LoadFromFile(filename string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.wal, _ = util.NewWAL(config.GetEnv("WAL_FILE_NAME", "default"))

	if db.wal == nil {
		return errors.New("WAL can't be opened")
	}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Starting session with a clean database...")
			return nil // No existing file is not an error
		}
		log.Println("Starting session with an exisitng database...")
		return err
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	key := []byte(config.GetEnv("SEC_KEY", "default")) // Use the same key as for encryption
	decryptedData, err := util.Decrypt(string(data), key)

	if err != nil {
		return err
	}

	var tables map[string]*Table
	if err := json.Unmarshal(decryptedData, &tables); err != nil {
		return err
	}

	db.tables = tables
	return nil
}

func (db *Database) CreateTable(name string, columns []string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.tables[name]; exists {
		return fmt.Errorf("table '%s' already exists", name)
	}

	entry := util.LogEntry {
		Operation: "create",
		TableName: name,
		Values: columns,
	}

	if err := db.wal.WriteEntry(entry); err != nil {
		return err
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

	entry := util.LogEntry {
		Operation: "drop",
		TableName: name,
	}

	if err := db.wal.WriteEntry(entry); err != nil {
		return err
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

	entry := util.LogEntry {
		Operation: "insert",
		TableName: tableName,
		Values: values,
	}

	if err := db.wal.WriteEntry(entry); err != nil {
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

	entry := util.LogEntry {
		Operation: "update",
		TableName: tableName,
		Values: updates,
		Condition: condition,
	}

	if err := db.wal.WriteEntry(entry); err != nil {
		return -1, err
	}

	return table.Update(updates, condition), nil
}

func (db *Database) DeleteRecords(tableName string, condition func(Record) bool) (int, error) {
	table, err := db.GetTable(tableName)
	if err != nil {
		return 0, err
	}

	entry := util.LogEntry {
		Operation: "delete",
		TableName: tableName,
		Condition: condition,
	}

	if err := db.wal.WriteEntry(entry); err != nil {
		return -1, err
	}

	return table.Delete(condition), nil
}

func (db *Database) RecoverFromWAL(wal *util.WAL) error {
	entries, err := wal.ReadEntries()
	if err != nil {
		return fmt.Errorf("failed to read WAL entries: %w", err)
	}

	for _, entry := range entries {
		switch entry.Operation {
		case "insert":
			db.InsertRecord(entry.TableName, entry.Values.([]string))
			// Handle other operations like update and delete.
		default:
			return fmt.Errorf("unknown operation in WAL: %s", entry.Operation)
		}
	}

	return nil
}