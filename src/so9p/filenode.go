package so9p

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"syscall"
)

type fileFS struct {
	LocalFileNode
}

type LocalFileNode struct {
	File *os.File
}

func (node *LocalFileNode) Create(name string, flag int, perm os.FileMode) (Node, error) {
	if DebugPrint {
		log.Printf("server: filenode.Create")
	}
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	newNode := &LocalFileNode{File: file}
	return newNode, err
}

func (node *LocalFileNode) Mkdir(name string, int, perm os.FileMode) error {
	err := os.Mkdir(name, perm)
	return err
}

func (node *LocalFileNode) FI(name string) (FileInfo, error) {
	var fi FileInfo
	if DebugPrint {
		log.Printf("server: FI %v\n", node)
	}
	err := syscall.Lstat(name, &fi.Stat)

	if err != nil {
		log.Print(err)
		return fi, err
	}

	fi.Link, _ = os.Readlink(name)

	if DebugPrint {
		log.Printf("server: FileInfo %v\n", fi)
	}
	fi.Name = name
	return fi, err
}

func (node *LocalFileNode) ReadAt(b []byte, Off int64) (int, error) {
	if DebugPrint {
		log.Printf("server: node %v\n", node)
	}
	if DebugPrint {
		log.Printf("server: file %v\n", node.File)
	}

	n, err := node.File.ReadAt(b, Off)
	if DebugPrint {
		log.Printf("server: Read %v, %v\n", n, err)
	}
	return n, err
}

func (node *LocalFileNode) Write(data []byte, Off int64) (size int, err error) {
	if DebugPrint {
		log.Printf("server: node %v\n", node)
	}
	if DebugPrint {
		log.Printf("server: file %v\n", node.File)
	}

	size, err = node.File.WriteAt(data, Off)
	if DebugPrint {
		log.Printf("server: Write %v, %v\n", size, err)
	}
	if err != nil {
		log.Print(err)
	}
	return size, err
}

func (node *LocalFileNode) Close() (err error) {
	if DebugPrint {
		log.Printf("server: filenode.Close node %v\n", node)
	}
	if DebugPrint {
		log.Printf("server: filenode.Close file %v\n", node.File)
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
 * What should we do if there are too many entries? Interesting question.
 */
func (node *LocalFileNode) ReadDir(name string) ([]FileInfo, error) {
	if DebugPrint {
		log.Printf("server: filenode.ReadDir node %v\n", node)
	}
	if DebugPrint {
		log.Printf("server: filenode.ReadDir file %v\n", node.File)
	}

	osfi, err := ioutil.ReadDir(name)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	fi := make([]FileInfo, len(osfi))
	for i, _ := range fi {
		fi[i].Name = osfi[i].Name()
		fullpath := path.Join(name, fi[i].Name)
		// Interesting problem. What if this one stat fails, and all others
		// succeed? In most Unix-like systems, the readdir will show the
		// name, and the stat will return as garbage. Not returning any results
		// because one Lstat failed is wrong; hiding the name because the Lstat
		// failed is wrong. If we log an error for every busted dirent we'll be
		// doing a LOT of logging. Conclusion: for now, ignore the error.
		_ = syscall.Lstat(fullpath, &fi[i].Stat)
	}
	return fi, err
}

func (node *LocalFileNode) Readlink(name string) (val string, err error) {
	if DebugPrint {
		log.Printf("server: filenode.Readlink node %v name %v\n", node, name)
	}

	val, err = os.Readlink(name)

	if err != nil {
		log.Print(err)
	}
	return val, err
}
