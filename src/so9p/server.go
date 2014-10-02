package so9p

import (
	"errors"
	"log"
	"os"
	"path"
)

var (
	fid2sFid    = make(map[Fid]*sFid, 128)
	serverFid   = Fid(2)
	path2Server = make(map[string]Node)
)

func FullPath(serverPath string, name string) string {
	/* push a / onto the front of path. Then clean it.
	 * This removes attempts to walk out of the tree.
	 */
	name = path.Clean(path.Join("/", name))
	finalPath := path.Join(serverPath, name)
	/* walk to whatever the new path is -- may be same as old */
	DebugPrintf("fullpath %v\n", finalPath)

	return name
}

func GetServerNode(aFid Fid) (Node, error) {
	if serverFid, ok := fid2sFid[aFid]; ok {
		return serverFid.Node, nil
	}
	log.Printf("Could not find fid %v in fid2sFid", aFid)
	return null, nil
}

func (server *So9ps) Attach(Args *Nameargs, Resp *Nameresp) (err error) {

	DebugPrintf("Attach: args %v\n", Args)

	name := FullPath(server.Path, Args.Name)
	n, ok := path2Server[Args.Name]
	if !ok {
		log.Printf("No node for root %v\n", err)
		return
	}

	Resp.FI, err = n.FI(name)
	if err != nil {
		log.Printf("FI fails for %v\n", name)
		return
	}
	Resp.Fid = Args.Fid

	Resp.Fid = serverFid
	fid2sFid[Resp.Fid] = &sFid{n}
	serverFid = serverFid + 1

	return
}

func (server *So9ps) Unattach(Args *Nameargs, Resp *Nameresp) (err error) {

	DebugPrintf("Unattach: args %v\n", Args)

	delete(fid2sFid, Args.Fid)
	return
}

func (server *So9ps) Stat(Args *Newargs, Resp *Nameresp) (err error) {

	DebugPrintf("Stat: args %v\n", Args)

	name := FullPath(server.Path, Args.Name)
	n, err := GetServerNode(Args.Fid)

	if fs, ok := n.(interface {
		FI(string) (FileInfo, error)
	}); ok {
		fi, err := fs.FI(name)
		Resp.FI = fi
		DebugPrintf("fs.FI returns (%v, %v)\n", fi, err)
	} else {
		DebugPrintf("Node has no FI method\n")
		err = errors.New("Unimplemented")
	}

	return
}

func (server *So9ps) Create(Args *Newargs, Resp *Nameresp) (err error) {

	DebugPrintf("Create: args %v\n", Args)

	name := FullPath(server.Path, Args.Name)
	n, err := GetServerNode(Args.Fid)

	if fs, ok := n.(interface {
		Create(string, int, os.FileMode) (Node, error)
	}); ok {
		newNode, err := fs.Create(name, Args.Mode, Args.Perm)
		if err != nil {
			return err
		}
		Resp.Fid = serverFid
		fid2sFid[Resp.Fid] = &sFid{newNode}
		serverFid = serverFid + 1
		DebugPrintf("fs.Create returns (%v)\n", newNode)
	} else {
		DebugPrintf("Node has no Create method\n")
		err = errors.New("Unimplemented")
	}

	return
}

func (server *So9ps) Read(Args *Ioargs, Resp *Ioresp) (err error) {
     
	DebugPrintf("Read: args %v\n", Args)

	n, err := GetServerNode(Args.Fid)

	if fs, ok := n.(interface {
		ReadAt([]byte, int64) (int, error)
	}); ok {
		Resp.Data = make([]byte, Args.Len)
		Resp.Len, err = fs.ReadAt(Resp.Data, Args.Off)
		DebugPrintf("fs.Read @ %v returns (%v, %v)\n", Args.Off, Resp.Len, err)
		// if we get ANYTHING, return no error.
		if Resp.Len > 0 {
			return nil
		}
	} else {
		DebugPrintf("Node has no Read method\n")
		err = errors.New("Unimplemented")
	}

	return
}

func (server *So9ps) Write(Args *Ioargs, Resp *Ioresp) (err error) {

	DebugPrintf("Write: args %v\n", Args)

	n, err := GetServerNode(Args.Fid)

	if fs, ok := n.(interface {
		Write([]byte, int64) (int, error)
	}); ok {
		size, err := fs.Write(Args.Data, Args.Off)
		Resp.Len = size
		DebugPrintf("fs.Write returns (%v,%v), fs now %v\n", size, err, fs)
	} else {
		DebugPrintf("Node has no Write method\n")
		err = errors.New("Unimplemented")
	}

	return
}

func (server *So9ps) Close(Args *Ioargs, Resp *Ioresp) (err error) {

	DebugPrintf("Close: args %v\n", Args)

	n, err := GetServerNode(Args.Fid)

	if fs, ok := n.(interface {
		Close() error
	}); ok {
		err = fs.Close()
		DebugPrintf("fs.Close returns (%v)\n", err)
	} else {
		DebugPrintf("Node has no Close method\n")
		err = errors.New("Unimplemented")
	}

	// Is this the right thing to do unconditionally?
	delete(fid2sFid, Args.Fid)
	return
}
func (server *So9ps) ReadDir(Args *Nameargs, Resp *FIresp) (err error) {

	DebugPrintf("ReadDir: args %v\n", Args)

	name := FullPath(server.Path, Args.Name)
	n, err := GetServerNode(Args.Fid)

	if fs, ok := n.(interface {
		ReadDir(string) ([]FileInfo, error)
	}); ok {
		Resp.FI, err = fs.ReadDir(name)
		DebugPrintf("fs.ReadDir returns (%v)\n", err)
	} else {
		DebugPrintf("Node has no ReadDir method\n")
		err = errors.New("Unimplemented")
	}

	return
}
