package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/sedimentary/lib"
)

var initCmd = &cobra.Command{
	Use:   "new",
	Short: "initialize a new sedimentary",
	RunE: func(_ *cobra.Command, _ []string) error {
		wd, newErr := os.Getwd()
		if newErr != nil {
			return newErr
		}

		sdmt, newErr := lib.New(wd)
		if newErr != nil {
			return newErr
		}

		if err := sdmt.WriteVersion(); err != nil {
			return err
		}

		if closeErr := sdmt.Close(); closeErr != nil {
			return errors.Join(closeErr, newErr)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
