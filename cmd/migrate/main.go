package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/config"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
)

func main() {
	var command string
	flag.StringVar(&command, "cmd", "up", "Migration command: up, down, rollback")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	switch command {
	case "up":
		log.Println("Running migrations...")
		if err := database.RunMigrations(cfg); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("✓ Migrations completed successfully")

	case "down":
		log.Println("Rolling back all migrations...")
		if err := database.MigrateDown(cfg); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		log.Println("✓ All migrations rolled back successfully")

	case "rollback":
		log.Println("Rolling back last migration...")
		if err := database.RollbackMigration(cfg); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		log.Println("✓ Last migration rolled back successfully")

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		fmt.Fprintf(os.Stderr, "Available commands: up, down, rollback\n")
		os.Exit(1)
	}
}
