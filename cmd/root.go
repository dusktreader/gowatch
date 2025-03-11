package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/dusktreader/gowatch/timer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show verbose logging output")
}

var rootCmd = &cobra.Command{
	Use:				"gowatch",
	Short:				"Go Stopwatch",
	Long:				"A command line stopwatch written in Go",
	PersistentPreRun:	preRun,
	Run:				rootMain,
}

func preRun(cmd *cobra.Command, args []string) {
	err := timer.EnsureDir(timer.GetCacheDir())
	MaybeDie(err)

	err = timer.EnsureDir(timer.GetConfigDir())
	MaybeDie(err)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	verbose, err := cmd.Flags().GetBool("verbose")
	MaybeDie(err)
	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
}

func rootMain(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}

func MaybeDie(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "There was an error:", err)
		os.Exit(1)
	}
}

func Die(msg string, flags ...interface{}) {
	msg = fmt.Sprintf(msg, flags...)
	fmt.Fprintln(os.Stderr, "Aborting:", msg)
	os.Exit(1)
}

func Execute() {
	err := rootCmd.Execute()
	MaybeDie(err)
}
