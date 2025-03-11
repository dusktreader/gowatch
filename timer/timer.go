package timer

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const APP_NAME = "gowatch"
const DEFAULT_TIMER_NAME = "default"

type Timer struct {
	TotalTime	time.Duration	`json:"total"`
	StartTime	time.Time		`json:"start"`
	EndTime		time.Time		`json:"end"`
}

type NamedTimer struct {
	Name	string
	Ticks	*Timer
}

type NowProvider interface {
	Now() time.Time
}

type RealNowProvider struct {}

func (RealNowProvider) Now() time.Time {
	return time.Now()
}

type FixedNowProvider struct {
	Moment time.Time
}

func (p FixedNowProvider) Now() time.Time {
	return p.Moment
}

func now(nowProviderArg []NowProvider) time.Time {
	var nowProvider NowProvider
	if len(nowProviderArg) == 0 {
		nowProvider = RealNowProvider{}
	} else {
		nowProvider = nowProviderArg[0]
	}
	return nowProvider.Now()
}

func GetConfigDir() string {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Sprintf("Error getting user config dir: %v", err))
	}

	configDir := filepath.Join(baseDir, APP_NAME)
	return configDir
}

func GetCacheDir() string {
	baseDir, err := os.UserCacheDir()
	if err != nil {
		panic(fmt.Sprintf("Error getting user cache dir: %v", err))
	}

	cacheDir := filepath.Join(baseDir, APP_NAME)
	return cacheDir
}

func EnsureDir(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("Couldn't create directory: %v", err)
	}
	return nil
}

func (t *Timer) String() string {
	return fmt.Sprintf(
		"(%s -- %s) -> %s",
		t.StartTime.Format(time.RFC3339),
		t.EndTime.Format(time.RFC3339),
		t.ElapsedString(),
	)
}

func (nt *NamedTimer) String() string {
	return fmt.Sprintf(
		"%s: %s\n",
		nt.Name,
		nt.Ticks,
	)
}

func (t *Timer) Elapsed(nowProviderArg ...NowProvider) time.Duration {
	if !t.StartTime.IsZero() && t.EndTime.IsZero() {
		endTime := now(nowProviderArg)
		totalTime := endTime.Sub(t.StartTime)
		return totalTime
	}

	return t.TotalTime
}

func (t *Timer) ElapsedString(nowProviderArg ...NowProvider) string {
	elapsed := t.Elapsed(nowProviderArg...)
	return elapsed.Round(time.Millisecond).String()
}

func allTimerFiles(cacheDir string) ([]os.DirEntry, error) {
	allFiles, err := os.ReadDir(cacheDir)
	if err != nil {
		msg := "Error finding files"
		slog.Error(msg, "error", err)
		return nil, fmt.Errorf(msg + ": %v", err)
	}

	matchFiles := make([]os.DirEntry, 0)
	for _, file := range allFiles {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		ext := filepath.Ext(filename)
		if ext != ".json" {
			continue
		}
		matchFiles = append(matchFiles, file)
	}
	return matchFiles, nil
}

func Clear(name string, cacheDir string) error {
	path := filepath.Join(cacheDir, name + ".json")
	slog.Debug("Clearing timer file", "path", path)

	err := os.Remove(path)
	if err != nil {
		msg := "Error clearing timer data"
		slog.Error(msg, "error", err)
		return fmt.Errorf(msg + ": %v", err)
	}
	return nil
}

func ClearAll(cacheDir string) error {
	allFiles, err := allTimerFiles(cacheDir)
	if err != nil {
		return err
	}

	failures := make([]os.DirEntry, 0)
	for _, file := range allFiles {
		path := filepath.Join(cacheDir, file.Name())
		err := os.Remove(path)
		if err != nil {
			slog.Error("Couldn't remove timer file!", "path", path, "error", err)
			failures = append(failures, file)
		}
	}

	if len(failures) > 0 {
		return fmt.Errorf("Some files couldn't be removed: %v", failures)
	}

	return nil
}

func Load(name string, cacheDir string, mustExist ...bool) (*Timer, error) {
	path := filepath.Join(cacheDir, name + ".json")
	slog.Debug("Loading timer from file", "path", path)

	t := new(Timer)

	data, err := os.ReadFile(path)
	if err != nil {
		if len(mustExist) > 0 && mustExist[0] {
			msg := "File does not exist"
			slog.Error(msg, "path", path)
			return nil, fmt.Errorf(msg + ": %v", path)
		}
		return t, nil
	}

	slog.Debug("Deserializing data")
	err = json.Unmarshal(data, t)
	if err != nil {
		msg := "Error loading timer data"
		slog.Error(msg, "error", err)
		return nil, fmt.Errorf(msg + ": %v", err)
	}

	return t, nil
}

func LoadAll(cacheDir string) ([]*NamedTimer, error) {
	allFiles, err := allTimerFiles(cacheDir)
	if err != nil {
		return nil, err
	}

	namedTimers := make([]*NamedTimer, 0)
	for _, file := range allFiles {
		filename := file.Name()
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)

		ticks, err := Load(name, cacheDir, true)
		if err != nil {
			slog.Warn("Skipping timer that failed to load", "name", name, "error", err)
			continue
		}

		namedTimers = append(
			namedTimers,
			&NamedTimer{
				Name: name,
				Ticks: ticks,
			},
		)
	}

	return namedTimers, nil
}

func (t *Timer) Dump(name string, cacheDir string) error {
	path := filepath.Join(cacheDir, name + ".json")

	slog.Debug("Serializing data")
	data, err := json.Marshal(t)
	if err != nil {
		msg := "Error dumping timer data"
		slog.Error(msg, "error", err)
		return fmt.Errorf(msg + ": %v", err)
	}

	slog.Debug("Dumping timer to file", "path", path)
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		msg := "Error writing timer data to file"
		slog.Error(msg, "error", err)
		return fmt.Errorf(msg + ": %v", err)
	}

	return nil
}

func (t *Timer) IsRunning() bool {
	return !t.StartTime.IsZero() && t.EndTime.IsZero()
}

func (t *Timer) Start(nowProviderArg ...NowProvider) error {
	if t.IsRunning() {
		return fmt.Errorf("Timer is already running")
	}

	t.StartTime = now(nowProviderArg)
	t.EndTime = time.Time{}
	return nil
}

func (t *Timer) Stop(nowProviderArg ...NowProvider) error {
	if !t.IsRunning() {
		return fmt.Errorf("Timer is not running")
	}

	t.EndTime = now(nowProviderArg)
	t.TotalTime += t.EndTime.Sub(t.StartTime)
	return nil
}

func (t *Timer) Toggle(nowProviderArg ...NowProvider) bool {
	wasStopped := false
	if t.IsRunning() {
		err := t.Stop(nowProviderArg...)
		if err != nil {
			panic("Error stopping timer (this should not happen): " + err.Error())
		}
		wasStopped = true
	} else {
		err := t.Start(nowProviderArg...)
		if err != nil {
			panic("Error starting timer (this should not happen): " + err.Error())
		}
	}
	return wasStopped
}

func (t *Timer) Reset() {
	t.StartTime = time.Time{}
	t.EndTime = time.Time{}
	t.TotalTime = 0
}
