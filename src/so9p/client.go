package so9p

import (
	"io"
	"log"
	"os"
)

// Attach attaches to a so9p server.
func (client *Conn) Attach(name string, args ...interface{}) (*Client, error) {
	a := &AttachArgs{Name: name, Args: args}
	var reply Attachresp
	err := client.Call("Server.Attach", a, &reply)
	fi := reply.FI
	if debugPrint {
		log.Printf("client: clientattach: %v gets %v\n", name, err)
	}
	if err != nil {
		return nil, err
	}
	return &Client{Conn: client, fi: fi, Fid: reply.Fid}, err
}

// Unattach disconnext from a server.
func (client *Client) Unattach() error {
	args := &NameArgs{Fid: client.Fid}
	var reply Nameresp
	err := client.Client.Call("Server.Unattach", args, &reply)
	if debugPrint {
		log.Printf("Unattach: gets %v\n", err)
	}
	return err
}

// Open opens a file, creating if needed.
func (client *Client) Open(name string, mode int) (*File, error) {
	args := &NameArgs{Fid: client.Fid, Name: name, Mode: (mode & (^os.O_CREATE))}
	var reply Nameresp
	err := client.Client.Call("Server.Create", args, &reply)
	if debugPrint {
		log.Printf("client: Open: %v gets %v\n", name, err)
	}
	return &File{Client: client, Fid: reply.Fid}, err
}

// Create creates a file
func (client *Client) Create(name string, mode int, perm os.FileMode) (*File, error) {
	args := &NewArgs{Fid: client.Fid, Name: name, Mode: mode | os.O_CREATE, Perm: perm}
	var reply Nameresp
	err := client.Client.Call("Server.Create", args, &reply)
	if debugPrint {
		log.Printf("client: Create(: %v gets %v\n", name, err)
	}
	return &File{Client: client, Fid: reply.Fid}, err
}

// Stat implements os.Stat
func (client *Client) Stat(name string) (FileInfo, error) {
	args := &NewArgs{Fid: client.Fid, Name: name}
	var reply Nameresp
	err := client.Client.Call("Server.Stat", args, &reply)
	if debugPrint {
		log.Printf("client: Stat: %v gets %v, %v\n", name, reply.FI.Stat, err)
	}
	//	reply.FI.Name = path.Base(name)
	return reply.FI, err
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

// ReadDir reads an entire directory.
func (client *Client) ReadDir(name string) ([]FileInfo, error) {
	args := &NameArgs{Fid: client.Fid, Name: name}
	var reply FIresp
	err := client.Client.Call("Server.ReadDir", args, &reply)
	if debugPrint {
		log.Printf("client: ReadDir: %v\n", err)
	}
	return reply.FI, err
}

// Readlink implements os.ReadLink
func (client *Client) Readlink(name string) (string, error) {
	args := &NameArgs{Fid: client.Fid, Name: name}
	var reply FileInfo
	err := client.Client.Call("Server.Stat", args, &reply)
	if debugPrint {
		log.Printf("client: Readlink: %v gets %v, %v\n", name, reply, err)
	}
	return reply.Link, err
}
