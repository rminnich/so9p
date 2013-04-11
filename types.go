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
	Fs fileFS
}

type so9pc struct {
	Client *rpc.Client
}

type Nameargs struct {
	Name string
	Fid  fid
	NFid  fid
}

type FileInfo struct {
    SName	string
    SSize int64 
    SMode os.FileMode     
    SModTime time.Time
    SIsDir bool      
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
	FI() (FileInfo, error)
}

type fileNode struct {
	FullPath, Name string
}

