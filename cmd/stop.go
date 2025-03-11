package cmd

import (
	"fmt"
	"log/slog"

	"github.com/dusktreader/gowatch/timer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:	"stop",
	Short:	"Stop a timer",
	Long:	"Stop a named timer",
	Args:	cobra.MaximumNArgs(1),
	Run:	stopMain,
}

func stopMain(_ *cobra.Command, args []string){
	var name string
	if len(args) == 0 {
		name = timer.DEFAULT_TIMER_NAME
	} else {
		name = args[0]
	}

	cacheDir := timer.GetCacheDir()

	t, err := timer.Load(name, cacheDir)
	MaybeDie(err)

	slog.Debug("Stopping timer", "Name", name)
	err = t.Stop()
	MaybeDie(err)

	err = t.Dump(name, cacheDir)
	MaybeDie(err)

	slog.Debug("Timer started", "Name", name, "Timer", t)
	fmt.Println(t.ElapsedString())
}
