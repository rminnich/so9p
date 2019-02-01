package so9p

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"syscall"
)

type fileFS struct {
	localFileNode
}

type localFileNode struct {
	File *os.File
}

// Attach implements a server attach for local file nodes
func (n *localFileNode) Attach(Args *AttachArgs, Resp *Attachresp) (err error) {
	Resp.FI, err = n.FI(Args.Name)
	if err != nil {
		log.Printf("FI fails for %v\n", Args.Name)
	}
	return err
}

// Create implements Create for local file nodes.
func (n *localFileNode) Create(name string, flag int, perm os.FileMode) (Node, error) {
	debugPrintf("filenode: Create(%s, %#x, %o", name, flag, perm)
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	newNode := &localFileNode{File: file}
	return newNode, err
}

// Mkdir implements os.Mkdir
func (n *localFileNode) Mkdir(name string, int, perm os.FileMode) error {
	err := os.Mkdir(name, perm)
	return err
}

// FI returns an os.FileInnfo
func (n *localFileNode) FI(name string) (FileInfo, error) {
	var fi FileInfo
	if debugPrint {
		log.Printf("server: FI %v\n", n)
	}
	err := syscall.Lstat(name, &fi.Stat)

	if err != nil {
		log.Print(err)
		return fi, err
	}

	fi.Link, _ = os.Readlink(name)

	if debugPrint {
		log.Printf("server: FileInfo %v\n", fi)
	}
	fi.Name = name
	return fi, err
}

// ReadAt implements ReadAt
func (n *localFileNode) ReadAt(b []byte, Off int64) (int, error) {
	if debugPrint {
		log.Printf("server: node %v\n", n)
	}
	if debugPrint {
		log.Printf("server: file %v\n", n.File)
	}

	amt, err := n.File.ReadAt(b, Off)
	if debugPrint {
		log.Printf("server: Read %v, %v\n", amt, err)
	}
	return amt, err
}

// Write implements os.Write
func (n *localFileNode) Write(data []byte, Off int64) (size int, err error) {
	if debugPrint {
		log.Printf("server: node %v\n", n)
	}
	if debugPrint {
		log.Printf("server: file %v\n", n.File)
	}

	size, err = n.File.WriteAt(data, Off)
	if debugPrint {
		log.Printf("server: Write %v, %v\n", size, err)
	}
	if err != nil {
		log.Print(err)
	}
	return size, err
}

// Close implements os.Close
func (n *localFileNode) Close() (err error) {
	if debugPrint {
		log.Printf("server: filenode.Close node %v\n", n)
	}
	if debugPrint {
		log.Printf("server: filen.Close file %v\n", n.File)
	}

	err = n.File.Close()
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

// ReadDir implements os.Readdir
func (n *localFileNode) ReadDir(name string) ([]FileInfo, error) {
	if debugPrint {
		log.Printf("server: filenode.ReadDir node %v\n", n)
	}
	if debugPrint {
		log.Printf("server: filenode.ReadDir file %v\n", n.File)
	}

	osfi, err := ioutil.ReadDir(name)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	fi := make([]FileInfo, len(osfi))
	for i := range fi {
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

// ReadLink implements ReadLink
func (n *localFileNode) Readlink(name string) (val string, err error) {
	if debugPrint {
		log.Printf("server: filenode.Readlink node %v name %v\n", n, name)
	}

	val, err = os.Readlink(name)

	if err != nil {
		log.Print(err)
	}
	return val, err
}
