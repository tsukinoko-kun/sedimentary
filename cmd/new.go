package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/sedimentary/lib/libdb"
)

var initCmd = &cobra.Command{
	Use:   "new",
	Short: "initialize a new sedimentary",
	RunE: func(_ *cobra.Command, _ []string) error {
		db, err := libdb.Init()
		if db != nil {
			defer func() {
				if err := db.Close(); err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				}
			}()
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
