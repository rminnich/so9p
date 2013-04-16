package main

import (
	"fmt"
	"log"
	"path"
	"os"
	"time"
)
func (fs *fileFS) Root() (node Node, err error) {
	node, err = &fileNode{FullPath:"/", Name:"/"}, nil
	return
}


func (fi FileInfo) Name() string {
	return fi.SName
}

func (fi FileInfo) Size() int64 {
	return fi.SSize
}

func (fi FileInfo) Mode() os.FileMode {
	return fi.SMode
}

func (fi FileInfo) ModTime() time.Time {
	return fi.SModTime
}

func (fi FileInfo) IsDir() bool {
	return fi.SIsDir
}

func (fi FileInfo) Sys() interface{} {
	return nil
}


func (node *fileNode) Walk(walkTo string) (Node, error) {
	/* push a / onto the front of path. Then clean it.
	 * This removes attempts to walk out of the tree.
	 */
	walkTo = path.Clean(path.Join("/", walkTo))
	/* walk to whatever the new path is -- may be same as old */
	fi, err := os.Stat(path.Join(node.Name, walkTo))
	if err != nil {
		return nil, err
	}

	newNode := &fileNode{FullPath:walkTo, Name:fi.Name()}
	return newNode, err
}

func (node *fileNode) FI() (FileInfo, error) {
	var fi FileInfo
	fmt.Printf("FI %v\n", node)
	osfi, err := os.Stat(node.FullPath)

	if err != nil {
		log.Print(err)
		return fi, err
	}
	fi.SName = osfi.Name()
	fi.SSize = osfi.Size()
	fi.SMode = osfi.Mode()
	fi.SModTime = osfi.ModTime()
	fi.SIsDir = osfi.IsDir()
	return fi, err
}

func (node *fileNode) Open() (err error) {
	fmt.Printf("node %v\n", node)
	node.File, err = os.Open(node.FullPath)

	if err != nil {
		log.Print(err)
	}
	fmt.Printf("node, %v, file, %v\n", node, node.File)
	return err
}

func (node *fileNode) Read(Size int, Off int64) (b []byte, err error) {
	fmt.Printf("node %v\n", node)
	b = make([]byte, Size)
	fmt.Printf("file %v\n", node.File)
	
	n, err := node.File.ReadAt(b, Off)
	fmt.Printf("Read %v, %v\n", n, err)
	if err != nil {
		log.Print(err)
	}
	return b, err
}

func (node *fileNode) Write(data[]byte, Off int64) (size int, err error) {
	fmt.Printf("node %v\n", node)
	fmt.Printf("file %v\n", node.File)
	
	size, err = node.File.WriteAt(data, Off)
	fmt.Printf("Write %v, %v\n", size, err)
	if err != nil {
		log.Print(err)
	}
	return size, err
}

func (node *fileNode) Close() (err error) {
	fmt.Printf("node %v\n", node)
	fmt.Printf("file %v\n", node.File)
	
	err = node.File.Close()
	fmt.Printf("Close %v\n", err)
	if err != nil {
		log.Print(err)
	}
	return err
}

