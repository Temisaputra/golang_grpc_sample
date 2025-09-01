package cmd

import (
	"log"

	"github.com/Temisaputra/warOnk/internal/domain/entity"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(_ *cobra.Command, _ []string) {
		log.Println("use migrate [up|down|fresh] to run migrations")
	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all migrations",
	Run: func(_ *cobra.Command, _ []string) {
		startMigrate("up")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Drop all tables",
	Run: func(_ *cobra.Command, _ []string) {
		startMigrate("down")
	},
}

var migrateFreshCmd = &cobra.Command{
	Use:   "fresh",
	Short: "Drop all tables and re-run migrations",
	Run: func(_ *cobra.Command, _ []string) {
		startMigrate("fresh")
	},
}

func startMigrate(migrationType string) {
	deps := InitDependencies() // 🔑 ambil dari bootstrap.go
	db := deps.DB

	switch migrationType {
	case "up":
		if err := entity.Migrate(db); err != nil {
			log.Fatalf("migration up failed: %v", err)
		}
		log.Println("✅ migration up success")
	case "down":
		if err := entity.Drop(db); err != nil {
			log.Fatalf("migration down failed: %v", err)
		}
		log.Println("✅ migration down success")
	case "fresh":
		if err := entity.Drop(db); err != nil {
			log.Fatalf("migration fresh (drop) failed: %v", err)
		}
		if err := entity.Migrate(db); err != nil {
			log.Fatalf("migration fresh (migrate) failed: %v", err)
		}
		log.Println("✅ migration fresh success")
	}
}
