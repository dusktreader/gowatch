package cmd

import (
	"log/slog"

	"github.com/dusktreader/gowatch/timer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:	"start",
	Short:	"Start a timer",
	Long:	"Start a named timer",
	Args:	cobra.MaximumNArgs(1),
	Run:	startMain,
}

func startMain(_ *cobra.Command, args []string){
	var name string
	if len(args) == 0 {
		name = timer.DEFAULT_TIMER_NAME
	} else {
		name = args[0]
	}

	cacheDir := timer.GetCacheDir()

	t, err := timer.Load(name, cacheDir)
	MaybeDie(err)

	slog.Debug("Starting timer", "Name", name)
	err = t.Start()
	MaybeDie(err)

	err = t.Dump(name, cacheDir)
	MaybeDie(err)

	slog.Debug("Timer started", "Name", name, "Timer", t)
}
