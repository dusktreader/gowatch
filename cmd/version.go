package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func init() {
	versionCmd.PersistentFlags().BoolP("full", "f", false, "Show the full version info")
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:	"version",
	Short:	"Show the version",
	Long:	"Show the version",
	Run:	versionMain,
}

func versionMain(cmd *cobra.Command, args []string){
	_, err := cmd.Flags().GetBool("full")
	MaybeDie(err)

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		Die("Couldn't read build info")
	}

	fmt.Println(buildInfo.Main.Version)
}
