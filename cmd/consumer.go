package cmd

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/online-bnsp/backend/constant"
	"github.com/online-bnsp/backend/consumer"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "consumer",
		Short: "Start Consumer server",
		Run: func(cmd *cobra.Command, args []string) {
			db, err := di.GetDatabase()
			if err != nil {
				log.Fatal("init server error:", err)
			}

			handlers := consumer.New(db)
			// register all consumers below
			mbi.Register("Calculate Coin Views", constant.SampleConsumer, "cerita_kaos", handlers.SampleConsumer) // sample

			// run all consumers
			mbi.Run()

			log.Println("Consumer server started")

			// wait
			mbi.Wait()
		},
	}
	rootCmd.AddCommand(cmd)
}
