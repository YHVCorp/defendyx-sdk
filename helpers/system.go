package helpers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/elastic/go-sysinfo"
)

type OSInfo struct {
	Arch        string
	Hostname    string
	Os          string
	CurrentUser string
	Mac         string
	Addresses   string
}

func GetMyPath() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func ExecuteWithResult(c string, dir string, arg ...string) (string, bool) {
	cmd := exec.Command(c, arg...)

	cmd.Dir = dir
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}

	out, err := cmd.Output()
	if err != nil {
		return string(out[:]) + err.Error(), true
	}

	return string(out[:]), false
}

func Execute(c string, dir string, arg ...string) error {
	cmd := exec.Command(c, arg...)
	cmd.Dir = dir

	return cmd.Run()
}

func GetOsInfo() (*OSInfo, error) {
	var info OSInfo

	hostInfo, err := sysinfo.Host()
	if err != nil {
		return nil, fmt.Errorf("error getting host info: %v", err)
	}

	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	info.Arch = hostInfo.Info().Architecture
	info.Hostname = hostInfo.Info().Hostname
	info.Os = fmt.Sprintf("%s %s %s", hostInfo.Info().OS.Type, hostInfo.Info().OS.Name, hostInfo.Info().OS.Version)
	info.CurrentUser = currentUser.Username
	info.Mac = strings.Join(hostInfo.Info().MACs, ",")
	info.Addresses = strings.Join(hostInfo.Info().IPs, ",")

	return &info, nil
}
