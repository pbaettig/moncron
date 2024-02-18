package model

import (
	"bufio"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type OperatingSystem struct {
	Name          string
	Version       string
	KernelVersion string
}

type Host struct {
	Name string
	OS   OperatingSystem
}

func unquote(s string) string {
	if strings.Count(s, "\"") == 2 {
		uq, _ := strconv.Unquote(s)
		return uq
	}

	return s
}

func parseOSRelease() (OperatingSystem, error) {
	var (
		o   OperatingSystem = OperatingSystem{}
		err error
	)
	fd, err := os.Open("/etc/os-release")
	if err != nil {
		return o, err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		ls := strings.SplitN(scanner.Text(), "=", 2)
		if ls[0] == "NAME" {
			o.Name = unquote(ls[1])
			continue
		}

		if ls[0] == "VERSION_CODENAME" {
			o.Version = unquote(ls[1])
			continue
		}

	}
	return o, nil
}

func runUname() (string, error) {
	cmd := exec.Command("uname", "-r")
	stdout := new(strings.Builder)
	cmd.Stdout = stdout
	if err := cmd.Run(); err != nil {
		return "", nil
	}

	return strings.TrimSpace(stdout.String()), nil
}

func GetHost() (Host, error) {
	var (
		info Host = Host{}
		err  error
	)
	info.Name, err = os.Hostname()
	if err != nil {
		return info, err
	}

	info.OS, err = parseOSRelease()
	if err != nil {
		return info, err
	}

	info.OS.KernelVersion, err = runUname()
	return info, err
}
