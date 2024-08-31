package cmd

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

func init() {
	var dir string
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migration",
		Run: func(cmd *cobra.Command, args []string) {
			db, err := di.GetDatabase()
			if err != nil {
				log.Fatal(err)
			}

			driver, err := postgres.WithInstance(db, &postgres.Config{})
			if err != nil {
				log.Fatal(err)
			}

			if dir == "" {
				dir, err = os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				dir += "/files/db_schema/"
			}

			m, err := migrate.NewWithDatabaseInstance("file://"+dir, "postgres", driver)
			if err != nil {
				log.Fatal("file error:", err)
			}

			err = m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
			if err != nil {
				log.Fatal("up error: ", err)
			}
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory of migration files")
	rootCmd.AddCommand(cmd)
}
