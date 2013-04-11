package main

import (
	"fmt"
	"log"
	"path"
	"os"
	"time"
)

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

	newNode := &fileNode{walkTo, fi.Name()}
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

func (fs *fileFS) Root() (node Node, err error) {
	node, err = &fileNode{"/", "/"}, nil
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

