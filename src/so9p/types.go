package so9p

import (
	"net/rpc"
	"os"
	"syscall"

	"github.com/google/uuid"
)

type Fid uuid.UUID

// Server contains data for a server instance
type Server struct {
	Server   *rpc.Server
	FullPath string
	Fid      Fid
	Fs       fileFS
	Node     Node
	Files    map[Fid]Node
}

// Conn is a connection to a server from a client.
type Conn struct {
	*rpc.Client
}

// Client is a client struct for a file fid.
type Client struct {
	*Conn
	fi  FileInfo
	Fid Fid
}

// File is the client struct for a file.
type File struct {
	*Client
	Fid Fid
	Off int64
	EOF bool
}

// Ioargs are args for an IO, either read or write.
type Ioargs struct {
	Fid  Fid
	Len  int
	Off  int64
	Data []byte
}

// Ioresp is the response to an IO.
type Ioresp struct {
	Len  int
	Data []byte
	EOF  bool
}

// FIresp is the response to a stat or readdir.
type FIresp struct {
	FI []FileInfo
}

// AttachArgs is used in an attach, and has a server name type and args.
type AttachArgs struct {
	Name string
	Args []interface{}
}

// Attachresp is the resonse to an attach.
type Attachresp struct {
	FI  FileInfo
	Fid Fid
}

// NameArgs have a Name, Fid, and mode.
type NameArgs struct {
	Name string
	Fid  Fid
	Mode int
}

// NewArgs are the args for creation.
type NewArgs struct {
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

// FileInfo has a stat, name, and link value.
type FileInfo struct {
	Stat syscall.Stat_t
	Name string
	Link string
}

// Nameresp is a response for operations on names. It has a FI and Fid.
type Nameresp struct {
	FI  FileInfo
	Fid Fid
}

// FS is the interface for file servers.
type FS interface {
	Root() (Node, error)
}

// Node is the interface for an Nodde, requiring implementations for Attach and FI.
type Node interface {
	Attach(*AttachArgs, *Attachresp) error
	FI(name string) (FileInfo, error)
}
