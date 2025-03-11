package timer_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/dusktreader/gowatch/timer"
)

func home() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}

func span(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic("Couldn't parse duration: " + s)
	}
	return d
}

// func boom(t *testing.T, msgs ...interface{}) {
// 	if len(msgs) > 0 {
// 		t.Fatalf("BOOM: %v", fmt.Sprintf(msgs[0].(string), msgs[1:]...))
// 	} else {
// 		t.Fatalf("BOOM!")
// 	}
// }

func moment(m string) time.Time {
	t, err := time.Parse(time.RFC3339, m)
	if err != nil {
		panic("Couldn't parse time: " + m)
	}
	return t
}

func TestGetConfigDir(t *testing.T) {
	want := filepath.Join(home(), "/.config", timer.APP_NAME)
	got := timer.GetConfigDir()
	if got != want {
		t.Errorf("Wrong config dir: wanted %v, got %v", want, got)
	}
}

func TestGetCacheDir(t *testing.T) {
	want := filepath.Join(home(), "/.cache", timer.APP_NAME)
	got := timer.GetCacheDir()
	if got != want {
		t.Errorf("Wrong config dir: wanted %v, got %v", want, got)
	}
}

func TestClear_NoFile(t *testing.T) {
	cacheDir := t.TempDir()

	err := timer.Clear("valid", cacheDir)
	if err == nil {
		t.Fatalf("Clear didn't return error even though file doesn't exist")
	}
}

func TestClear_ValidFile(t *testing.T) {
	cacheDir := t.TempDir()

	path := filepath.Join(cacheDir, "valid.json")
	err := os.WriteFile(path, []byte("some data"), 0644)
	if err != nil {
		t.Fatalf("Couldn't write dummy data to file: %v", err)
	}

	err = timer.Clear("valid", cacheDir)
	if err != nil {
		t.Errorf("Clear returned an error when trying to clear file %v: %v", path, err)
	}

}

func TestLoad_NoFile_notExistOk(t *testing.T) {
	cacheDir := t.TempDir()
	got, err := timer.Load("nonexistent", cacheDir)
	want := new(timer.Timer)
	if err != nil {
		t.Errorf("Load returned an error on nonexistent watch: %v", err)
	} else if !reflect.DeepEqual(want, got) {
		t.Errorf("Load returned wrong timer for nonexistent watch: wanted %v, got %v", want, got)
	}
}

func TestLoad_NoFile_notExistErr(t *testing.T) {
	cacheDir := t.TempDir()
	_, err := timer.Load("nonexistent", cacheDir, true)
	if err == nil {
		t.Errorf("Load didn't error on a non existing file when mustExist was set")
	}
}

func TestLoad_InvalidFile(t *testing.T) {
	cacheDir := t.TempDir()

	path := filepath.Join(cacheDir, "invalid.json")
	err := os.WriteFile(path, []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("Couldn't write invalid watch file: %v", err)
	}

	got, err := timer.Load("invalid", cacheDir)
	if err == nil {
		t.Errorf("Load didn't return an error on invalid watch file")
	} else if got != nil {
		t.Errorf("Load returned a timer on invalid watch file: %v", got)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	cacheDir := t.TempDir()

	want := &timer.Timer{
		TotalTime:	span("3m"),
		StartTime:	moment("2025-03-11T11:02:00Z"),
		EndTime:	moment("2025-03-11T11:05:00Z"),
	}
	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("Couldn't marshal timer: %v", err)
	}

	path := filepath.Join(cacheDir, "valid.json")
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		t.Fatalf("Couldn't write valid watch file: %v", err)
	}

	got, err := timer.Load("valid", cacheDir)
	if err != nil {
		t.Errorf("Load returned an error for a valid watch file: %v", err)
	} else if !reflect.DeepEqual(got, want) {
		t.Errorf("Load returned wrong timer for valid watch: wanted %v, got %v", want, got)
	}
}

func TestLoadAll_BadDir(t *testing.T) {
	_, err := timer.LoadAll("nonexistent")
	if err == nil {
		t.Errorf("LoadAll didn't return an error on a nonexistent directory")
	}
}

func TestLoadAll_Valid(t *testing.T) {
	cacheDir := t.TempDir()

	timer1 := &timer.Timer{
		TotalTime:	span("1m"),
		StartTime:	moment("2025-03-11T17:39:00Z"),
		EndTime:	moment("2025-03-11T17:40:00Z"),
	}
	err := timer1.Dump("good1", cacheDir)
	if err != nil {
		t.Fatalf("Couldn't dump good1: %v", err)
	}

	timer2 := &timer.Timer{
		TotalTime:	span("2m"),
		StartTime:	moment("2025-03-11T17:40:00Z"),
		EndTime:	moment("2025-03-11T17:42:00Z"),
	}
	err = timer2.Dump("good2", cacheDir)
	if err != nil {
		t.Fatalf("Couldn't dump good2: %v", err)
	}

	timer3 := &timer.Timer{
		TotalTime:	span("3m"),
		StartTime:	moment("2025-03-11T17:42:00Z"),
		EndTime:	moment("2025-03-11T17:45:00Z"),
	}
	err = timer3.Dump("good3", cacheDir)
	if err != nil {
		t.Fatalf("Couldn't dump good3: %v", err)
	}

	subDir := filepath.Join(cacheDir, "sub")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Couldn't make subdir: %v", err)
	}

	timer4 := &timer.Timer{
		TotalTime:	span("4m"),
		StartTime:	moment("2025-03-11T17:45:00Z"),
		EndTime:	moment("2025-03-11T17:49:00Z"),
	}
	err = timer4.Dump("good4", subDir)
	if err != nil {
		t.Fatalf("Couldn't dump good4: %v", err)
	}

	other := filepath.Join(cacheDir, "other.txt")
	err = os.WriteFile(other, []byte("excluded file"), 0644)
	if err != nil {
		t.Fatalf("Couldn't write other file: %v", err)
	}

	want := []*timer.NamedTimer{
		{
			Name: "good1",
			Ticks: timer1,
		},
		{
			Name: "good2",
			Ticks: timer2,
		},
		{
			Name: "good3",
			Ticks: timer3,
		},
	}

	got, err := timer.LoadAll(cacheDir)
	if err != nil {
		t.Fatalf("LoadAll returned an error: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("LoadAll didn't return the right timers:\n\nwanted\n%v,\n\ngot\n%v", want, got)
	}
}

func TestDump_WriteFileError(t *testing.T) {
	cacheDir := t.TempDir()

	dummyData := []byte("dummy data")
	path := filepath.Join(cacheDir, "forbidden.json")
	err := os.WriteFile(path, dummyData, 0000)
	if err != nil {
		t.Fatalf("Couldn't write dummy watch file: %v", err)
	}

	ticks := new(timer.Timer)
	err = ticks.Dump("forbidden", cacheDir)
	if err == nil {
		t.Errorf("Dump watch didn't return an error on a protected file")
	}
}

func TestDump_ValidFile(t *testing.T) {
	cacheDir := t.TempDir()

	want := &timer.Timer{
		TotalTime:	span("3m"),
		StartTime:	moment("2025-03-11T11:02:00Z"),
		EndTime:	moment("2025-03-11T11:05:00Z"),
	}

	err := want.Dump("valid", cacheDir)
	if err != nil {
		t.Errorf("Dump returned an error for a valid watch file: %v", err)
	}

	path := filepath.Join(cacheDir, "valid.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Couldn't read data from valid watch file: %v", err)
	}

	got := new(timer.Timer)
	err = json.Unmarshal(data, got)
	if err != nil {
		t.Fatalf("Couldn't unmarshal data from watch file: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Dump didn't dump data correctly: wanted %v, got %v", want, got)
	}
}

func TestIsRunning(t *testing.T) {
	ticks := new(timer.Timer)

	if ticks.IsRunning() {
		t.Errorf("IsRunning returned true on a timer with no start time")
	}

	ticks.StartTime = moment("2025-03-11T11:47:00Z")
	if !ticks.IsRunning() {
		t.Errorf("IsRunning returned false on a timer with only a start time")
	}

	ticks.EndTime = moment("2025-03-11T11:49:00Z")
	if ticks.IsRunning() {
		t.Errorf("IsRunning returned true on a timer with a start time and end time")
	}
}

func freeze(t *testing.T, moment string) timer.FixedNowProvider {
	tt, err := time.Parse(time.RFC3339, moment)
	if err != nil {
		t.Fatalf("Couldn't parse time: %v", err)
	}
	return timer.FixedNowProvider{
		Moment: tt,
	}
}

func TestStart(t *testing.T) {
	ticks := new(timer.Timer)

	np := freeze(t, "2025-03-11T12:05:00Z")

	err := ticks.Start(np)
	if err != nil {
		t.Errorf("Start returned an error even though timer is not running: %v", err)
	}

	got := ticks.StartTime
	want := np.Moment
	if want != got {
		t.Errorf("Start didn't set the start time correctly: wanted %v, got %v", want, got)
	}

	err = ticks.Start(np)
	if err == nil {
		t.Errorf("Start didn't return an error even though timer is running")
	}
}

func TestStop(t *testing.T) {
	ticks := &timer.Timer{
		TotalTime: span("3m"),
		StartTime: moment("2025-03-11T12:27:00Z"),
	}

	np := freeze(t, "2025-03-11T12:29:00Z")

	err := ticks.Stop(np)
	if err != nil {
		t.Errorf("Stop returned an error even though timer is running: %v", err)
	}

	got := ticks.EndTime
	want := np.Moment
	if want != got {
		t.Errorf("Stop didn't set the end time correctly: wanted %v, got %v", want, got)
	}

	wantTotal := span("5m")
	gotTotal := ticks.TotalTime
	if wantTotal != gotTotal {
		t.Errorf("Stop didn't update total time correctly: wanted %v, got %v", want, got)
	}

	err = ticks.Stop(np)
	if err == nil {
		t.Errorf("Stop didn't return an error even though timer is not running")
	}
}

func TestToggle(t *testing.T) {
	ticks := new(timer.Timer)

	np := freeze(t, "2025-03-11T12:42:00Z")
	m1 := np.Moment
	st := ticks.Toggle(np)
	if st {
		t.Errorf("Toggle reported a stoppage on a timer that wasn't running")
	}
	want := &timer.Timer{
		StartTime: m1,
	}
	if !reflect.DeepEqual(want, ticks) {
		t.Errorf("Toggle didn't update timer correctly: wanted %v, got %v", want, ticks)
	}

	np = freeze(t, "2025-03-11T12:47:00Z")
	m2 := np.Moment
	st = ticks.Toggle(np)
	if !st {
		t.Errorf("Toggle reported no stoppage on a timer that was running")
	}
	want = &timer.Timer{
		StartTime: m1,
		EndTime: m2,
		TotalTime: span("5m"),
	}
	if !reflect.DeepEqual(want, ticks) {
		t.Errorf("Toggle didn't update timer correctly: wanted %v, got %v", want, ticks)
	}

	np = freeze(t, "2025-03-11T12:56:00Z")
	m3 := np.Moment
	st = ticks.Toggle(np)
	if st {
		t.Errorf("Toggle reported a stoppage on a timer that wasn't running")
	}
	want = &timer.Timer{
		StartTime: m3,
		TotalTime: span("5m"),
	}
	if !reflect.DeepEqual(want, ticks) {
		t.Errorf("Toggle didn't update timer correctly: wanted %v, got %v", want, ticks)
	}

	np = freeze(t, "2025-03-11T12:57:00Z")
	m4 := np.Moment
	st = ticks.Toggle(np)
	if !st {
		t.Errorf("Toggle reported no stoppage on a timer that was running")
	}
	want = &timer.Timer{
		StartTime: m3,
		EndTime: m4,
		TotalTime: span("6m"),
	}
	if !reflect.DeepEqual(want, ticks) {
		t.Errorf("Toggle didn't update timer correctly: wanted %v, got %v", want, ticks)
	}
}

func TestReset(t *testing.T) {
	ticks := &timer.Timer{
		TotalTime:	span("3m"),
		StartTime:	moment("2025-03-11T11:02:00Z"),
		EndTime:	moment("2025-03-11T11:05:00Z"),
	}
	ticks.Reset()

	want := new(timer.Timer)
	if !reflect.DeepEqual(want, ticks) {
		t.Errorf("Reset didn't clear timer correctly: wanted %v, got %v", want, ticks)
	}
}
