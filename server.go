package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

func (server *So9ps) FullPath(name string) string {
	/* push a / onto the front of path. Then clean it.
	 * This removes attempts to walk out of the tree.
	 */
	name = path.Clean(path.Join("/", name))
	finalPath := path.Join(server.Path, name)
	/* walk to whatever the new path is -- may be same as old */
	if debugprint {
		fmt.Printf("full %v\n", finalPath)
	}
	return name
}
func (server *So9ps) Attach(Args *Nameargs, Resp *Nameresp) (err error) {
	if debugprint {
		fmt.Printf("attach args %v resp %v\n", Args, Resp)
	}
	_, err = os.Stat(Args.Name)
	if err != nil {
		log.Print("Attach", err)
		return err
	}

	n, err := server.Fs.Root()
	if err != nil {
	   log.Printf("No node for root %v\n", err)
	   return nil
	   }
	Resp.FI, err = n.FI(Args.Name)
	Resp.Fid = Args.Fid
	server.Node = n
	servermap[Args.Fid] = &sfid{n}
	return err
}

func (server *So9ps) Create(Args *Newargs, Resp *Nameresp) (err error) {
	if debugprint {
		fmt.Printf("Create args %v resp %v\n", Args, Resp)
	}

	name := server.FullPath(Args.Name)
	if debugprint {
		fmt.Printf("Create: fullpath is %v\n", name)
	}

	n := server.Node
	if fs, ok := n.(interface {
		Create(string, int, os.FileMode) (Node, error)
	}); ok {
		newNode, err := fs.Create(name, Args.mode, Args.perm)
		if debugprint {
			fmt.Printf("fs.Create returns (%v, %v)\n", newNode, err)
		}
		if err != nil {
			log.Print("create", err)
			return nil
		}
		Resp.Fid = Args.Fid
		servermap[Args.Fid] = &sfid{newNode}
	}

	return err
}

func (server *So9ps) Open(Args *Nameargs, Resp *Nameresp) (err error) {
	if debugprint {
		fmt.Printf("Open args %v resp %v\n", Args, Resp)
	}
	name := server.FullPath(Args.Name)
	if debugprint {
		fmt.Printf("Open: fullpath is %v\n", name)
	}
	n := server.Node
	if fs, ok := n.(interface {
		Open(name string) (Node, error)
	}); ok {
		newNode, err := fs.Open(name)
		if debugprint {
			fmt.Printf("fs.Open returns (%v), fs now %v\n", err, fs)
		}
		if err != nil {
			log.Print("open: error", err)
			return nil
		}
		Resp.Fid = Args.Fid
		servermap[Args.Fid] = &sfid{newNode}
	}

	return err
}

func (server *So9ps) Read(Args *Ioargs, Resp *Ioresp) (err error) {
	var serverfid *sfid
	var ok bool
	if debugprint {
		fmt.Printf("Read args %v resp %v\n", Args, Resp)
	}
	ofid := Args.Fid
	if serverfid, ok = servermap[ofid]; !ok {
		return err
	}

	if debugprint {
		fmt.Printf("read ofid %v\n", serverfid)
	}

	n := serverfid.Node

	if fs, ok := n.(interface {
		Read(int, int64) ([]byte, error)
	}); ok {
		data, err := fs.Read(Args.Len, Args.Off)
		if debugprint {
			fmt.Printf("fs.Read returns (%v,%v), fs now %v\n", data, err, fs)
		}
		Resp.Data = data
	}

	return err
}

func (server *So9ps) Write(Args *Ioargs, Resp *Ioresp) (err error) {
	var serverfid *sfid
	var ok bool
	if debugprint {
		fmt.Printf("Write args %v resp %v\n", Args, Resp)
	}
	ofid := Args.Fid
	if serverfid, ok = servermap[ofid]; !ok {
		return err
	}

	if debugprint {
		fmt.Printf("write ofid %v\n", serverfid)
	}

	n := serverfid.Node

	if fs, ok := n.(interface {
		Write([]byte, int64) (int, error)
	}); ok {
		size, err := fs.Write(Args.Data, Args.Off)
		if debugprint {
			fmt.Printf("fs.Write returns (%v,%v), fs now %v\n",
				size, err, fs)
		}
		Resp.Len = size
	}

	return err
}

func (server *So9ps) Close(Args *Ioargs, Resp *Ioresp) (err error) {
	var serverfid *sfid
	var ok bool
	if debugprint {
		fmt.Printf("Close args %v resp %v\n", Args, Resp)
	}
	ofid := Args.Fid
	if serverfid, ok = servermap[ofid]; !ok {
		return err
	}

	if debugprint {
		fmt.Printf("close ofid %v\n", serverfid)
	}

	n := serverfid.Node

	if fs, ok := n.(interface {
		Close() error
	}); ok {
		err := fs.Close()
		if debugprint {
			fmt.Printf("fs.Close returns (%v)\n", err)
		}
	}

	/* either way it's gone */
	delete(servermap, ofid)
	return err
}
func (server *So9ps) ReadDir(Args *Ioargs, Resp *FIresp) (err error) {
	var serverfid *sfid
	var ok bool
	if debugprint {
		fmt.Printf("ReadDir args %v resp %v\n", Args, Resp)
	}
	ofid := Args.Fid
	if serverfid, ok = servermap[ofid]; !ok {
		return err
	}

	if debugprint {
		fmt.Printf("readdir ofid %v\n", serverfid)
	}

	n := serverfid.Node

	if fs, ok := n.(interface {
		ReadDir() ([]FileInfo, error)
	}); ok {
		Resp.FI, err = fs.ReadDir()
		if debugprint {
			fmt.Printf("fs.ReadDir returns (%v)\n", err)
		}
	}

	return err
}
