package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
)

var servermap map[fid]*sfid
var clientfid = fid(2)
var debugprint = false

func main() {

	if len(os.Args) > 1 && os.Args[1] == "s" {

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
		rootfid := fid(1)
		fi, err := client.attach("/", rootfid)
		if err != nil {
			log.Fatal("attach", err)
		}
		etcfid, fi, err := client.walk(rootfid, "etc")
		if err != nil {
			log.Fatal("walk", err)
		}
		if debugprint {
			fmt.Printf("Walk: %v, %v, %v\n", etcfid, fi, err)
		}
		filelist, err := client.readdir(etcfid)
		if err != nil {
			log.Fatal("readdir", err)
		}
		fmt.Printf("etc has %v\n", filelist)

		hostfid, fi, err := client.walk(etcfid, "hosts")
		if err != nil {
			log.Fatal("walk", err)
		}
		if debugprint {
			fmt.Printf("Walk to hosts: %v, %v, %v\n", hostfid, fi, err)
		}
		err = client.open(hostfid)
		if err != nil {
			log.Fatal("open", err)
		}
		data, err := client.read(hostfid, 1<<20, 0)

		if err != nil {
			log.Fatal("read", err)
		}
		if debugprint {
			fmt.Printf("Read: %v, %v\n", data, err)
		}
		if len(os.Args) < 2 {
			return
		}
		hostfid, fi, err = client.walk(rootfid, os.Args[1])
		if debugprint {
			fmt.Printf("Walk to %v: %v, %v, %v\n", os.Args[1], hostfid, fi, err)
		}
		err = client.open(hostfid)
		if err != nil {
			log.Fatal("open", err)
		}
		_, err = client.read(hostfid, 1<<20, 0)
		if err != nil {
			log.Fatal("read", err)
		}
		for i := 1; i < 1048576; i = i * 2 {
			start := time.Now()
			data, err = client.read(hostfid, i, 0)
			cost := time.Since(start)
			if err != nil {
				log.Fatal("read", err)
			}
			fmt.Printf("%v took %v\n", len(data), cost)

		}
		err = client.close(hostfid)
		if err != nil {
			log.Fatal("close", err)
		}

	}

}
