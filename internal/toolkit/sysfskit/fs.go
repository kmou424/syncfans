package sysfskit

import (
	"github.com/kmou424/ero"
	"os"
	"path/filepath"
	"strings"
)

const maxSymlinkDepth = 10

// CheckSysfsNode checks if a given path is a valid sysfs node.
func CheckSysfsNode(path string) (bool, error) {
	return checkSysfsNodeWithDepth(path, 0)
}

func checkSysfsNodeWithDepth(path string, depth int) (bool, error) {
	if strings.TrimSpace(path) == "" {
		return false, ero.New("path is empty")
	}

	if depth > maxSymlinkDepth {
		return false, ero.New("too many levels of symbolic links")
	}

	if !checkPathExists(path) {
		return false, ero.Newf("path not exists: %s", path)
	}

	if !isSysfsNode(path) {
		if isLink(path) {
			target, err := getLinkTarget(path)
			if err != nil {
				return false, ero.Newf("failed to get link target: %v", err)
			}

			if !filepath.IsAbs(target) {
				target = filepath.Join(filepath.Dir(path), target)
			}

			return checkSysfsNodeWithDepth(target, depth+1)
		} else {
			return false, ero.Newf("path is not a sysfs node: %s", path)
		}
	}

	if !isAccessible(path) {
		return false, ero.Newf("sysfs node is not accessible: %s", path)
	}

	return true, nil
}

func checkPathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func isRegularFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

func isLink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink == os.ModeSymlink
}

func getLinkTarget(path string) (string, error) {
	return os.Readlink(path)
}

func isSysfsNode(path string) bool {
	if !strings.HasPrefix(path, "/sys/") {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	return strings.HasPrefix(absPath, "/sys/")
}

func isAccessible(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		if info, statErr := os.Stat(path); statErr == nil && info.IsDir() {
			_, err = os.ReadDir(path)
			return err == nil
		}
		return false
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	return true
}
