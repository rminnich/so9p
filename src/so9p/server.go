package so9p

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var (
	servers     = make(map[Fid]*Server, 128)
	path2Server = make(map[string]*Server)
)

// AddFS adds a file system type.
func AddFS(fsName string, node Node) {
	if _, ok := path2Server[fsName]; ok {
		log.Fatalf("Someone tried to add %v but it already exists", fsName)
	}

	path2Server[fsName] = &Server{Node: node, Fid: Fid(uuid.New())}
}

// FullPath returns the full clean path of a file name.
func FullPath(serverPath string, name string) string {
	/* push a / onto the front of path. Then clean it.
	 * This removes attempts to walk out of the tree.
	 */
	name = filepath.Clean(filepath.Join("/", name))
	finalPath := filepath.Join(serverPath, name)
	/* walk to whatever the new path is -- may be same as old */
	debugPrintf("fullpath %v\n", finalPath)

	return name
}

// GetServerNode gets a server node, using a FID
func GetServerNode(aFid Fid) (Node, error) {
	if s, ok := servers[aFid]; ok {
		return s.Node, nil
	}
	log.Printf("Could not find fid %v in Servers", aFid)
	return null, nil
}

// Attach is the server response ot an attach
func (server *Server) Attach(Args *AttachArgs, Resp *Attachresp) (err error) {

	debugPrintf("Attach: args %v\n", Args)

	name := FullPath(server.FullPath, Args.Name)
	n, ok := path2Server[Args.Name]
	if !ok {
		log.Printf("No node for root %v\n", err)
		return
	}

	Resp.FI, err = n.Node.FI(name)
	if err != nil {
		log.Printf("FI fails for %v\n", name)
		return
	}
	Resp.Fid = server.Fid

	return
}

// Unattach is the server side of an unsttach
func (server *Server) Unattach(Args *NameArgs, Resp *Nameresp) (err error) {

	debugPrintf("Unattach: args %v\n", Args)

	return
}

// Stat is the server side of a stat
func (server *Server) Stat(Args *NewArgs, Resp *Nameresp) (err error) {

	debugPrintf("Stat: args %v\n", Args)

	name := FullPath(server.FullPath, Args.Name)
	n, err := GetServerNode(Args.Fid)

	if fs, ok := n.(interface {
		FI(string) (FileInfo, error)
	}); ok {
		fi, err := fs.FI(name)
		Resp.FI = fi
		debugPrintf("fs.FI returns (%v, %v)\n", fi, err)
	} else {
		debugPrintf("Node has no FI method\n")
		err = errors.New("Unimplemented")
	}

	return
}

// Create implements a server create
func (server *Server) Create(Args *NewArgs, Resp *Nameresp) (err error) {

	debugPrintf("Create: args %v\n", Args)

	name := FullPath(server.FullPath, Args.Name)
	n, err := GetServerNode(Args.Fid)

	if fs, ok := n.(interface {
		Create(string, int, os.FileMode) (Node, error)
	}); ok {
		newNode, err := fs.Create(name, Args.Mode, Args.Perm)
		if err != nil {
			return err
		}
		Resp.Fid = server.Fid
		debugPrintf("fs.Create returns (%v)\n", newNode)
	} else {
		debugPrintf("Node has no Create method\n")
		err = errors.New("Unimplemented")
	}

	return
}

func (server *Server) Read(Args *Ioargs, Resp *Ioresp) (err error) {

	debugPrintf("Read: args %v\n", Args)

	n, err := GetServerNode(Args.Fid)

	if fs, ok := n.(interface {
		ReadAt([]byte, int64) (int, error)
	}); ok {
		Resp.Data = make([]byte, Args.Len)
		Resp.Len, err = fs.ReadAt(Resp.Data, Args.Off)
		debugPrintf("fs.Read @ %v returns (%v, %v)\n", Args.Off, Resp.Len, err)
		// The RPC package has a few limits. The error return combines an error
		// for the RPC and an error for what the RPC is doing. The result is that
		// we have to be careful for the error return for io.EOF.
		// So it goes, it's too nice not to do it this way anyway.
		// if we get ANYTHING, return no error.
		if err == io.EOF {
			Resp.EOF = true
			debugPrintf("server: EOF on read fo %d bytes", Resp.Len)
		}
		if Resp.Len > 0 {
			return nil
		}
	} else {
		debugPrintf("Node has no Read method\n")
		err = errors.New("Unimplemented")
	}

	return
}

func (server *Server) Write(Args *Ioargs, Resp *Ioresp) (err error) {

	debugPrintf("Write: args %v\n", Args)

	n, err := GetServerNode(Args.Fid)

	if fs, ok := n.(interface {
		Write([]byte, int64) (int, error)
	}); ok {
		size, err := fs.Write(Args.Data, Args.Off)
		Resp.Len = size
		debugPrintf("fs.Write returns (%v,%v), fs now %v\n", size, err, fs)
	} else {
		debugPrintf("Node has no Write method\n")
		err = errors.New("Unimplemented")
	}

	return
}

// Close closes a FID
func (server *Server) Close(Args *Ioargs, Resp *Ioresp) (err error) {

	debugPrintf("Close: args %v\n", Args)

	n, err := GetServerNode(Args.Fid)

	if fs, ok := n.(interface {
		Close() error
	}); ok {
		err = fs.Close()
		debugPrintf("fs.Close returns (%v)\n", err)
	} else {
		debugPrintf("Node has no Close method\n")
		err = errors.New("Unimplemented")
	}

	return
}

// ReadDir reads a directory
func (server *Server) ReadDir(Args *NameArgs, Resp *FIresp) (err error) {

	debugPrintf("ReadDir: args %v\n", Args)

	name := FullPath(server.FullPath, Args.Name)
	n, err := GetServerNode(Args.Fid)

	if fs, ok := n.(interface {
		ReadDir(string) ([]FileInfo, error)
	}); ok {
		Resp.FI, err = fs.ReadDir(name)
		debugPrintf("fs.ReadDir returns (%v)\n", err)
	} else {
		debugPrintf("Node has no ReadDir method\n")
		err = errors.New("Unimplemented")
	}

	return
}
