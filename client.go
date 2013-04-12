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

func (client *so9pc) open(file fid) (error) {
	args := &Nameargs{Fid: file}
	var reply Nameresp
	err := client.Client.Call("So9ps.Open", args, &reply)
	fmt.Printf("clientopen: %v gets %v, %v\n", file, err)
	return err
}

func (client *so9pc) read(file fid, Len int, Off int64) ([]byte, error) {
	args := &Ioargs{Fid: file, Len: Len, Off: Off}
	var reply Ioresp
	err := client.Client.Call("So9ps.Read", args, &reply)
	fmt.Printf("clientopen: %v gets %v, %v\n", reply, err)
	return reply.Data, err
}
