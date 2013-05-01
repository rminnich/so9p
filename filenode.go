package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func (fs *fileFS) Root() (node Node, err error) {
	node, err = &fileNode{}, nil
	return
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

func (node *fileNode) Create(name string, flag int, perm os.FileMode) (Node,error){
	if debugprint {
		fmt.Printf("filenode.Create, node is %v\n", node)
	}
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	newNode := &fileNode{File: file}
	return newNode, err
}

func (node *fileNode) Mkdir(name string, int, perm os.FileMode) error{
	err := os.Mkdir(name, perm)
	return err
}

func osFI2FI(osfi os.FileInfo, fi *FileInfo) {
	fi.SSize = osfi.Size()
	fi.SMode = osfi.Mode()
	fi.SModTime = osfi.ModTime()
	fi.SIsDir = osfi.IsDir()
}

func (node *fileNode) FI(name string) (FileInfo, error) {
	var fi FileInfo
	if debugprint {
		fmt.Printf("FI %v\n", node)
	}
	osfi, err := os.Stat(name)

	if err != nil {
		log.Print(err)
		return fi, err
	}
	osFI2FI(osfi, &fi)
	return fi, err
}

func (node *fileNode) Read(Size int, Off int64) ([]byte, error) {
	if debugprint {
		fmt.Printf("node %v\n", node)
	}
	b := make([]byte, Size)
	if debugprint {
		fmt.Printf("file %v\n", node.File)
	}

	n, err := node.File.ReadAt(b, Off)
	if debugprint {
		fmt.Printf("Read %v, %v\n", n, err)
	}
	if err != nil {
		log.Print(err)
	}
	return b[0:n], err
}

func (node *fileNode) Write(data []byte, Off int64) (size int, err error) {
	if debugprint {
		fmt.Printf("node %v\n", node)
	}
	if debugprint {
		fmt.Printf("file %v\n", node.File)
	}

	size, err = node.File.WriteAt(data, Off)
	if debugprint {
		fmt.Printf("Write %v, %v\n", size, err)
	}
	if err != nil {
		log.Print(err)
	}
	return size, err
}

func (node *fileNode) Close() (err error) {
	if debugprint {
		fmt.Printf("filenode.Close node %v\n", node)
	}
	if debugprint {
		fmt.Printf("filenode.Close file %v\n", node.File)
	}

	err = node.File.Close()
	if err != nil {
		log.Print(err)
	}
	return err
}

/* we don't even implement opendir because it never
 * made any sense. Just call ReadDir with a node
 * you walked to and we're done.
 */
func (node *fileNode) ReadDir(name string) ([]FileInfo, error) {
	if debugprint {
		fmt.Printf("filenode.ReadDir node %v\n", node)
	}
	if debugprint {
		fmt.Printf("filenode.ReadDir file %v\n", node.File)
	}

	osfi, err := ioutil.ReadDir(name)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	fi := make([]FileInfo, len(osfi))
	for i,_ := range(fi) {
		osFI2FI(osfi[i], &fi[i])
	}
	return fi, err
}
