package so9p

// the null server does nothing, badly.

import (
	"os"
	"syscall"
)

type nullFS struct {
}

type nullNode struct {
}

var (
	null = &nullNode{}
)

func (node *nullNode) FI(name string) (FileInfo, error) {
	return FileInfo{}, nil
}

func (node *nullNode) Create(string, int, os.FileMode) (Node, error) {
	return nil, syscall.EBADFD
}

func (node *nullNode) Mkdir(string, int, os.FileMode) error {
	return syscall.EBADFD
}

func (node *nullNode) Read(int, int64) ([]byte, error) {
	return nil, syscall.EBADFD
}

func (node *nullNode) Write([]byte, int64) (int, error) {
	return -1, syscall.EBADFD
}

func (node *nullNode) Close() error {
	return syscall.EBADFD
}

func (node *nullNode) ReadDir(string) ([]FileInfo, error) {
	return nil, syscall.EBADFD
}

func (node *nullNode) Readlink(string) (string, error) {
	return "", syscall.EBADFD
}
