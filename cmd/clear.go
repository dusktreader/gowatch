package cmd

import (
	"log/slog"

	"github.com/dusktreader/gowatch/timer"
	"github.com/spf13/cobra"
)

func init() {
	clearCmd.PersistentFlags().BoolP("all", "A", false, "Remove all timers")
	rootCmd.AddCommand(clearCmd)
}

var clearCmd = &cobra.Command{
	Use:	"clear",
	Short:	"Clear timers",
	Long:	"Clear timers",
	Args:	cobra.MaximumNArgs(1),
	Run:	clearMain,
}

func clearMain(cmd *cobra.Command, args []string){
	all, err := cmd.Flags().GetBool("all")
	MaybeDie(err)

	cacheDir := timer.GetCacheDir()

	if all {
		slog.Debug("Clearing all timers")
		err = timer.ClearAll(cacheDir)
		MaybeDie(err)
	} else {
		var name string
		if len(args) == 0 {
			name = timer.DEFAULT_TIMER_NAME
		} else {
			name = args[0]
		}

		slog.Debug("Clearing timer", "Name", name)
		err = timer.Clear(name, cacheDir)
		MaybeDie(err)
	}
}
