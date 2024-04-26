package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/sedimentary/build"
	"github.com/tsukinoko-kun/sedimentary/lib"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version of this binary and its libraries",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Printf("Version %s\n", build.Version)
		fmt.Printf("Commit  %s\n", build.Commit)
		fmt.Printf("Date    %s\n", build.Date)

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		sdmt, err := lib.Open(wd)
		if err != nil {
			return err
		}

		if v, err := sdmt.ReadVersion(); err != nil {
			_ = sdmt.Close()
			return err
		} else {
			fmt.Printf("Created %s\n", v)
		}

		if err := sdmt.Close(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
