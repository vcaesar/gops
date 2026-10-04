package ps

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/vcaesar/tt"
)

func TestGetPid(t *testing.T) {
	i := GetPid()
	tt.NotZero(t, i)
}

func TestPids(t *testing.T) {
	ids, err := Pids()
	tt.NotZero(t, len(ids))
	tt.Nil(t, err)
}

func TestProcess(t *testing.T) {
	ps, err := Process()
	tt.NotZero(t, len(ps))
	tt.Nil(t, err)
}

func TestIsMainBinary(t *testing.T) {
	cases := []struct {
		goos, path, parent string
		want               bool
	}{
		{"darwin", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "/sbin/launchd", true},
		{"darwin", "/Applications/Google Chrome.app/Contents/Frameworks/Google Chrome Helper.app/Contents/MacOS/Google Chrome Helper", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", false},
		{"darwin", "/usr/local/bin/chrome", "/bin/zsh", false},
		{"windows", `C:\Program Files\Google\Chrome\Application\chrome.exe`, `C:\Windows\explorer.exe`, true},
		{"windows", `C:\Program Files\Google\Chrome\Application\chrome.exe`, `c:\program files\google\chrome\application\CHROME.EXE`, false},
		{"windows", `C:\Program Files\App\apphelper.exe`, `C:\Windows\explorer.exe`, false},
		{"windows", `C:\Program Files\App\app.exe`, "", true},
		{"linux", "/opt/google/chrome/chrome", "/usr/lib/systemd/systemd", true},
		{"linux", "/opt/google/chrome/chrome", "/opt/google/chrome/chrome", false},
		{"linux", "/usr/lib/app/app-helper", "/bin/bash", false},
		{"freebsd", "/usr/local/bin/app", "/bin/sh", true},
	}

	for _, c := range cases {
		tt.Equal(t, c.want, isMainBinary(c.goos, c.path, c.parent), c.goos+" "+c.path)
	}
}

func TestFindMainIds(t *testing.T) {
	path, err := FindPath(GetPid())
	tt.Nil(t, err)

	ids, err := FindMainIds(filepath.Base(path))
	tt.Nil(t, err)
	tt.True(t, slices.Contains(ids, GetPid()))

	ids, err = FindMainIds("no-such-process-name-xyz")
	tt.Nil(t, err)
	tt.Zero(t, len(ids))
}

func TestFindNames(t *testing.T) {
	name, err := FindNames()
	tt.NotZero(t, len(name))
	tt.Nil(t, err)
}
