package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"path"
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

type Nameresp struct {
	FI  os.FileInfo
	Fid fid
}

type FS interface {
	Root() (Node, error)
}

type fileFS struct {
	fileNode
}

type Node interface {
	FI() (os.FileInfo, error)
}

type fileNode struct {
	Name string
}

var servermap map[fid]*sfid
var clientfid fid

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

	newNode := &fileNode{fi.Name()}
	return newNode, err
}

func (node *fileNode) FI() (os.FileInfo, error) {
	fi, err := os.Stat(node.Name)
	return fi, err
}

func (fs *fileFS) Root() (node Node, err error) {
	node, err = &fileNode{"/"}, nil
	return
}

func (server *So9ps) Srvattach(Args *Nameargs, Resp *Nameresp) (err error) {
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

func (server *So9ps) Srvwalk(Args *Nameargs, Resp *Nameresp) (err error) {
	var serverfid *sfid
	var ok bool
	fmt.Printf("Walk args %v resp %v\n", Args, Resp)
	/* ofid valid? */
	ofid := Args.Fid
	if serverfid, ok = servermap[ofid]; !ok {
		return err
	}

	nfid := Args.NFid

	/* shortcut: new name is 0 length */
	if len(Args.Name) == 0 {
		servermap[nfid] = servermap[ofid]
		return nil
	}
	n := serverfid.Node

	walkTo := Args.Name
	/* walk to whatever the new path is -- may be same as old */
	if fs, ok := n.(interface {
		Walk(string) (Node, error)
	}); ok {
		newNode, err := fs.Walk(walkTo)
		Resp.FI, err = newNode.FI()
		Resp.Fid = Args.Fid
		servermap[Args.Fid] = &sfid{newNode}
		if err != nil {
			log.Print("walk", err)
			return err
		}
	}

	if err != nil {
		log.Print("Walk", err)
	}

	return err
}

func (client *so9pc) attach(name string, file fid) (os.FileInfo, error) {
	var fi os.FileInfo
	args := &Nameargs{name, file, file}
	var reply Nameresp
	err := client.Client.Call("So9ps.Srvattach", args, &reply)
	fmt.Printf("clientattach: %v gets %v, %v\n", name, fi, err)
	return fi, err
}

func (client *so9pc) walk(file fid, name string) (fid, os.FileInfo, error) {
	var fi os.FileInfo
	clientfid++
	newfid := clientfid
	args := &Nameargs{name, file, newfid}
	var reply Nameresp
	err := client.Client.Call("So9ps.Walk", args, &reply)
	fmt.Printf("clientwalk: %v gets %v, %v\n", name, fi, err)
	return newfid, fi, err
}

func main() {

	if os.Args[1] == "s" {
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
