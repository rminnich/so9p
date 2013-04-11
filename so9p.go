package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"path"
	"time"
)

type fid int64

type sfid struct {
	Node Node
}

type So9ps struct {
	Server *rpc.Server
	Fs fileFS
}

type so9pc struct {
	Client *rpc.Client
}

type Nameargs struct {
	Name string
	Fid  fid
	NFid  fid
}

type FileInfo struct {
    SName	string
    SSize int64 
    SMode os.FileMode     
    SModTime time.Time
    SIsDir bool      
}

type Nameresp struct {
	FI  FileInfo
	Fid fid
}

type FS interface {
	Root() (Node, error)
}

type fileFS struct {
	fileNode
}

type Node interface {
	FI() (FileInfo, error)
}

type fileNode struct {
	FullPath, Name string
}

var servermap map[fid]*sfid
var clientfid fid

func (fi FileInfo) Name() string {
	return fi.SName
}

func (fi FileInfo) Size() int64 {
	return fi.SSize
}

func (fi FileInfo) Mode() os.FileMode {
	return fi.SMode
}

func (fi FileInfo) ModTime() time.Time {
	return fi.SModTime
}

func (fi FileInfo) IsDir() bool {
	return fi.SIsDir
}

func (fi FileInfo) Sys() interface{} {
	return nil
}

func (node *fileNode) Walk(walkTo string) (Node, error) {
	/* push a / onto the front of path. Then clean it.
	 * This removes attempts to walk out of the tree.
	 */
	walkTo = path.Clean(path.Join("/", walkTo))
	/* walk to whatever the new path is -- may be same as old */
	fi, err := os.Stat(path.Join(node.Name, walkTo))
	if err != nil {
		return nil, err
	}

	newNode := &fileNode{walkTo, fi.Name()}
	return newNode, err
}

func (node *fileNode) FI() (FileInfo, error) {
	var fi FileInfo
	fmt.Printf("FI %v\n", node)
	osfi, err := os.Stat(node.FullPath)

	if err != nil {
		log.Print(err)
		return fi, err
	}
	fi.SName = osfi.Name()
	fi.SSize = osfi.Size()
	fi.SMode = osfi.Mode()
	fi.SModTime = osfi.ModTime()
	fi.SIsDir = osfi.IsDir()
	return fi, err
}

func (fs *fileFS) Root() (node Node, err error) {
	node, err = &fileNode{"/", "/"}, nil
	return
}

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

	walkTo := path.Join(dirfi.SName, Args.Name)
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
		Resp.Fid = Args.Fid
		servermap[Args.Fid] = &sfid{newNode}
	}

	return err
}

func (client *so9pc) attach(name string, file fid) (os.FileInfo, error) {
	var fi os.FileInfo
	args := &Nameargs{name, file, file}
	var reply Nameresp
	err := client.Client.Call("So9ps.Attach", args, &reply)
	fmt.Printf("clientattach: %v gets %v, %v\n", name, fi, err)
	fi = reply.FI
	return fi, err
}

func (client *so9pc) walk(file fid, name string) (fid, os.FileInfo, error) {
	var fi os.FileInfo
	clientfid++
	newfid := clientfid
	args := &Nameargs{name, file, newfid}
	var reply Nameresp
	err := client.Client.Call("So9ps.Walk", args, &reply)
	fi = reply.FI
	fmt.Printf("clientwalk: %v gets %v, %v\n", name, fi, err)
	return newfid, fi, err
}

func main() {

	if os.Args[1] == "s" {

		servermap = make(map[fid]*sfid, 128)
		S := new(So9ps)
		S.Fs.Name = "/"
		rpc.Register(S)
		l, err := net.Listen("tcp", ":1234")
		if err != nil {
			log.Fatal(err)
		}
		rpc.Accept(l)
	} else {
		var client so9pc
		var err error
		client.Client, err = rpc.Dial("tcp", "localhost"+":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}
		fi, err := client.attach("/", 1)
		if err != nil {
			log.Fatal("attach", err)
		}
		newfid, fi, err := client.walk(1, "etc")
		if err != nil {
			log.Fatal("walk", err)
		}
		fmt.Printf("newfid %v fi %v err %v\n", newfid, fi, err)
	}

}
