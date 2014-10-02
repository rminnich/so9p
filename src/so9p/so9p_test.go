// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package so9p

import (
	"bytes"
	"io"
	"net"
	"net/rpc"
	"os"
	"testing"
	"time"
)

func TestStartServer(t *testing.T) {
	go func() {
		DebugPrint = false
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

func TestRunLocalFS(t *testing.T) {
	var conn So9pConn
	var client *So9pc
	var err error
	conn.Client, err = rpc.Dial("tcp", "localhost"+":1234")
	if err != nil {
		t.Fatal("test: dialing:", err)
	}
	if client, err = conn.Attach("/"); err != nil {
		t.Fatal("test: attach", err)
	}
	defer client.Unattach()

	t.Logf("test:attach client %v\n", client)
	ents, err := client.ReadDir("/etc")
	if err != nil && err != io.EOF {
		t.Fatal("test: ReadDIr", err)
	}
	t.Logf("test: readdir %v: %v,%v\n", "/etc", len(ents), err)

	hosts, err := client.Open("/etc/hosts", os.O_RDONLY)
	if err != nil {
		t.Fatal("test: Open", err)
	}
	t.Logf("test: open %v: %v, %v\n", "/etc/hosts", hosts, err)
	data := make([]byte, 1024)
	start := time.Now()
	for i := 1; i < 1048576; i = i * 2 {
		amt, err := hosts.ReadAt(data, 0)
		if err != nil && err != io.EOF {
			t.Fatalf("test: read loop iter %v %v %v %v", i, amt, err, io.EOF)
		}

	}
	cost := time.Since(start)
	t.Logf("test: 1M iterations took %v\n", cost)

	copyto, err := client.Create("/tmp/x", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		t.Fatal("test: Create", err)
	}
	t.Logf("test: create %v: %v, %v\n", "/tmp/x", copyto, err)
	for i := 1; i < 1048576; i = i * 2 {
		start := time.Now()
		i, err := hosts.ReadAt(data, 0)
		if err != nil && err != io.EOF {
			t.Fatal("test: Read", err)
		}
		_, err = copyto.WriteAt(data, 0)
		cost := time.Since(start)
		if err != nil {
			t.Fatal("test: write", err)
		}
		t.Logf("test: %v took %v\n", i, cost)

	}

}

func TestRAMFS(t *testing.T) {
	var conn So9pConn
	var client *So9pc
	var err error
	AddRamFS()
	conn.Client, err = rpc.Dial("tcp", "localhost"+":1234")
	if err != nil {
		t.Fatal("test: dialing:", err)
	}
	if client, err = conn.Attach("/ramfs"); err != nil {
		t.Fatal("test: attach", err)
	}
	defer client.Unattach()

	ents, err := client.ReadDir("/")
	if err != nil {
		t.Fatal("test: ReadDIr", err)
	}
	t.Logf("test: readdir %v: %v,%v\n", "/etc", len(ents), err)
	t.Logf("test: attach client %v\n", client)
	_, err = client.Open("x", os.O_RDONLY)
	if err == nil {
		t.Fatal("test: ramfs open 'x' succeeded, should have failed")
	}

	copyto, err := client.Create("x", os.O_WRONLY, 0666)
	if err != nil {
		t.Fatal("test: Create", err)
	}
	t.Logf("test: create %v: %v\n", "x", copyto)
	writedata := []byte("Hi there")
	readdata := writedata
	_, err = copyto.WriteAt(writedata, 0)
	if err != nil {
		t.Fatal("test: write", err)
	}
	if false {
	_, err = copyto.ReadAt(readdata, 0)
	if err != nil {
		t.Fatal("test: read", err)
	}
	t.Logf("test: read ramfs x :%v:", readdata)
	if !bytes.Equal(writedata[:], readdata[:]) {
		t.Fatal("test: writedata and readdata did not match: '%v' != '%v'", writedata, readdata)
	}

	ents, err = client.ReadDir("/")
	if err != nil {
		t.Fatal("test: ReadDIr", err)
	}
	t.Logf("test: readdir %v: %v,%v\n", "/etc", len(ents), err)
	}

}

func TestIO(t *testing.T) {
	var conn So9pConn
	var source, dest *So9pc
	var err error
	AddRamFS()
	conn.Client, err = rpc.Dial("tcp", "localhost"+":1234")
	if err != nil {
		t.Fatal("test: dialing:", err)
	}
	if dest, err = conn.Attach("/ramfs"); err != nil {
		t.Fatal("test: attach ramfs", err)
	}
	if source, err = conn.Attach("/"); err != nil {
		t.Fatal("test: local attach", err)
	}
	fs, err := source.Open("/etc/hosts", os.O_RDONLY)
	if err != nil {
		t.Fatal("test: Open", err)
	}
	fd, err := dest.Open("/x", os.O_RDONLY)
	if err != nil {
		t.Fatal("test: Open", err)
	}
	total, err := io.Copy(fd, fs)
	if err == nil {
		t.Fatal("test: io.Copy did not fail but should have")
	}
	fd, err = dest.Create("/hosts", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		t.Fatal("test: Create", err)
	}
	total, err = io.Copy(fd, fs)
	t.Logf("test: Copied %v bytes", total)
}
