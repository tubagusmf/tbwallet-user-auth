package console

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/tubagusmf/tbwallet-user-auth/internal/config"
	"github.com/tubagusmf/tbwallet-user-auth/internal/helper"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

var (
	direction string
	step      int
)

func init() {
	rootCmd.AddCommand(migrationCMD)

	migrationCMD.Flags().StringVarP(&direction, "direction", "d", "up", "Migration direction (up or down)")
	migrationCMD.Flags().IntVarP(&step, "step", "s", 1, "Number of migrations to apply")
}

var migrationCMD = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `This command is used to apply or rollback database migrations.`,
	Run:   migrateDB,
}

func migrateDB(cmd *cobra.Command, args []string) {
	config.LoadWithViper()

	connDB, err := sql.Open("postgres", helper.GetConnectionString())
	log.Printf("Connection String: %s", helper.GetConnectionString())

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer connDB.Close()

	migrations := &migrate.FileMigrationSource{Dir: "./database/migrations"}

	var n int
	if direction == "down" {
		n, err = migrate.ExecMax(connDB, "postgres", migrations, migrate.Down, step)
	} else {
		n, err = migrate.ExecMax(connDB, "postgres", migrations, migrate.Up, step)
	}

	if err != nil {
		log.Fatalf("Error applying migrations: %v", err)
	}

	log.Printf("Successfully applied %d migrations", n)
}
