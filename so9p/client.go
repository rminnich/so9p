package so9p

import (
	"fmt"
	"os"
	//"path"
)

func (client *So9pc) Attach(name string, file Fid) (FileInfo, error) {
	args := &Nameargs{Name:name, Fid:file}
	var reply Nameresp
	err := client.Client.Call("So9ps.Attach", args, &reply)
	fi := reply.FI
	if DebugPrint {
		fmt.Printf("clientattach: %v gets %v, %v\n", name, fi, err)
	}
	return fi, err
}

func (client *So9pc) Open(name string, mode int) (Fid, error)  {
	args := &Nameargs{Name: name, Mode: (mode&(^os.O_CREATE))}
	var reply Nameresp
	err := client.Client.Call("So9ps.Create", args, &reply)
	if DebugPrint {
		fmt.Printf("Open: %v gets %v, %v\n", name, err)
	}
	return reply.Fid, err
}

func (client *So9pc) Create(name string, mode int, perm os.FileMode) (Fid, error)  {
	args := &Newargs{Name: name, Mode: mode|os.O_CREATE, Perm: perm}
	var reply Nameresp
	err := client.Client.Call("So9ps.Create", args, &reply)
	if DebugPrint {
		fmt.Printf("Create(: %v gets %v, %v\n", name, err)
	}
	return reply.Fid, err
}

func (client *So9pc) Stat(name string) (FileInfo, error)  {
	args := &Newargs{Name: name}
	var reply Nameresp
	err := client.Client.Call("So9ps.Stat", args, &reply)
	if DebugPrint {
		fmt.Printf("Stat: %v gets %v, %v\n", name, reply.FI.Stat, err)
	}
//	reply.FI.Name = path.Base(name)
	return reply.FI, err
}

func (client *So9pc) Read(file Fid, Len int, Off int64) ([]byte, error) {
	args := &Ioargs{Fid: file, Len: Len, Off: Off}
	var reply Ioresp
	err := client.Client.Call("So9ps.Read", args, &reply)
	if DebugPrint {
		fmt.Printf("Read: %v gets %v, %v\n", reply, err)
	}
	return reply.Data, err
}

func (client *So9pc) Write(file Fid, Data []byte, Off int64) (int, error) {
	args := &Ioargs{Fid: file, Data: Data, Off: Off}
	var reply Ioresp
	err := client.Client.Call("So9ps.Write", args, &reply)
	if DebugPrint {
		fmt.Printf("Write: %v gets %v, %v\n", reply, err)
	}
	return reply.Len, err
}

func (client *So9pc) Close(file Fid) error {
	args := &Ioargs{Fid: file}
	var reply Ioresp
	err := client.Client.Call("So9ps.Close", args, &reply)
	if DebugPrint {
		fmt.Printf("Close: %v gets %v, %v\n", reply, err)
	}
	return err
}

func (client *So9pc) ReadDir(name string) ([]FileInfo, error) {
	args := &Nameargs{Name:name}
	var reply FIresp
	err := client.Client.Call("So9ps.ReadDir", args, &reply)
	if DebugPrint {
		fmt.Printf("ReadDir: %v gets %v, %v\n", reply, err)
	}
	return reply.FI, err
}

func (client *So9pc) Readlink(name string) (string, error) {
	args := &Nameargs{Name:name}
	var reply FileInfo
	err := client.Client.Call("So9ps.Stat", args, &reply)
	if DebugPrint {
		fmt.Printf("Readlink: %v gets %v, %v\n", reply, err)
	}
	return reply.Link, err
}
