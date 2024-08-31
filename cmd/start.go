package cmd

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start API server",
		Run: func(cmd *cobra.Command, args []string) {
			srv, err := di.GetAPIServer()
			if err != nil {
				log.Fatal("init server error:", err)
			}
			log.Println("Starting server at", srv.Addr)
			err = srv.ListenAndServe()
			if err != nil {
				log.Fatal("unable to start server:", err)
			}
		},
	}
	rootCmd.AddCommand(cmd)
}
