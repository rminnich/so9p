package so9p

import (
	"io"
	"log"
	"os"
)

func (client *So9pConn) Attach(name string) (*So9pc, error) {
	args := &Nameargs{Name: name}
	var reply Nameresp
	err := client.Call("So9ps.Attach", args, &reply)
	fi := reply.FI
	if DebugPrint {
		log.Printf("client: clientattach: %v gets %v\n", name, err)
	}
	if err != nil {
		return nil, err
	}
	return &So9pc{So9pConn: client, fi: fi, Fid: reply.Fid}, err
}

func (client *So9pc) Unattach() error {
	args := &Nameargs{Fid: client.Fid}
	var reply Nameresp
	err := client.Client.Call("So9ps.Unattach", args, &reply)
	if DebugPrint {
		log.Printf("Unattach: gets %v\n", err)
	}
	return err
}

func (client *So9pc) Open(name string, mode int) (*So9file, error) {
	args := &Nameargs{Fid: client.Fid, Name: name, Mode: (mode & (^os.O_CREATE))}
	var reply Nameresp
	err := client.Client.Call("So9ps.Create", args, &reply)
	if DebugPrint {
		log.Printf("client: Open: %v gets %v\n", name, err)
	}
	return &So9file{So9pc: client, Fid: reply.Fid}, err
}

func (client *So9pc) Create(name string, mode int, perm os.FileMode) (*So9file, error) {
	args := &Newargs{Fid: client.Fid, Name: name, Mode: mode | os.O_CREATE, Perm: perm}
	var reply Nameresp
	err := client.Client.Call("So9ps.Create", args, &reply)
	if DebugPrint {
		log.Printf("client: Create(: %v gets %v\n", name, err)
	}
	return &So9file{So9pc: client, Fid: reply.Fid}, err
}

func (client *So9pc) Stat(name string) (FileInfo, error) {
	args := &Newargs{Fid: client.Fid, Name: name}
	var reply Nameresp
	err := client.Client.Call("So9ps.Stat", args, &reply)
	if DebugPrint {
		log.Printf("client: Stat: %v gets %v, %v\n", name, reply.FI.Stat, err)
	}
	//	reply.FI.Name = path.Base(name)
	return reply.FI, err
}

func (client *So9file) ReadAt(b []byte, Off int64) (int, error) {
	// if you got an EOF indication, give them one zero-byte
	// read back and then clear the EOF indication.
	if client.EOF {
		client.EOF = false
		return 0, io.EOF
	}
	args := &Ioargs{Fid: client.Fid, Len: len(b), Off: Off}
	var reply Ioresp
	err := client.Client.Call("So9ps.Read", args, &reply)
	if DebugPrint {
		log.Printf("client: ReadAt: %v gets %v\n", reply, err)
	}
	if reply.EOF {
		client.EOF = true
	}
	copy(b, reply.Data)
	return reply.Len, err
}

func (client *So9file) Read(b []byte) (int, error) {
	amt, err := client.ReadAt(b, client.Off)
	if err == nil {
		client.Off += int64(amt)
	}
	return amt, err
}

func (client *So9file) WriteAt(Data []byte, Off int64) (int, error) {
	args := &Ioargs{Fid: client.Fid, Data: Data, Off: Off}
	var reply Ioresp
	err := client.Client.Call("So9ps.Write", args, &reply)
	if DebugPrint {
		log.Printf("client: Write: %v gets %v\n", reply, err)
	}
	return reply.Len, err
}

func (client *So9file) Write(b []byte) (int, error) {
	amt, err := client.WriteAt(b, client.Off)
	if err == nil {
		client.Off += int64(amt)
	}
	return amt, err
}

func (client *So9file) Close() error {
	args := &Ioargs{Fid: client.Fid}
	var reply Ioresp
	err := client.Client.Call("So9ps.Close", args, &reply)
	if DebugPrint {
		log.Printf("client: Close: %v gets %v\n", reply, err)
	}
	return err
}

func (client *So9pc) ReadDir(name string) ([]FileInfo, error) {
	args := &Nameargs{Fid: client.Fid, Name: name}
	var reply FIresp
	err := client.Client.Call("So9ps.ReadDir", args, &reply)
	if DebugPrint {
		log.Printf("client: ReadDir: %v\n", err)
	}
	return reply.FI, err
}

func (client *So9pc) Readlink(name string) (string, error) {
	args := &Nameargs{Fid: client.Fid, Name: name}
	var reply FileInfo
	err := client.Client.Call("So9ps.Stat", args, &reply)
	if DebugPrint {
		log.Printf("client: Readlink: %v gets %v, %v\n", reply, err)
	}
	return reply.Link, err
}
