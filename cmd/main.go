package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sushantbelsare/go-make-a-db/internal/cli"
	"github.com/sushantbelsare/go-make-a-db/internal/database"
)

func main() {
	// Create a new database
	db := database.NewDatabase()

	// Create a new CLI instance
	cli := cli.NewCLI(db)

	// Check if there are command-line arguments
	if len(os.Args) > 1 {
		// Execute the command and exit
		if err := cli.ExecuteCommand(os.Args[1:]); err != nil {
			log.Fatal(err)
		}
		return
	}

	// If no arguments, start interactive mode
	fmt.Println("Welcome to SimpleDB. Type 'help' for a list of commands.")
	for {
		if err := cli.InteractiveMode(); err != nil {
			if err.Error() == "exit" {
				fmt.Println("Goodbye!")
				return
			}
			fmt.Printf("Error: %v\n", err)
		}
	}
}