package so9p

import (
	"fmt"
	"os"
	//"path"
)

func (client *So9pConn) Attach(name string, file Fid) (*So9pc, error) {
	args := &Nameargs{Name: name, Fid: file}
	var reply Nameresp
	err := client.Call("So9ps.Attach", args, &reply)
	fi := reply.FI
	if DebugPrint {
		fmt.Printf("clientattach: %v gets %v\n", name, err)
	}
	if err != nil {
		return nil, err
	}
	return &So9pc{So9pConn: client, fi: fi, Fid: reply.Fid}, err
}

func (client *So9pc) Open(name string, mode int) (*So9file, error) {
	args := &Nameargs{Fid: client.Fid, Name: name, Mode: (mode & (^os.O_CREATE))}
	var reply Nameresp
	err := client.Client.Call("So9ps.Create", args, &reply)
	if DebugPrint {
		fmt.Printf("Open: %v gets %v\n", name, err)
	}
	return &So9file{So9pc: client, Fid: reply.Fid}, err
}

func (client *So9pc) Create(name string, mode int, perm os.FileMode) (*So9file, error) {
	args := &Newargs{Fid: client.Fid, Name: name, Mode: mode | os.O_CREATE, Perm: perm}
	var reply Nameresp
	err := client.Client.Call("So9ps.Create", args, &reply)
	if DebugPrint {
		fmt.Printf("Create(: %v gets %v\n", name, err)
	}
	return &So9file{So9pc: client, Fid: reply.Fid}, err
}

func (client *So9pc) Stat(name string) (FileInfo, error) {
	args := &Newargs{Fid: client.Fid, Name: name}
	var reply Nameresp
	err := client.Client.Call("So9ps.Stat", args, &reply)
	if DebugPrint {
		fmt.Printf("Stat: %v gets %v, %v\n", name, reply.FI.Stat, err)
	}
	//	reply.FI.Name = path.Base(name)
	return reply.FI, err
}

// Read implements the io.ReaderAt interface.
func (client *So9file) Read(Len int, Off int64) ([]byte, error) {
	args := &Ioargs{Fid: client.Fid, Len: Len, Off: Off}
	var reply Ioresp
	err := client.Client.Call("So9ps.Read", args, &reply)
	if DebugPrint {
		fmt.Printf("Read: %v gets %v\n", reply, err)
	}
	return reply.Data, err
}

func (client *So9file) Write(Data []byte, Off int64) (int, error) {
	args := &Ioargs{Fid: client.Fid, Data: Data, Off: Off}
	var reply Ioresp
	err := client.Client.Call("So9ps.Write", args, &reply)
	if DebugPrint {
		fmt.Printf("Write: %v gets %v\n", reply, err)
	}
	return reply.Len, err
}

func (client *So9file) Close() error {
	args := &Ioargs{Fid: client.Fid}
	var reply Ioresp
	err := client.Client.Call("So9ps.Close", args, &reply)
	if DebugPrint {
		fmt.Printf("Close: %v gets %v\n", reply, err)
	}
	return err
}

func (client *So9pc) ReadDir(name string) ([]FileInfo, error) {
	args := &Nameargs{Fid: client.Fid, Name: name}
	var reply FIresp
	err := client.Client.Call("So9ps.ReadDir", args, &reply)
	if DebugPrint {
		fmt.Printf("ReadDir: %v gets %v\n", reply, err)
	}
	return reply.FI, err
}

func (client *So9pc) Readlink(name string) (string, error) {
	args := &Nameargs{Fid: client.Fid, Name: name}
	var reply FileInfo
	err := client.Client.Call("So9ps.Stat", args, &reply)
	if DebugPrint {
		fmt.Printf("Readlink: %v gets %v, %v\n", reply, err)
	}
	return reply.Link, err
}
