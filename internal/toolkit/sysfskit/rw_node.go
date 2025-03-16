package sysfskit

import (
	"github.com/kmou424/ero"
	"os"
	"strconv"
	"strings"
)

func Read(path string) (string, error) {
	if ok, err := CheckSysfsNode(path); !ok {
		return "", ero.Newf("path %s does not exist: %v", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", ero.Newf("Read failed for %s: %v", path, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func Write(path, value string) error {
	if ok, err := CheckSysfsNode(path); !ok {
		return ero.Newf("path %s does not exist: %v", path, err)
	}

	if !isRegularFile(path) {
		return ero.Newf("path %s is not a regular file", path)
	}

	err := os.WriteFile(path, []byte(value), 0644)
	if err != nil {
		return ero.Newf("write failed for %s: %v", path, err)
	}
	return nil
}

func ReadInt(path string) (int, error) {
	s, err := Read(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

func WriteInt(path string, value int) error {
	return Write(path, strconv.Itoa(value))
}

func ReadAs[T any](path string, parseFunc func(string) (T, error)) (T, error) {
	s, err := Read(path)
	if err != nil {
		var zero T
		return zero, err
	}
	return parseFunc(s)
}

type SysfsNode struct {
	path string
}

func NewSysfsNode(path string) *SysfsNode {
	return &SysfsNode{path: path}
}

func (n *SysfsNode) Read() (string, error) {
	return Read(n.path)
}

func (n *SysfsNode) Write(value string) error {
	return Write(n.path, value)
}

func (n *SysfsNode) ReadInt() (int, error) {
	return ReadInt(n.path)
}

func (n *SysfsNode) WriteInt(value int) error {
	return WriteInt(n.path, value)
}

func (n *SysfsNode) ReadAs(parseFunc func(string) (interface{}, error)) (interface{}, error) {
	s, err := n.Read()
	if err != nil {
		return nil, err
	}
	return parseFunc(s)
}
