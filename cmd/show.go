package cmd

import (
	"fmt"
	"log/slog"

	"github.com/dusktreader/gowatch/timer"
	"github.com/spf13/cobra"
)

func init() {
	showCmd.PersistentFlags().BoolP("full", "f", false, "Show the full timer")
	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:	"show",
	Short:	"Show a timer",
	Long:	"Show a named timer",
	Args:	cobra.MaximumNArgs(1),
	Run:	showMain,
}

func showMain(cmd *cobra.Command, args []string){
	full, err := cmd.Flags().GetBool("full")
	MaybeDie(err)

	var name string
	if len(args) == 0 {
		name = timer.DEFAULT_TIMER_NAME
	} else {
		name = args[0]
	}

	cacheDir := timer.GetCacheDir()

	t, err := timer.Load(name, cacheDir, true)
	MaybeDie(err)

	if full {
		slog.Debug("Showing full timer", "Name", name)
		fmt.Println(t)
	} else {
		slog.Debug("Showing compact timer", "Name", name)
		fmt.Println(t.ElapsedString())
	}
}
