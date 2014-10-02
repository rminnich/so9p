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
	Path   string
	Fs     fileFS
}

type So9pConn struct {
	*rpc.Client
}

type So9pc struct {
	*So9pConn
	fi  FileInfo
	Fid Fid
}

type So9file struct {
	*So9pc
	Fid Fid
	Off int64
	EOF bool
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
	EOF  bool
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
 * in a FileInfo.
 * Toss in the symlink. There are lots of them in Linux,
 * and saving an RPC is always a nice idea. It's super
 * cheap just to do it.
 */
type FileInfo struct {
	Stat syscall.Stat_t
	Name string
	Link string
}

type Nameresp struct {
	FI  FileInfo
	Fid Fid
}

type FS interface {
	Root() (Node, error)
}

type Node interface {
	FI(name string) (FileInfo, error)
}
