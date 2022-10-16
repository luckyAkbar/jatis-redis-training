package console

import (
	"strconv"

	"github.com/luckyAkbar/jatis-redis-training/internal/config"
	"github.com/luckyAkbar/jatis-redis-training/internal/db"

	utils "github.com/kumparan/go-utils"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var runMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "run database migration",
	Long:  "Used to run databse migration defined in migration folder",
	Run:   migration,
}

func init() {
	runMigrate.PersistentFlags().Int("step", 0, "maximum migration steps")
	runMigrate.PersistentFlags().String("direction", "up", "migration direction")
	RootCmd.AddCommand(runMigrate)
}

func migration(cmd *cobra.Command, _ []string) {
	direction := cmd.Flag("direction").Value.String()
	stepStr := cmd.Flag("step").Value.String()
	step, err := strconv.Atoi(stepStr)
	if err != nil {
		logrus.WithField("stepStr", stepStr).Fatal("Failed to parse step to int: ", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "./db/migration",
	}

	migrate.SetTable("schema_migrations")

	db.InitializePostgresConn()

	pgdb, err := db.PostgresDB.DB()
	if err != nil {
		logrus.WithField("DatabaseDSN", config.PostgresDSN()).Fatal("failed to run migration")
	}

	var n int
	if direction == "down" {
		n, err = migrate.ExecMax(pgdb, "postgres", migrations, migrate.Down, step)
	} else {
		n, err = migrate.ExecMax(pgdb, "postgres", migrations, migrate.Up, step)
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"migrations": utils.Dump(migrations),
			"direction":  direction}).
			Fatal("Failed to migrate database: ", err)
	}

	logrus.Infof("Applied %d migrations!\n", n)
}
