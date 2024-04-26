package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/sedimentary/lib"
)

var (
	m         string
	commitCmd = &cobra.Command{
		Use:   "commit",
		Short: "add changes to the sedimentary",
		RunE: func(_ *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			sdmt, err := lib.Open(wd)
			if err != nil {
				return err
			}

			if err := sdmt.Commit(args, m); err != nil {
				_ = sdmt.Close()
				return err
			}

			if err := sdmt.Close(); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	commitCmd.Flags().StringVarP(&m, "message", "m", "", "use the given message as the commit message")
	rootCmd.AddCommand(commitCmd)
}
