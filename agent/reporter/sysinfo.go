package reporter

import (
	"bytes"
	"encoding/json"
	"github.com/kmou424/ero"
	"github.com/kmou424/syncfans/internal/caused"
	"github.com/kmou424/syncfans/internal/conf"
	"github.com/kmou424/syncfans/internal/proto"
	"github.com/spf13/cast"
	"os"
	"os/exec"
	"runtime"
)

type SysInfo map[string]any

func (s SysInfo) Report() (*proto.ReportSysInfo, error) {
	dst := &proto.ReportSysInfo{}
	marshal, err := json.Marshal(s)
	if err != nil {
		return nil, caused.ValueError(ero.Wrap(err, "failed to marshal sysinfo"))
	}
	err = json.Unmarshal(marshal, dst)
	if err != nil {
		return nil, caused.ValueError(ero.Wrap(err, "failed to unmarshal sysinfo"))
	}
	return dst, nil
}

func GetSysInfo() (SysInfo, error) {
	sysInfo := SysInfo{}

	config := conf.GetAgentConfig()
	for entry, base := range config.Sysinfo {
		value, err := querySysInfo(base)
		switch base.Type {
		case "string":
			sysInfo[entry] = value
		case "int":
			sysInfo[entry] = cast.ToInt(value)
		case "float":
			sysInfo[entry] = cast.ToFloat64(value)
		default:
			return nil, caused.ValueError(ero.New("unsupported sysinfo type"))
		}
		if err != nil {
			return nil, ero.Wrap(err, "failed to query sysinfo entry: %s", entry)
		}
	}

	return sysInfo, nil
}

func querySysInfo(base *conf.SysInfoBase) (string, error) {
	query := base.Query
	switch base.Method {
	case "file":
		file, err := os.ReadFile(query)
		if err != nil {
			return "", caused.FileSystemError(ero.Wrap(err, "failed to read file"))
		}
		return string(file), nil
	case "shell":
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", query)
		} else {
			// For Unix-like systems (Linux, macOS)
			cmd = exec.Command("sh", "-c", query)
		}
		out, err := cmd.Output()
		out = bytes.TrimSpace(out)
		if err != nil {
			return "", caused.RuntimeError(ero.Wrap(err, "failed to execute command"))
		}
		return string(out), nil
	default:
		return "", caused.RuntimeError(ero.New("unsupported sysinfo query method"))
	}
}
