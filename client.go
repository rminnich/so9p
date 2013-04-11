package main

import (
	"fmt"
	"os"
)

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

