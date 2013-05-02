package so9p

import (
	"fmt"
	"os"
)

func (client *So9pc) Attach(name string, file fid) (FileInfo, error) {
	args := &Nameargs{Name:name, Fid:file}
	var reply Nameresp
	err := client.Client.Call("So9ps.Attach", args, &reply)
	fi := reply.FI
	if DebugPrint {
		fmt.Printf("clientattach: %v gets %v, %v\n", name, fi, err)
	}
	return fi, err
}

func (client *So9pc) Open(name string, mode int) (fid, error)  {
	args := &Nameargs{Name: name, Mode: (mode&(^os.O_CREATE))}
	var reply Nameresp
	err := client.Client.Call("So9ps.Create", args, &reply)
	if DebugPrint {
		fmt.Printf("clientopen: %v gets %v, %v\n", name, err)
	}
	return reply.Fid, err
}

func (client *So9pc) Create(name string, mode int, perm os.FileMode) (fid, error)  {
	args := &Newargs{Name: name, Mode: mode|os.O_CREATE, Perm: perm}
	var reply Nameresp
	err := client.Client.Call("So9ps.Create", args, &reply)
	if DebugPrint {
		fmt.Printf("clientopen: %v gets %v, %v\n", name, err)
	}
	return reply.Fid, err
}

func (client *So9pc) Read(file fid, Len int, Off int64) ([]byte, error) {
	args := &Ioargs{Fid: file, Len: Len, Off: Off}
	var reply Ioresp
	err := client.Client.Call("So9ps.Read", args, &reply)
	if DebugPrint {
		fmt.Printf("clientopen: %v gets %v, %v\n", reply, err)
	}
	return reply.Data, err
}

func (client *So9pc) Write(file fid, Data []byte, Off int64) (int, error) {
	args := &Ioargs{Fid: file, Data: Data, Off: Off}
	var reply Ioresp
	err := client.Client.Call("So9ps.Write", args, &reply)
	if DebugPrint {
		fmt.Printf("clientopen: %v gets %v, %v\n", reply, err)
	}
	return reply.Len, err
}

func (client *So9pc) Close(file fid) error {
	args := &Ioargs{Fid: file}
	var reply Ioresp
	err := client.Client.Call("So9ps.Close", args, &reply)
	if DebugPrint {
		fmt.Printf("clientopen: %v gets %v, %v\n", reply, err)
	}
	return err
}

func (client *So9pc) Readdir(file fid) ([]FileInfo, error) {
	args := &Ioargs{Fid: file}
	var reply FIresp
	err := client.Client.Call("So9ps.ReadDir", args, &reply)
	if DebugPrint {
		fmt.Printf("clientopen: %v gets %v, %v\n", reply, err)
	}
	return reply.FI, err
}
