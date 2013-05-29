package so9p

import (
	"fmt"
	"log"
	"os"
	"path"
)
var servermap map[Fid]*sFid
var serverFid = Fid(2)

func (server *So9ps) FullPath(name string) string {
	/* push a / onto the front of path. Then clean it.
	 * This removes attempts to walk out of the tree.
	 */
	name = path.Clean(path.Join("/", name))
	finalPath := path.Join(server.Path, name)
	/* walk to whatever the new path is -- may be same as old */
	if DebugPrint {
		fmt.Printf("full %v\n", finalPath)
	}
	return name
}
func (server *So9ps) Attach(Args *Nameargs, Resp *Nameresp) (err error) {
	if DebugPrint {
		fmt.Printf("attach args %v resp %v\n", Args, Resp)
	}

	name := server.FullPath(Args.Name)
	if DebugPrint {
		fmt.Printf("Stat: fullpath is %v\n", name)
	}

	n, err := server.Fs.Root()
	if err != nil {
	   log.Printf("No node for root %v\n", err)
	   return nil
	   }
	Resp.FI, err = n.FI(Args.Name)
	if err != nil {
	   log.Printf("FI fails for %v\n", Args.Name)
	   return nil
	   }
	Resp.Fid = Args.Fid
	server.Node = n
	servermap = make(map[Fid]*sFid, 128)
	servermap[Args.Fid] = &sFid{n}
	return err
}

func (server *So9ps) Stat(Args *Newargs, Resp *Nameresp) (err error) {
	if DebugPrint {
		fmt.Printf("Stat: args %v resp %v\n", Args, Resp)
	}

	name := server.FullPath(Args.Name)
	if DebugPrint {
		fmt.Printf("Stat: fullpath is %v\n", name)
	}

	n := server.Node
	if fs, ok := n.(interface {
		FI(string) (FileInfo, error)
	}); ok {
		fi, err := fs.FI(name)
		if DebugPrint {
			fmt.Printf("fs.FI returns (%v, %v)\n", fi, err)
		}
		if err != nil {
			log.Print("Stat", err)
			return nil
		}
		Resp.FI = fi
		return nil
	}

	if (DebugPrint){
	   fmt.Printf("no FI method\n")
	   }
	return err
}

func (server *So9ps) Create(Args *Newargs, Resp *Nameresp) (err error) {
	if DebugPrint {
		fmt.Printf("Create args %v resp %v\n", Args, Resp)
	}

	name := server.FullPath(Args.Name)
	if DebugPrint {
		fmt.Printf("Create: fullpath is %v\n", name)
	}

	n := server.Node
	if fs, ok := n.(interface {
		Create(string, int, os.FileMode) (Node, error)
	}); ok {
		newNode, err := fs.Create(name, Args.Mode, Args.Perm)
		if DebugPrint {
			fmt.Printf("fs.Create returns (%v, %v)\n", newNode, err)
		}
		if err != nil {
			log.Print("create", err)
			return nil
		}
		Resp.Fid = serverFid
		servermap[Resp.Fid] = &sFid{newNode}
		serverFid = serverFid + 1
	}

	return err
}

func (server *So9ps) Read(Args *Ioargs, Resp *Ioresp) (err error) {
	var serverFid *sFid
	var ok bool
	if DebugPrint {
		fmt.Printf("Read args %v resp %v\n", Args, Resp)
	}
	oFid := Args.Fid
	if serverFid, ok = servermap[oFid]; !ok {
		return err
	}

	if DebugPrint {
		fmt.Printf("read oFid %v\n", serverFid)
	}

	n := serverFid.Node

	if fs, ok := n.(interface {
		Read(int, int64) ([]byte, error)
	}); ok {
		data, err := fs.Read(Args.Len, Args.Off)
		if DebugPrint {
			fmt.Printf("fs.Read returns (%v,%v), fs now %v\n", data, err, fs)
		}
		Resp.Data = data
	}

	return err
}

func (server *So9ps) Write(Args *Ioargs, Resp *Ioresp) (err error) {
	var serverFid *sFid
	var ok bool
	if DebugPrint {
		fmt.Printf("Write args %v resp %v\n", Args, Resp)
	}
	oFid := Args.Fid
	if serverFid, ok = servermap[oFid]; !ok {
		return err
	}

	if DebugPrint {
		fmt.Printf("write oFid %v\n", serverFid)
	}

	n := serverFid.Node

	if fs, ok := n.(interface {
		Write([]byte, int64) (int, error)
	}); ok {
		size, err := fs.Write(Args.Data, Args.Off)
		if DebugPrint {
			fmt.Printf("fs.Write returns (%v,%v), fs now %v\n",
				size, err, fs)
		}
		Resp.Len = size
	}

	return err
}

func (server *So9ps) Close(Args *Ioargs, Resp *Ioresp) (err error) {
	var serverFid *sFid
	var ok bool
	if DebugPrint {
		fmt.Printf("Close args %v resp %v\n", Args, Resp)
	}
	oFid := Args.Fid
	if serverFid, ok = servermap[oFid]; !ok {
		return err
	}

	if DebugPrint {
		fmt.Printf("close oFid %v\n", serverFid)
	}

	n := serverFid.Node

	if fs, ok := n.(interface {
		Close() error
	}); ok {
		err := fs.Close()
		if DebugPrint {
			fmt.Printf("fs.Close returns (%v)\n", err)
		}
	}

	/* either way it's gone */
	delete(servermap, oFid)
	return err
}
func (server *So9ps) ReadDir(Args *Nameargs, Resp *FIresp) (err error) {

	if DebugPrint {
		fmt.Printf("ReadDir args %v resp %v\n", Args, Resp)
	}

	name := server.FullPath(Args.Name)
	if DebugPrint {
		fmt.Printf("Create: fullpath is %v\n", name)
	}

	n := server.Node

	if fs, ok := n.(interface {
		ReadDir(string) ([]FileInfo, error)
	}); ok {
		Resp.FI, err = fs.ReadDir(name)
		if DebugPrint {
			fmt.Printf("fs.ReadDir returns (%v)\n", err)
		}
	}

	return err
}

