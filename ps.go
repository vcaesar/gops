// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.
//

package ps

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

// Nps process struct
type Nps struct {
	// Pid  int32
	Pid  int
	Name string
}

// ToInt convert []int32 to []int
func ToInt(pid []int32) (res []int) {
	for _, v := range pid {
		res = append(res, int(v))
	}

	return
}

// GetPid get the process id
func GetPid() int {
	return os.Getpid()
}

// Pids get the all process id
func Pids() ([]int, error) {
	ids, err := process.Pids()
	return ToInt(ids), err
}

// PidExists determine whether the process exists
func PidExists(pid int) (bool, error) {
	return process.PidExists(int32(pid))
}

// Process get the all process struct
func Process() ([]Nps, error) {
	var npsArr []Nps
	pid, err := process.Pids()
	if err != nil {
		return npsArr, err
	}

	for i := 0; i < len(pid); i++ {
		nps, _ := process.NewProcess(pid[i])
		names, _ := nps.Name()

		np := Nps{
			int(pid[i]),
			names,
		}

		npsArr = append(npsArr, np)
	}

	return npsArr, err
}

// FindName find the process name by the process id
func FindName(pid int) (string, error) {
	nps, err := process.NewProcess(int32(pid))
	if err != nil {
		return "", err
	}

	return nps.Name()
}

// FindNames find the all process name
func FindNames() ([]string, error) {
	var strArr []string
	pid, err := process.Pids()

	if err != nil {
		return strArr, err
	}

	for i := 0; i < len(pid); i++ {
		nps, _ := process.NewProcess(pid[i])
		names, _ := nps.Name()

		strArr = append(strArr, names)
	}

	return strArr, err
}

// FindMainIds finds all main process PIDs by matching the executable path.
// It looks for processes whose path contains "name" (case insensitive)
// and are main binaries, excluding helper processes:
//   - macOS: the binary is under .app/Contents/MacOS/
//   - Windows, Linux and others: the parent process runs a different
//     executable (children of multi-process apps share the parent's binary)
//
// Returns all matching main PIDs, or all matching PIDs if no main one is found.
func FindMainIds(name string) ([]int, error) {
	pids, err := Pids()
	if err != nil {
		return nil, err
	}

	name = strings.ToLower(name)
	var mainPids []int
	var fallbackPids []int

	for _, pid := range pids {
		path, err := FindPath(pid)
		if err != nil {
			continue
		}

		pathLower := strings.ToLower(path)
		if !strings.Contains(pathLower, name) {
			continue
		}

		parentPath := ""
		if runtime.GOOS != "darwin" {
			parentPath = findParentPath(pid)
		}

		if isMainBinary(runtime.GOOS, path, parentPath) {
			mainPids = append(mainPids, pid)
		} else {
			fallbackPids = append(fallbackPids, pid)
		}
	}

	if len(mainPids) > 0 {
		return mainPids, nil
	}
	return fallbackPids, nil
}

// isMainBinary reports whether path is a main app binary (not a helper).
// parentPath is the executable path of the parent process, used off macOS.
func isMainBinary(goos, path, parentPath string) bool {
	pathLower := strings.ToLower(path)
	if strings.Contains(pathLower, "helper") {
		return false
	}

	if goos == "darwin" {
		return strings.Contains(pathLower, ".app/contents/macos/")
	}

	if goos == "windows" {
		return !strings.EqualFold(path, parentPath)
	}
	return path != parentPath
}

// findParentPath returns the executable path of the parent process,
// or "" if it cannot be determined.
func findParentPath(pid int) string {
	nps, err := process.NewProcess(int32(pid))
	if err != nil {
		return ""
	}

	ppid, err := nps.Ppid()
	if err != nil {
		return ""
	}

	path, err := FindPath(int(ppid))
	if err != nil {
		return ""
	}
	return path
}

// FindId finds the main process by matching the executable path.
// It looks for a process whose path contains "name" (case insensitive)
// and prioritizes main binaries (see FindMainIds), excluding helper processes.
// Returns the PID or -1 if not found.
func FindId(name string) (int, error) {
	pids, err := FindMainIds(name)
	if err != nil {
		return -1, err
	}

	if len(pids) == 0 {
		return -1, nil
	}
	return pids[0], nil
}

// FindIds finds the all processes named with a subset
// of "name" (case insensitive),
// return matched IDs.
func FindIds(name string) ([]int, error) {
	var pids []int
	nps, err := Process()
	if err != nil {
		return pids, err
	}

	name = strings.ToLower(name)
	for i := 0; i < len(nps); i++ {
		psname := strings.ToLower(nps[i].Name)
		abool := strings.Contains(psname, name)
		if abool {
			pids = append(pids, nps[i].Pid)
		}
	}

	return pids, err
}

// FindPath find the process path by the process pid
func FindPath(pid int) (string, error) {
	nps, err := process.NewProcess(int32(pid))
	if err != nil {
		return "", err
	}

	return nps.Exe()
}

// Run command shell
func Run(path string) ([]byte, error) {
	cmdName := "/bin/bash"
	params := "-c"
	if runtime.GOOS == "windows" {
		cmdName = "cmd"
		params = "/c"
	}

	cmd := exec.Command(cmdName, params, path)
	output, err := cmd.Output()

	return output, err
}

// IsRun return the process is runing or not
func IsRun(pid int) (bool, error) {
	nps, err := process.NewProcess(int32(pid))
	if err != nil {
		return false, err
	}

	return nps.IsRunning()
}

// Status return the process status
func Status(pid int) ([]string, error) {
	nps, err := process.NewProcess(int32(pid))
	if err != nil {
		return []string{}, err
	}

	return nps.Status()
}

// Kill kill the process by PID
func Kill(pid int) error {
	p := os.Process{Pid: pid}
	return p.Kill()
}
