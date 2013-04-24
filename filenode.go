package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

func (fs *fileFS) Root() (node Node, err error) {
	node, err = &fileNode{FullPath: "/", Name: "/"}, nil
	return
}

func (fi FileInfo) FullPath() string {
	return fi.SFullPath
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
	if debugprint {
		fmt.Printf("filenode.Walk, node is %v\n", node)
	}
	walkTo = path.Clean(path.Join("/", walkTo))
	finalPath := path.Join(node.FullPath, walkTo)
	/* walk to whatever the new path is -- may be same as old */
	if debugprint {
		fmt.Printf("full %v\n", finalPath)
	}

	fi, err := os.Stat(finalPath)
	if err != nil {
		return nil, err
	}

	newNode := &fileNode{FullPath: finalPath, Name: fi.Name()}
	return newNode, err
}

func (node *fileNode) Create(name string, flag int, perm os.FileMode) (Node,error){
	/* push a / onto the front of path. Then clean it.
	 * This removes attempts to walk out of the tree.
	 */
	if debugprint {
		fmt.Printf("filenode.Create, node is %v\n", node)
	}
	name = path.Clean(path.Join("/", name))
	finalPath := path.Join(node.FullPath, name)
	/* walk to whatever the new path is -- may be same as old */
	if debugprint {
		fmt.Printf("full %v\n", finalPath)
	}

	file, err := os.OpenFile(finalPath, flag, perm)
	if err != nil {
		return nil, err
	}

	newNode := &fileNode{FullPath: finalPath, Name: name, File: file}
	return newNode, err
}

func (node *fileNode) Mkdir(name string, int, perm os.FileMode) error{
	/* push a / onto the front of path. Then clean it.
	 * This removes attempts to walk out of the tree.
	 */
	if debugprint {
		fmt.Printf("filenode.Mkdir is %v\n", node)
	}
	name = path.Clean(path.Join("/", name))
	finalPath := path.Join(node.FullPath, name)
	/* walk to whatever the new path is -- may be same as old */
	if debugprint {
		fmt.Printf("full %v\n", finalPath)
	}

	err := os.Mkdir(finalPath, perm)
	return err
}

func osFI2FI(osfi os.FileInfo, fi *FileInfo) {
	fi.SName = osfi.Name()
	fi.SSize = osfi.Size()
	fi.SMode = osfi.Mode()
	fi.SModTime = osfi.ModTime()
	fi.SIsDir = osfi.IsDir()
}

func (node *fileNode) FI() (FileInfo, error) {
	var fi FileInfo
	if debugprint {
		fmt.Printf("FI %v\n", node)
	}
	osfi, err := os.Stat(node.FullPath)

	if err != nil {
		log.Print(err)
		return fi, err
	}
	osFI2FI(osfi, &fi)
	return fi, err
}

func (node *fileNode) Open() (err error) {
	if debugprint {
		fmt.Printf("Open: node %v\n", node)
	}
	node.File, err = os.Open(node.FullPath)

	if err != nil {
		log.Print("OPEN failed: ", err)
	}
	if debugprint {
		fmt.Printf("node, %v, file, %v\n", node, node.File)
	}
	return err
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
func (node *fileNode) ReadDir() ([]FileInfo, error) {
	if debugprint {
		fmt.Printf("filenode.ReadDir node %v\n", node)
	}
	if debugprint {
		fmt.Printf("filenode.ReadDir file %v\n", node.File)
	}

	osfi, err := ioutil.ReadDir(node.FullPath)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	fi := make([]FileInfo, len(osfi))
	for i,_ := range(fi) {
		osFI2FI(osfi[i], &fi[i])
		fi[i].SFullPath = path.Join(node.FullPath, fi[i].SName)
	}
	return fi, err
}
