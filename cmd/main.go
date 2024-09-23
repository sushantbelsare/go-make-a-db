package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sushantbelsare/go-make-a-db/internal/cli"
	"github.com/sushantbelsare/go-make-a-db/internal/database"
)

const dbFilename = "database.json"

func main() {
	db := database.NewDatabase()

	// Load the database from file
	err := db.LoadFromFile(dbFilename)
	if err != nil {
		log.Fatalf("Failed to load database: %v", err)
	}

	c := cli.NewCLI(db)

	// Set up a channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// go func() {
	// 	sig := <-sigChan
	// 	fmt.Printf("\nReceived %s signal. Saving database...\n", sig)

		

	// 	os.Exit(0) // Exit gracefully
	// }()

	fmt.Println("Welcome to SimpleDB. Type 'help' for a list of commands.")
	if err := c.InteractiveMode(); err != nil {
		log.Fatalf("Error running CLI: %v", err)
	}
}