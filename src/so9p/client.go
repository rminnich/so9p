package so9p

import (
	"io"
	"log"
	"os"
)

func NewClientConn(c *Conn) *ClientConn {
	return &ClientConn{Conn: c}
}

// Attach attaches to a so9p server.
func (client *ClientConn) Attach(name string, args ...string) error {
	a := &AttachArgs{Name: name, Args: args}
	var reply Attachresp
	err := client.Call("Server.Attach", a, &reply)
	if debugPrint {
		log.Printf("client: clientattach: %v gets (%v, %v)\n", name, client, err)
	}
	return err
}

// Unattach disconnext from a server.
func (client *ClientConn) Unattach() error {
	args := &NameArgs{Fid: client.Fid}
	var reply Nameresp
	err := client.Client.Call("Server.Unattach", args, &reply)
	if debugPrint {
		log.Printf("Unattach: gets %v\n", err)
	}
	return err
}

// Open opens a file
func (client *ClientConn) Open(Fid Fid, mode int) (*File, error) {
	args := &NameArgs{Fid: Fid, Mode: (mode & (^os.O_CREATE))}
	var reply Nameresp
	err := client.Client.Call("Server.Create", args, &reply)
	if debugPrint {
		log.Printf("client: Open: %v gets %v\n", Fid, err)
	}
	return &File{ClientConn: client, Fid: reply.Fid}, err
}

// Create creates a file
func (client *ClientConn) Create(Fid Fid, name string, mode int, perm os.FileMode) (*File, error) {
	args := &NewArgs{Fid: Fid, Name: name, Mode: mode | os.O_CREATE, Perm: perm}
	var reply Nameresp
	err := client.Client.Call("Server.Create", args, &reply)
	if debugPrint {
		log.Printf("client: Create(: %v gets %v\n", name, err)
	}
	return &File{ClientConn: client, Fid: reply.Fid}, err
}

// Stat implements os.Stat
// We use it so we don't need to implement Walk
func (client *ClientConn) Stat(Fid Fid, name string) (FileInfo, error) {
	args := &StatArgs{Fid: Fid}
	var reply Nameresp
	err := client.Client.Call("Server.Stat", args, &reply)
	if debugPrint {
		log.Printf("client: Stat: %v gets %v, %v\n", name, reply.FI.Stat, err)
	}
	//	reply.FI.Name = path.Base(name)
	return reply.FI, err
}

// ReadDir reads an entire directory.
func (client *ClientConn) ReadDir(Fid Fid) ([]FileInfo, error) {
	args := &NameArgs{Fid: Fid}
	var reply FIresp
	err := client.Client.Call("Server.ReadDir", args, &reply)
	if debugPrint {
		log.Printf("client: ReadDir: %v\n", err)
	}
	return reply.FI, err
}

// Readlink implements os.ReadLink
func (client *ClientConn) Readlink(Fid Fid) (string, error) {
	args := &NameArgs{Fid: Fid}
	var reply FileInfo
	err := client.Client.Call("Server.Stat", args, &reply)
	if debugPrint {
		log.Printf("client: Readlink: %v gets %v, %v\n", Fid, reply, err)
	}
	return "", err
}

// ReadAt implements pread
func (client *File) ReadAt(b []byte, Off int64) (int, error) {
	// if you got an EOF indication, give them one zero-byte
	// read back and then clear the EOF indication.
	if client.EOF {
		client.EOF = false
		return 0, io.EOF
	}
	args := &Ioargs{Fid: client.Fid, Len: len(b), Off: Off}
	var reply Ioresp
	err := client.Client.Call("Server.Read", args, &reply)
	if debugPrint {
		log.Printf("client: ReadAt: %v gets %v\n", reply, err)
	}
	if reply.EOF {
		client.EOF = true
	}
	copy(b, reply.Data)
	return reply.Len, err
}

// Read implemens io.Read
func (client *File) Read(b []byte) (int, error) {
	amt, err := client.ReadAt(b, client.Off)
	if err == nil {
		client.Off += int64(amt)
	}
	return amt, err
}

// WriteAt implements io.WriteAt
func (client *File) WriteAt(Data []byte, Off int64) (int, error) {
	args := &Ioargs{Fid: client.Fid, Data: Data, Off: Off}
	var reply Ioresp
	err := client.Client.Call("Server.Write", args, &reply)
	if debugPrint {
		log.Printf("client: Write: %v gets %v\n", reply, err)
	}
	return reply.Len, err
}

// Write implements io.Write
func (client *File) Write(b []byte) (int, error) {
	amt, err := client.WriteAt(b, client.Off)
	if err == nil {
		client.Off += int64(amt)
	}
	return amt, err
}

// Close implements io.Close
func (client *File) Close() error {
	args := &Ioargs{Fid: client.Fid}
	var reply Ioresp
	err := client.Client.Call("Server.Close", args, &reply)
	if debugPrint {
		log.Printf("client: Close: %v gets %v\n", reply, err)
	}
	return err
}
