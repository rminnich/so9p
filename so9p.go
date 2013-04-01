package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"net/rpc"
)

type fid int64

type sfid struct {
	Path string
	File *os.File
	Fid fid
}

type So9ps struct {
	Server *rpc.Server
	fileFS
}

type so9pc struct {
	Client *rpc.Client
}

type Nameargs struct {
	Name []string
	Fid  []fid
}

type Nameresp struct {
	FI []os.FileInfo
	Fid []fid
	Err error
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

var servermap map[fid] *sfid

func (node *fileNode) FI() (os.FileInfo, error){
	fi, err := os.Stat(node.Name)
	return fi, err
}

func (fs *fileFS) Root() (node Node, err error) {
	node, err = &fileNode{"/"}, nil
	return
}

func (server *So9ps) Srvattach(Args *Nameargs, Resp *Nameresp) (err error) {
	fmt.Printf("attach args %v resp %v\n", Args, Resp)
	_, err = os.Stat(Args.Name[0])
	if err != nil {
		log.Print("Attach", err)
		return err
	}

	n, err := server.Root()
	Resp.FI[0], err = n.FI()
	Resp.Fid[0] = Args.Fid[0]
	servermap[Args.Fid[0]] = &sfid{Path:Args.Name[0], Fid: Args.Fid[0]}
	return err
}
	
func (server *So9ps) Srvwalk(Args *Nameargs, Resp *Nameresp) (err error) {
	fmt.Printf("Walk args %v resp %v\n", Args, Resp)
	/* ofid valid? */
	ofid := Args.Fid[0]
	if sfid, ok := servermap[ofid]; ! ok {
		return err
	}

	nfid := Args.Fid[1]

	/* shortcut: new name is 0 length */
	if len(Args.Name[0] == 0 {
		servermap[nfid] = servermap[ofid]
		return nil
	}

	/* walk to whatever the new path is -- may be same as old */
	new
	_, err = os.Stat(Args.Name[0])
	if err != nil {
		log.Print("Attach", err)
		return err
	}

	n, err := server.Root()
	Resp.FI[0], err = n.FI()
	Resp.Fid[0] = Args.Fid[0]
	servermap[Args.Fid[0]] = &sfid{Path:Args.Name[0], Fid: Args.Fid[0]}
	return err
}
	
func (client *so9pc) attach (name string, file fid) (os.FileInfo, error) {
	var fi os.FileInfo
	args := &Nameargs{[]string{name}, []fid{file}}
	var reply Nameresp
	err := client.Client.Call("So9ps.Srvattach", args, &reply)
	fmt.Printf("clientattach: %v gets %v, %v\n", name, fi, err)
	return fi, err
}

func (client *so9pc) walk (name string, file fid) (os.FileInfo, error) {
	var fi os.FileInfo
	args := &Nameargs{[]string{name}, []fid{file}}
	var reply Nameresp
	err := client.Client.Call("So9ps.Walk", args, &reply)
	fmt.Printf("clientwalk: %v gets %v, %v\n", name, fi, err)
	return fi, err
}

func main(){

	if os.Args[1] == "s" {
		S := new(So9ps)
		S.Name = "/"
		rpc.Register(S)
		l, err := net.Listen("tcp", ":1234")
		if err != nil {
			log.Fatal(err)
		}
		rpc.Accept(l)
	} else {
		var client so9pc
		var err error
		client.Client, err = rpc.Dial("tcp", "localhost" + ":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}
		fi, err := client.attach("/", 1)
		if err != nil {
		log.Fatal("attach", err)
	}
	fmt.Printf("fi %v err %v\n", fi, err)
	}
	
}
