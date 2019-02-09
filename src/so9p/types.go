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
	Fid    Fid
	Server *rpc.Server
	Fs     fileFS
}

// Conn is a connection to a server from a client.
type Conn struct {
	*rpc.Client
}

// Client is a client struct for a file fid.
// It's not clear we should have embedded the Conn.
type ClientConn struct {
	*Conn
	FI  FileInfo
	Fid Fid
}

// File is the client struct for a file.
// cureent ops are
// ReadAt, WriteAt, Read, Write, Statx
type File struct {
	*ClientConn
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
	Args []string
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

// StatArgs have a Fid
type StatArgs struct {
	Fid Fid
}

// StatResp is the response to a stat
type StatResp struct {
	F FileInfo
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

// FileInfo has a stat, name, and Fid
type FileInfo struct {
	Stat syscall.Stat_t
	Name string
	Fid  Fid
}

// Nameresp is a response for operations on names. It has a FI and Fid.
type Nameresp struct {
	FI  FileInfo
	Fid Fid
}

// FS holds information about file servers
type FS interface {
	Attach(...string) (Node, error)
}

// Node is the interface for a server, requiring implementations for Attach and FI.
type Node interface {
	Stat() (*FileInfo, error)
	Sys() FileInfo
	Readdir() ([]Node, error)
}
