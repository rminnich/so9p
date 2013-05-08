package so9p

import (
	"net/rpc"
	"os"
	"syscall"
)

type Fid int64

type sFid struct {
	Node Node
}

type So9ps struct {
	Server *rpc.Server
	Path string
	Node Node
	Fs     fileFS
}

type So9pc struct {
	Client *rpc.Client
}

type Ioargs struct {
	Fid  Fid
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
	Fid  Fid
	Mode int
}

type Newargs struct {
	Name string
	Fid  Fid
	NFid Fid
	Perm os.FileMode
	Mode int
}

/* There's too much stuff we need that's too abstracted
 * in a FileInfo. Just embed the Stat_t.
 */
type FileInfo struct {
     Stat syscall.Stat_t
     Name string
}

type Nameresp struct {
	FI  FileInfo
	Fid Fid
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
