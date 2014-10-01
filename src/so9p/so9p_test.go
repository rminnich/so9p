// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package so9p

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"testing"
	"time"
)

var (
	dbg        = flag.Bool("dbg", false, "set true for debug")
	server     = flag.Bool("s", true, "set true for server")
	debugprint = false
)

func TestRunLocalFS(t *testing.T) {
	flag.Parse()
	debugprint = *dbg
	Args := flag.Args()
	go func() {
		S := new(So9ps)
		S.Path = "/"
		rpc.Register(S)
		l, err := net.Listen("tcp", ":1234")
		if err != nil {
			log.Fatal(err)
		}
		rpc.Accept(l)
	}()

	time.Sleep(time.Second)
	var client So9pc
	var err error
	rootfid := Fid(1)
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

func TestRAMFS(t *testing.T) {
	flag.Parse()
	debugprint = *dbg
	go func() {
		S := new(So9ps)
		S.Path = "/"
		rpc.Register(S)
		l, err := net.Listen("tcp", ":2345")
		if err != nil {
			log.Fatal(err)
		}
		rpc.Accept(l)
	}()

	time.Sleep(time.Second)
	var client So9pc
	var err error
	rootfid := Fid(1)
	client.Client, err = rpc.Dial("tcp", "localhost"+":2345")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	fi, err := client.Attach("/ramfs", rootfid)
	if err != nil {
		log.Fatal("attach", err)
	}
	if debugprint {
		fmt.Printf("attach fi %v\n", fi)
	}
	ents, err := client.ReadDir("/")
	if err != nil {
		log.Fatal("ReadDIr", err)
	}
	if debugprint {
		fmt.Printf("readdir %v: %v,%v\n", "/etc", ents, err)
	}

	hostfid, err := client.Open("x", os.O_RDONLY)
	if err == nil {
		log.Fatal("ramfs open 'x' succeeded, should have failed")
	}

	copyfid, err := client.Create("x", os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Create", err)
	}
	if debugprint {
		fmt.Printf("create %v: %v, %v\n", "x", hostfid, err)
	}
	_, err = client.Write(copyfid, []byte("Hi there"), 0)
	if err != nil {
		log.Fatal("write", err)
	}
	data, err := client.Read(hostfid, 128, 0)
	if err != nil {
		log.Fatal("read", err)
	}
	log.Printf("read ramfs x :%v:", data)

}
