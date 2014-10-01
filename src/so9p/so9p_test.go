// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package so9p

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"testing"
	"time"
)

func TestStartServer(t *testing.T) {
	go func() {
		DebugPrint = true
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
}

func TestRunLocalFS(t *testing.T) {
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
	fmt.Printf("attach fi %v\n", fi)
	ents, err := client.ReadDir("/etc")
	if err != nil {
		log.Fatal("ReadDIr", err)
	}
	fmt.Printf("readdir %v: %v,%v\n", "/etc", ents, err)

	hostfid, err := client.Open("/etc/hosts", os.O_RDONLY)
	if err != nil {
		log.Fatal("Open", err)
	}
	fmt.Printf("open %v: %v, %v\n", "/etc/hosts", hostfid, err)
	for i := 1; i < 1048576; i = i * 2 {
		start := time.Now()
		data, err := client.Read(hostfid, i, 0)
		cost := time.Since(start)
		if err != nil {
			log.Fatal("read", err)
		}
		fmt.Printf("%v took %v\n", len(data), cost)

	}

	copyfid, err := client.Create("/tmp/x", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Create", err)
	}
	fmt.Printf("create %v: %v, %v\n", "/tmp/x", hostfid, err)
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

	time.Sleep(time.Second)
	var client So9pc
	var err error
	rootfid := Fid(1)
	client.Client, err = rpc.Dial("tcp", "localhost"+":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	fi, err := client.Attach("/ramfs", rootfid)
	if err != nil {
		log.Fatal("attach", err)
	}
	fmt.Printf("attach fi %v\n", fi)
	ents, err := client.ReadDir("/")
	if err != nil {
		log.Fatal("ReadDIr", err)
	}
	fmt.Printf("readdir %v: %v,%v\n", "/etc", ents, err)

	_, err = client.Open("x", os.O_RDONLY)
	if err == nil {
		log.Fatal("ramfs open 'x' succeeded, should have failed")
	}

	copyfid, err := client.Create("x", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Create", err)
	}
	fmt.Printf("create %v: %v\n", "x", copyfid)
	_, err = client.Write(copyfid, []byte("Hi there"), 0)
	if err != nil {
		log.Fatal("write", err)
	}
	data, err := client.Read(copyfid, 128, 0)
	if err != nil {
		log.Fatal("read", err)
	}
	log.Printf("read ramfs x :%v:", data)

}
