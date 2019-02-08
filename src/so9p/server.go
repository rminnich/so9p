package so9p

import (
	//	"errors"
	"fmt"
	//	"io"
	"log"
	//	"os"
	"path/filepath"
)

var (
	fs      = make(map[Fid]FS, 128)
	path2FS = make(map[string]FS)
	nodes   = make(map[Fid]Node)
)

func GetServerNode(f Fid) (Node, error) {
	n, ok := nodes[f]
	if !ok {
		return nil, fmt.Errorf("how the hell did this happen, %v is gone", f)
	}
	return n, nil
}

// AddFS adds a file system type.
func AddFS(fsName string, n FS) {
	if _, ok := path2FS[fsName]; ok {
		log.Fatalf("Someone tried to add %v but it already exists", fsName)
	}

	f := newFid()
	fs[f] = n
	path2FS[fsName] = n
	debugPrintf("AddFS: %v, %v", fsName, path2FS[fsName])
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

// GetGS gets an FS, using a FID
func GetFS(aFid Fid) (FS, error) {
	if s, ok := fs[aFid]; ok {
		return s, nil
	}
	log.Printf("Could not find fid %v in Servers, map %v", aFid, fs)
	return null, nil
}

// Attach is the server response ot an attach
func (server *Server) Attach(Args *AttachArgs, Resp *Attachresp) (err error) {

	debugPrintf("Attach: args %v\n", Args)

	name := FullPath(server.FullPath, Args.Name)
	n, ok := path2FS[Args.Name]
	if !ok {
		log.Printf("No node for root %v\n", err)
		return
	}
	debugPrintf("Attach: found %v", n)
	node, err := n.Attach("") // filepath.Join(Args.Args...))
	if err != nil {
		log.Printf("attach filed at FS: %v", err)
		return
	}

	Resp.FI, err = node.Stat()
	if err != nil {
		log.Printf("FI fails for %v\n", name)
		return
	}
	fid := newFid()
	Resp.Fid = fid
	nodes[fid] = node
	debugPrintf("Attach: resp is %v", Resp)
	return
}

// Stat is the server side of a stat
func (server *Server) Stat(Args *StatArgs, Resp *StatResp) (err error) {
	debugPrintf("Stat: %v", Args)
	n, err := GetServerNode(Args.Fid)
	if err != nil {
		return err
	}

	fi, err := n.Stat()
	if err != nil {
		log.Printf("Stat fails: %v", err)
		return err
	}
	Resp.F = fi
	debugPrintf("fs.FI returns (%v, %v)\n", fi, err)
	return nil
}

// // Create implements a server create
// func (server *Server) Create(Args *NewArgs, Resp *Nameresp) (err error) {

// 	debugPrintf("Create: args %v\n", Args)

// 	name := FullPath(server.FullPath, Args.Name)
// 	n, err := GetServerNode(Args.Fid)
// 	if err != nil {
// 		return err
// 	}

// 	if fs, ok := n.Node.(interface {
// 		Create(string, int, os.FileMode) (Node, error)
// 	}); ok {
// 		newNode, err := fs.Create(name, Args.Mode, Args.Perm)
// 		if err != nil {
// 			return err
// 		}
// 		Resp.Fid = server.Fid
// 		debugPrintf("fs.Create returns (%v)\n", newNode)
// 	} else {
// 		debugPrintf("Node has no Create method\n")
// 		err = errors.New("Unimplemented")
// 	}

// 	return
// }

// // Read implements server read
// func (server *Server) Read(Args *Ioargs, Resp *Ioresp) (err error) {

// 	debugPrintf("Read: args %v\n", Args)

// 	n, err := GetServerNode(Args.Fid)
// 	if err != nil {
// 		return err
// 	}

// 	if fs, ok := n.Node.(interface {
// 		ReadAt([]byte, int64) (int, error)
// 	}); ok {
// 		Resp.Data = make([]byte, Args.Len)
// 		Resp.Len, err = fs.ReadAt(Resp.Data, Args.Off)
// 		debugPrintf("fs.Read @ %v returns (%v, %v)\n", Args.Off, Resp.Len, err)
// 		// The RPC package has a few limits. The error return combines an error
// 		// for the RPC and an error for what the RPC is doing. The result is that
// 		// we have to be careful for the error return for io.EOF.
// 		// So it goes, it's too nice not to do it this way anyway.
// 		// if we get ANYTHING, return no error.
// 		if err == io.EOF {
// 			Resp.EOF = true
// 			debugPrintf("server: EOF on read fo %d bytes", Resp.Len)
// 		}
// 		if Resp.Len > 0 {
// 			return nil
// 		}
// 	} else {
// 		debugPrintf("Node has no Read method\n")
// 		err = errors.New("Unimplemented")
// 	}

// 	return
// }

// // Write implements write.
// func (server *Server) Write(Args *Ioargs, Resp *Ioresp) (err error) {

// 	debugPrintf("Write: args %v\n", Args)

// 	n, err := GetServerNode(Args.Fid)
// 	if err != nil {
// 		return err
// 	}

// 	if fs, ok := n.Node.(interface {
// 		Write([]byte, int64) (int, error)
// 	}); ok {
// 		size, err := fs.Write(Args.Data, Args.Off)
// 		Resp.Len = size
// 		debugPrintf("fs.Write returns (%v,%v), fs now %v\n", size, err, fs)
// 	} else {
// 		debugPrintf("Node has no Write method\n")
// 		err = errors.New("Unimplemented")
// 	}

// 	return
// }

// // Close closes a FID
// func (server *Server) Close(Args *Ioargs, Resp *Ioresp) (err error) {

// 	debugPrintf("Close: args %v\n", Args)

// 	n, err := GetServerNode(Args.Fid)
// 	if err != nil {
// 		return err
// 	}

// 	if fs, ok := n.Node.(interface {
// 		Close() error
// 	}); ok {
// 		err = fs.Close()
// 		debugPrintf("fs.Close returns (%v)\n", err)
// 	} else {
// 		debugPrintf("Node has no Close method\n")
// 		err = errors.New("Unimplemented")
// 	}

// 	return
// }

// // ReadDir reads a directory
// func (server *Server) ReadDir(Args *NameArgs, Resp *FIresp) (err error) {

// 	debugPrintf("ReadDir: args %v\n", Args)

// 	name := FullPath(server.FullPath, Args.Name)
// 	n, err := GetServerNode(Args.Fid)
// 	if err != nil {
// 		return err
// 	}

// 	if fs, ok := n.Node.(interface {
// 		ReadDir(string) ([]FileInfo, error)
// 	}); ok {
// 		Resp.FI, err = fs.ReadDir(name)
// 		debugPrintf("fs.ReadDir returns (%v)\n", err)
// 	} else {
// 		debugPrintf("Node has no ReadDir method\n")
// 		err = errors.New("Unimplemented")
// 	}

// 	return
// }
