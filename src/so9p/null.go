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

func (n *nullNode) Attach(string) (Node, error) {
	return nil, syscall.EPERM
}

func (n *nullNode) UnAttach(Fid) error {
	return syscall.EPERM
}

func (n *nullNode) FI(name string) (FileInfo, error) {
	return FileInfo{}, nil
}

func (n *nullNode) Create(string, int, os.FileMode) (Node, error) {
	return nil, syscall.EBADFD
}

func (n *nullNode) Mkdir(string, int, os.FileMode) error {
	return syscall.EBADFD
}

func (n *nullNode) Read(int, int64) ([]byte, error) {
	return nil, syscall.EBADFD
}

func (n *nullNode) Write([]byte, int64) (int, error) {
	return -1, syscall.EBADFD
}

func (n *nullNode) Close() error {
	return syscall.EBADFD
}

func (n *nullNode) ReadDir(string) ([]FileInfo, error) {
	return nil, syscall.EBADFD
}

func (n *nullNode) Readlink(string) (string, error) {
	return "", syscall.EBADFD
}
