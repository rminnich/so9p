package main

import (
       "fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
)

var debugprint = false

func main() {

	if len(os.Args) > 1 && os.Args[1] == "s" {
		servermap = make(map[fid]*sfid, 128)
		S := new(So9ps)
		S.Path = "/"
		rpc.Register(S)
		l, err := net.Listen("tcp", ":1234")
		if err != nil {
			log.Fatal(err)
		}
		rpc.Accept(l)
	} else {
		var client so9pc
		var err error
		rootfid := fid(1)
		client.Client, err = rpc.Dial("tcp", "localhost"+":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}

		fi, err := client.attach("/", rootfid)
		if err != nil {
			log.Fatal("attach", err)
		}
		if debugprint {
		   fmt.Printf("attach fi %v\n", fi)
		   }
		if len(os.Args) == 1 {
			return
		}
		hostfid, err := client.open(os.Args[1], os.O_RDONLY)
		if debugprint {
			fmt.Printf("open %v: %v, %v\n", os.Args[1], hostfid, err)
		}
		for i := 1; i < 1048576; i = i * 2 {
			start := time.Now()
			data, err := client.read(hostfid, i, 0)
			cost := time.Since(start)
			if err != nil {
				log.Fatal("read", err)
			}
			fmt.Printf("%v took %v\n", len(data), cost)

		}
		if len(os.Args) < 3 {
			return
		}

		copyfid, err := client.create(os.Args[2], os.O_WRONLY, 0666)
		if debugprint {
			fmt.Printf("create %v: %v, %v\n", os.Args[2], hostfid, err)
		}
		for i := 1; i < 1048576; i = i * 2 {
			start := time.Now()
			data, err := client.read(hostfid, i, 0)
			_, err = client.write(copyfid, data, 0)
			cost := time.Since(start)
			if err != nil {
				log.Fatal("read", err)
			}
			fmt.Printf("%v took %v\n", len(data), cost)

		}
	}

}
