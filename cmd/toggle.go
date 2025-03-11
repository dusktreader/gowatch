package cmd

import (
	"fmt"
	"log/slog"

	"github.com/dusktreader/gowatch/timer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(toggleCmd)
}

var toggleCmd = &cobra.Command{
	Use:	"toggle",
	Short:	"Toggle a timer",
	Long:	"Toggle a named timer",
	Args:	cobra.MaximumNArgs(1),
	Run:	toggleMain,
}

func toggleMain(_ *cobra.Command, args []string){
	var name string
	if len(args) == 0 {
		name = timer.DEFAULT_TIMER_NAME
	} else {
		name = args[0]
	}

	cacheDir := timer.GetCacheDir()

	t, err := timer.Load(name, cacheDir)
	MaybeDie(err)

	slog.Debug("Toggling timer", "Name", name)
	wasStopped := t.Toggle()
	slog.Debug("Timer toggled", "Name", name, "Timer", t)

	err = t.Dump(name, cacheDir)
	MaybeDie(err)

	if wasStopped {
		fmt.Println(t.ElapsedString())
	}
}
