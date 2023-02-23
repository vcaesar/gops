package ps

import (
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

func TestFindNames(t *testing.T) {
	name, err := FindNames()
	tt.NotZero(t, len(name))
	tt.Nil(t, err)
}
