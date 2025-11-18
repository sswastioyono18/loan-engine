package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
	"github.com/sswastioyono18/loan-engine/pkg/util"

	_ "github.com/lib/pq" // Import PostgreSQL driver
)

func main() {
	// Define command-line flags
	var (
		action = flag.String("action", "up", "Migration action: up, down, status")
		dir    = flag.String("dir", "./migrations", "Directory containing migration files")
		help   = flag.Bool("help", false, "Show help message")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Initialize database connection
	db, err := util.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Get the underlying *sql.DB from our custom DB wrapper
	sqlDB := db.GetDB()

	// Set the dialect for goose
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("Failed to set goose dialect:", err)
	}

	// Execute migration action
	switch *action {
	case "up":
		if err := goose.Up(sqlDB, *dir); err != nil {
			log.Fatal("Migration failed:", err)
		}
		log.Println("Migrations applied successfully")
	case "down":
		if err := goose.Down(sqlDB, *dir); err != nil {
			log.Fatal("Migration rollback failed:", err)
		}
		log.Println("Migration rolled back successfully")
	case "status":
		if err := goose.Status(sqlDB, *dir); err != nil {
			log.Fatal("Failed to get migration status:", err)
		}
	default:
		log.Fatalf("Unknown action: %s. Use 'up', 'down', or 'status'", *action)
	}
}

func showHelp() {
	fmt.Println("Loan Engine Migration Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  migrate [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -action    Migration action: up, down, status (default: up)")
	fmt.Println("  -dir       Directory containing migration files (default: ./migrations)")
	fmt.Println("  -help      Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  migrate -action up                    # Apply all pending migrations")
	fmt.Println("  migrate -action down                  # Rollback last migration")
	fmt.Println("  migrate -action status                # Show migration status")
	fmt.Println("  migrate -action up -dir ./my-migrations # Apply migrations from custom directory")
}
