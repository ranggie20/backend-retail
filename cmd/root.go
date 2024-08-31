package cmd

import (
	"log"

	"github.com/online-bnsp/backend/dep"
	"github.com/spf13/cobra"
)

var configFile string
var di *dep.DI
var mbi *dep.MBI

var rootCmd = &cobra.Command{
	Use:   "kaos",
	Short: "CeritaKaos backend application",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		di, err = dep.InitDI(configFile)
		if err != nil {
			log.Fatal(err)
		}

		mbi, err = dep.InitMBI(configFile)
		if err != nil {
			log.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	rootCmd.AddCommand(&cobra.Command{Use: "completion", Hidden: true})
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Config file to use")
	return rootCmd.Execute()
}
