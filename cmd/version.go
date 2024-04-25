package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/sedimentary/build"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version of this binary and its libraries",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Printf("Version %s\n", build.Version)
		fmt.Printf("Commit  %s\n", build.Commit)
		fmt.Printf("Date    %s\n", build.Date)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
