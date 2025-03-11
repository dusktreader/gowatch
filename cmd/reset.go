package cmd

import (
	"log/slog"

	"github.com/dusktreader/gowatch/timer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:	"reset",
	Short:	"Reset a timer",
	Long:	"Reset a named timer",
	Args:	cobra.MaximumNArgs(1),
	Run:	resetMain,
}

func resetMain(_ *cobra.Command, args []string){
	var name string
	if len(args) == 0 {
		name = timer.DEFAULT_TIMER_NAME
	} else {
		name = args[0]
	}

	cacheDir := timer.GetCacheDir()

	t, err := timer.Load(name, cacheDir)
	MaybeDie(err)

	slog.Debug("Resetting timer", "Name", name)
	t.Reset()

	err = t.Dump(name, cacheDir)
	MaybeDie(err)
}
