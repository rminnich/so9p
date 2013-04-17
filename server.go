package main

import (
	"fmt"
	"log"
	"os"
)

func (server *So9ps) Attach(Args *Nameargs, Resp *Nameresp) (err error) {
	fmt.Printf("attach args %v resp %v\n", Args, Resp)
	_, err = os.Stat(Args.Name)
	if err != nil {
		log.Print("Attach", err)
		return err
	}

	n, err := server.Fs.Root()
	Resp.FI, err = n.FI()
	Resp.Fid = Args.Fid
	servermap[Args.Fid] = &sfid{n}
	return err
}

func (server *So9ps) Walk(Args *Nameargs, Resp *Nameresp) (err error) {
	var serverfid *sfid
	var ok bool
	fmt.Printf("Walk args %v resp %v\n", Args, Resp)
	/* ofid valid? */
	ofid := Args.Fid
	if serverfid, ok = servermap[ofid]; !ok {
		return err
	}

	fmt.Printf("ofid %v\n", serverfid)
	nfid := Args.NFid

	/* shortcut: new name is 0 length */
	if len(Args.Name) == 0 {
		servermap[nfid] = servermap[ofid]
		return nil
	}
	n := serverfid.Node
	dirfi, err := n.FI()
	if err != nil {
		return err
	}
	fmt.Printf("WALK: dirfi is %v, fullpath is %v\n", dirfi, dirfi.FullPath())
	walkTo := Args.Name
	/* walk to whatever the new path is -- may be same as old */
	if fs, ok := n.(interface {
		Walk(string) (Node, error)
	}); ok {
		newNode, err := fs.Walk(walkTo)
		fmt.Printf("fs.Walk returns (%v, %v)\n", newNode, err)
		if err != nil {
			log.Print("walk", err)
			return nil
		}
		if stat, ok := newNode.(interface {
			FI() (FileInfo, error)
		}); ok {
			fmt.Printf("stat seems to exist, ...\n");
			Resp.FI, err = stat.FI()
			if err != nil {
				log.Print("walk", err)
				return nil
			}
		} else {
			return nil
		}
		Resp.Fid = Args.NFid
		servermap[Args.NFid] = &sfid{newNode}
	}

	return err
}

func (server *So9ps) Open(Args *Nameargs, Resp *Nameresp) (err error) {
	var serverfid *sfid
	var ok bool
	fmt.Printf("Open args %v resp %v\n", Args, Resp)
	/* ofid valid? */
	ofid := Args.Fid
	if serverfid, ok = servermap[ofid]; !ok {
	   	      log.Print(err)
		return err
	}

	fmt.Printf("ofid %v\n", serverfid)
	
	n := serverfid.Node

	if fs, ok := n.(interface {
		Open() (error)
	}); ok {
		err := fs.Open()
		fmt.Printf("fs.Open returns (%v), fs now %v\n", err, fs)
	}

	return err
}

func (server *So9ps) Read(Args *Ioargs, Resp *Ioresp) (err error) {
	var serverfid *sfid
	var ok bool
	fmt.Printf("Read args %v resp %v\n", Args, Resp)
	ofid := Args.Fid
	if serverfid, ok = servermap[ofid]; !ok {
		return err
	}

	fmt.Printf("read ofid %v\n", serverfid)
	
	n := serverfid.Node

	if fs, ok := n.(interface {
		Read(int, int64) ([]byte, error)
	}); ok {
		data, err := fs.Read(Args.Len, Args.Off)
		fmt.Printf("fs.Read returns (%v,%v), fs now %v\n",
				    data,err, fs)
		Resp.Data = data
	}

	return err
}

func (server *So9ps) Write(Args *Ioargs, Resp *Ioresp) (err error) {
	var serverfid *sfid
	var ok bool
	fmt.Printf("Write args %v resp %v\n", Args, Resp)
	ofid := Args.Fid
	if serverfid, ok = servermap[ofid]; !ok {
		return err
	}

	fmt.Printf("write ofid %v\n", serverfid)
	
	n := serverfid.Node

	if fs, ok := n.(interface {
		Write([]byte, int64) (int, error)
	}); ok {
		size, err := fs.Write(Args.Data, Args.Off)
		fmt.Printf("fs.Write returns (%v,%v), fs now %v\n",
				    size, err, fs)
		Resp.Len = size
	}

	return err
}

func (server *So9ps) Close(Args *Ioargs, Resp *Ioresp) (err error) {
	var serverfid *sfid
	var ok bool
	fmt.Printf("Close args %v resp %v\n", Args, Resp)
	ofid := Args.Fid
	if serverfid, ok = servermap[ofid]; !ok {
		return err
	}

	fmt.Printf("close ofid %v\n", serverfid)
	
	n := serverfid.Node

	if fs, ok := n.(interface {
		Close() (error)
	}); ok {
		err := fs.Close()
		fmt.Printf("fs.Close returns (%v)\n",err)
	}

	/* either way it's gone */
	delete(servermap, ofid)
	return err
}

