package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/sedimentary/lib"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "resets the working directory to the tracked state",
	RunE: func(_ *cobra.Command, _ []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		sdmt, err := lib.Open(wd)
		if err != nil {
			return err
		}

		if err := sdmt.Reset(); err != nil {
			return err
		}

		if err := sdmt.Close(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
