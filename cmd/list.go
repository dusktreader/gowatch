package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/dusktreader/gowatch/timer"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.PersistentFlags().BoolP("full", "f", false, "Show the full timers")
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:	"list",
	Short:	"List all timers",
	Long:	"List all timers with their names and durations",
	Args:	cobra.MaximumNArgs(1),
	Run:	listMain,
}

func listMain(cmd *cobra.Command, args []string){
	full, err := cmd.Flags().GetBool("full")
	MaybeDie(err)

	cacheDir := timer.GetCacheDir()

	slog.Debug("Loading all timers")
	nts, err := timer.LoadAll(cacheDir)
	MaybeDie(err)

	if len(nts) == 0 {
		fmt.Fprintln(os.Stderr, "No timers found")
	}

	slog.Debug("Computing alignment for names")
	maxWidth := 0
	for _, nt := range nts {
		maxWidth = max(maxWidth, len(nt.Name))
	}

	slog.Debug("Listing timers", "Count", len(nts))
	for _, nt := range nts {
		if full {
			slog.Debug("Showing full timer", "Name", nt.Name)
			fmt.Printf("%*s: %s\n", maxWidth, nt.Name, nt.Ticks)
		} else {
			slog.Debug("Showing compact timer", "Name", nt.Name)
			fmt.Printf("%*s: %s\n", maxWidth, nt.Name, nt.Ticks.ElapsedString())
		}
	}
}
