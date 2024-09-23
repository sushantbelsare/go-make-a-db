package cli

import (
	_ "bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/sushantbelsare/go-make-a-db/internal/database"

	"io"
)

type CLI struct {
	db      *database.Database
	history []string
}

func NewCLI(db *database.Database) *CLI {
	return &CLI{
		db:      db,
		history: make([]string, 0),
	}
}

func (c *CLI) InteractiveMode() error {
	rl, err := readline.New("simple-db> ")
	if err != nil {
		return err
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // handle EOF or Ctrl+C
			if err == readline.ErrInterrupt || err == io.EOF {
				fmt.Println("Goodbye!")
				os.Exit(0)
				break // Exit the loop
			}
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		c.history = append(c.history, line)

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		err = c.ExecuteCommand(args)
		if err != nil {
			if strings.ToLower(args[0]) == "exit" {
				fmt.Println("Goodbye!")
				os.Exit(0)
				break // Exit the loop
			}
			fmt.Printf("Error: %v\n", err)
		}
	}

	return nil
}

func (c *CLI) ExecuteCommand(args []string) error {
	command := strings.ToLower(args[0])

	switch command {
	case "create":
		return c.handleCreate(args[1:])
	case "drop":
		return c.handleDrop(args[1:])
	case "list":
		return c.handleList(args[1:])
	case "insert":
		return c.handleInsert(args[1:])
	case "select":
		return c.handleSelect(args[1:])
	case "update":
		return c.handleUpdate(args[1:])
	case "delete":
		return c.handleDelete(args[1:])
	case "help":
		c.printHelp()
		return nil
	case "exit":
		return fmt.Errorf("exit")
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (c *CLI) handleCreate(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: create <table_name> <column1> <column2> ...")
	}

	tableName := args[0]
	columns := args[1:]

	err := c.db.CreateTable(tableName, columns)
	if err != nil {
		return err
	}

	fmt.Printf("Table '%s' created successfully.\n", tableName)
	return nil
}

func (c *CLI) handleDrop(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: drop <table_name>")
	}

	tableName := args[0]

	err := c.db.DropTable(tableName)
	if err != nil {
		return err
	}

	fmt.Printf("Table '%s' dropped successfully.\n", tableName)
	return nil
}

func (c *CLI) handleList(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: list")
	}

	tables := c.db.ListTables()
	if len(tables) == 0 {
		fmt.Println("No tables found.")
		return nil
	}

	fmt.Println("Tables:")
	for _, table := range tables {
		fmt.Printf("- %s\n", table)
	}
	return nil
}

func (c *CLI) handleInsert(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: insert <table_name> <value1> <value2> ...")
	}

	tableName := args[0]
	values := args[1:]

	err := c.db.InsertRecord(tableName, values)
	if err != nil {
		return err
	}

	fmt.Println("Record inserted successfully.")
	return nil
}

func (c *CLI) handleSelect(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: select <table_name> [<column>=<value>]")
	}

	tableName := args[0]
	var condition func(database.Record) bool

	if len(args) > 1 {
		parts := strings.SplitN(args[1], "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid condition format")
		}
		column, value := parts[0], parts[1]
		condition = func(r database.Record) bool {
			return r[column] == value
		}
	}

	records, err := c.db.SelectRecords(tableName, condition)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		fmt.Println("No records found.")
		return nil
	}

	c.printRecords(records)
	return nil
}

func (c *CLI) handleUpdate(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: update <table_name> <column>=<value> <condition_column>=<condition_value>")
	}

	tableName := args[0]
	updateParts := strings.SplitN(args[1], "=", 2)
	conditionParts := strings.SplitN(args[2], "=", 2)

	if len(updateParts) != 2 || len(conditionParts) != 2 {
		return fmt.Errorf("invalid update or condition format")
	}

	updates := map[string]string{updateParts[0]: updateParts[1]}
	condition := func(r database.Record) bool {
		return r[conditionParts[0]] == conditionParts[1]
	}

	count, err := c.db.UpdateRecords(tableName, updates, condition)
	if err != nil {
		return err
	}

	fmt.Printf("%d record(s) updated successfully.\n", count)
	return nil
}

func (c *CLI) handleDelete(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: delete <table_name> <column>=<value>")
	}

	tableName := args[0]
	parts := strings.SplitN(args[1], "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid condition format")
	}

	column, value := parts[0], parts[1]
	condition := func(r database.Record) bool {
		return r[column] == value
	}

	count, err := c.db.DeleteRecords(tableName, condition)
	if err != nil {
		return err
	}

	fmt.Printf("%d record(s) deleted successfully.\n", count)
	return nil
}

func (c *CLI) printRecords(records []database.Record) {
	if len(records) == 0 {
		return
	}

	columns := records[0].Columns()
	for _, col := range columns {
		fmt.Printf("%-15s", col)
	}
	fmt.Println()

	for _, record := range records {
		for _, col := range columns {
			fmt.Printf("%-15s", record[col])
		}
		fmt.Println()
	}
}

func (c *CLI) printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  create <table_name> <column1> <column2> ...  Create a new table")
	fmt.Println("  drop <table_name>                            Drop a table")
	fmt.Println("  list                                         List all tables")
	fmt.Println("  insert <table_name> <value1> <value2> ...    Insert a new record")
	fmt.Println("  select <table_name> [<column>=<value>]       Select records")
	fmt.Println("  update <table_name> <col>=<val> <cond_col>=<cond_val>  Update records")
	fmt.Println("  delete <table_name> <column>=<value>         Delete records")
	fmt.Println("  help                                         Show this help message")
	fmt.Println("  exit                                         Exit the program")
}
