package main

import (
	"net/rpc"
	"os"
	"time"
)

type fid int64

type sfid struct {
	Node Node
}

type So9ps struct {
	Server *rpc.Server
	Path string
	Node Node
	Fs     fileFS
}

type so9pc struct {
	Client *rpc.Client
}

type Ioargs struct {
	Fid  fid
	Len  int
	Off  int64
	Data []byte
}

type Ioresp struct {
	Len  int
	Data []byte
}
type FIresp struct {
	FI []FileInfo
}

type Nameargs struct {
	Name string
	Fid  fid
}

type Newargs struct {
	Name string
	Fid  fid
	NFid fid
	perm os.FileMode
	mode int
}

type FileInfo struct {
	SSize     int64
	SMode     os.FileMode
	SModTime  time.Time
	SIsDir    bool
}

type Nameresp struct {
	FI  FileInfo
	Fid fid
}

type FS interface {
	Root() (Node, error)
}

type fileFS struct {
	fileNode
}

type Node interface {
	FI(name string) (FileInfo, error)
}

type fileNode struct {
	File           *os.File
}
