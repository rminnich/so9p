package main

import (
       "flag"
       "fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
	"so9p"
)

var (
    dbg = flag.Bool("dbg", false, "set true for debug")
    server = flag.Bool("s", true, "set true for server")
    debugprint = false
   )

func main() {
     flag.Parse()
     debugprint = *dbg
     Args := flag.Args()
	if *server {
		S := new(so9p.So9ps)
		S.Path = "/"
		rpc.Register(S)
		l, err := net.Listen("tcp", ":1234")
		if err != nil {
			log.Fatal(err)
		}
		rpc.Accept(l)
	} else {
		var client so9p.So9pc
		var err error
		rootfid := so9p.Fid(1)
		client.Client, err = rpc.Dial("tcp", "localhost"+":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}

		fi, err := client.Attach("/", rootfid)
		if err != nil {
			log.Fatal("attach", err)
		}
		if debugprint {
		   fmt.Printf("attach fi %v\n", fi)
		   }
		ents, err := client.ReadDir("/etc")
		if err != nil {
			log.Fatal("ReadDIr", err)
		}
		if debugprint {
			fmt.Printf("readdir %v: %v,%v\n", "/etc", ents, err)
		}

		if len(Args) < 1 {
			return
		}
		hostfid, err := client.Open(Args[0], os.O_RDONLY)
		if err != nil {
			log.Fatal("Open", err)
		}
		if debugprint {
			fmt.Printf("open %v: %v, %v\n", Args[0], hostfid, err)
		}
		for i := 1; i < 1048576; i = i * 2 {
			start := time.Now()
			data, err := client.Read(hostfid, i, 0)
			cost := time.Since(start)
			if err != nil {
				log.Fatal("read", err)
			}
			fmt.Printf("%v took %v\n", len(data), cost)

		}
		if len(Args) < 2 {
			return
		}

		copyfid, err := client.Create(Args[1], os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal("Create", err)
		}
		if debugprint {
			fmt.Printf("create %v: %v, %v\n", Args[1], hostfid, err)
		}
		for i := 1; i < 1048576; i = i * 2 {
			start := time.Now()
			data, err := client.Read(hostfid, i, 0)
			if err != nil {
				log.Fatal("read", err)
			}
			_, err = client.Write(copyfid, data, 0)
			cost := time.Since(start)
			if err != nil {
				log.Fatal("write", err)
			}
			fmt.Printf("%v took %v\n", len(data), cost)

		}
	}

}
