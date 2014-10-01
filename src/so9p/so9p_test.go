// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package so9p

import (
	"bytes"
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
			t.Fatal(err)
		}
		rpc.Accept(l)
	}()
	time.Sleep(time.Second)
}

func TestBadFid(t *testing.T) {
	var conn So9pConn
	var client *So9pc
	var err error
	conn.Client, err = rpc.Dial("tcp", "localhost"+":1234")
	if err != nil {
		t.Fatal("dialing:", err)
	}
	if client, err = conn.Attach("/", 23); err != nil {
		t.Fatal("attach", err)
	}
	log.Printf("client is %v", client)
	// try to attach twice
	if _, err = client.Attach("/", 23); err == nil {
		t.Fatal("attach should have failed", err)
	}
}

func TestRunLocalFS(t *testing.T) {
	var conn So9pConn
	var client *So9pc
	var err error
	conn.Client, err = rpc.Dial("tcp", "localhost"+":1234")
	if err != nil {
		t.Fatal("dialing:", err)
	}
	if client, err = conn.Attach("/", 1); err != nil {
		t.Fatal("attach", err)
	}
	t.Logf("attach client %v\n", client)
	ents, err := client.ReadDir("/etc")
	if err != nil {
		t.Fatal("ReadDIr", err)
	}
	t.Logf("readdir %v: %v,%v\n", "/etc", ents, err)

	hostfid, err := client.Open("/etc/hosts", os.O_RDONLY)
	if err != nil {
		t.Fatal("Open", err)
	}
	t.Logf("open %v: %v, %v\n", "/etc/hosts", hostfid, err)
	data := make([]byte, 1024)
	for i := 1; i < 1048576; i = i * 2 {
		start := time.Now()
		i, err := hostfid.ReadAt(data, 0)
		cost := time.Since(start)
		if err != nil {
			t.Fatal("read", err)
		}
		t.Logf("%v took %v\n", i, cost)

	}

	copyfid, err := client.Create("/tmp/x", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		t.Fatal("Create", err)
	}
	t.Logf("create %v: %v, %v\n", "/tmp/x", hostfid, err)
	for i := 1; i < 1048576; i = i * 2 {
		start := time.Now()
		i, err := hostfid.ReadAt(data, 0)
		if err != nil {
			t.Fatal("read", err)
		}
		_, err = copyfid.WriteAt(data, 0)
		cost := time.Since(start)
		if err != nil {
			t.Fatal("write", err)
		}
		t.Logf("%v took %v\n", i, cost)

	}

}

func TestRAMFS(t *testing.T) {
	var conn So9pConn
	var client *So9pc
	var err error
	AddRamFS()
	conn.Client, err = rpc.Dial("tcp", "localhost"+":1234")
	if err != nil {
		t.Fatal("dialing:", err)
	}
	if client, err = conn.Attach("/ramfs", 4444); err != nil {
		t.Fatal("attach", err)
	}
	ents, err := client.ReadDir("/")
	if err != nil {
		t.Fatal("ReadDIr", err)
	}
	t.Logf("readdir %v: %v,%v\n", "/etc", ents, err)
	t.Logf("attach client %v\n", client)
	_, err = client.Open("x", os.O_RDONLY)
	if err == nil {
		t.Fatal("ramfs open 'x' succeeded, should have failed")
	}

	copyfid, err := client.Create("x", os.O_WRONLY, 0666)
	if err != nil {
		t.Fatal("Create", err)
	}
	t.Logf("create %v: %v\n", "x", copyfid)
	writedata := []byte("Hi there")
	readdata := writedata
	_, err = copyfid.WriteAt(writedata, 0)
	if err != nil {
		t.Fatal("write", err)
	}
	_, err = copyfid.ReadAt(readdata, 0)
	if err != nil {
		t.Fatal("read", err)
	}
	log.Printf("read ramfs x :%v:", readdata)
	if !bytes.Equal(writedata[:], readdata[:]) {
		t.Fatal("writedata and readdata did not match: '%v' != '%v'", writedata, readdata)
	}

	ents, err = client.ReadDir("/")
	if err != nil {
		t.Fatal("ReadDIr", err)
	}
	t.Logf("readdir %v: %v,%v\n", "/etc", ents, err)

}
