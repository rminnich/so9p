package so9p

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type ramFS struct {
	path string
}

type ramFSNode struct {
	name string
	b    bytes.Buffer
	FileInfo
}

var ramFSmap = make(map[string]*ramFSNode)

// Attach implements a server attach for local file nodes
func (f *ramFS) Attach(p ...string) (Node, error) {
	node := &ramFSNode{name: filepath.Join(f.path, filepath.Join(p...))}
	return node, nil
}

// Create implements Create for local file nodes.
func (n *ramFSNode) Create(name string, flag int, perm os.FileMode) (Node, error) {
	if debugPrint {
		log.Printf("ramfsnode.Create, node is %v\n", n)
	}
	if file, ok := ramFSmap[name]; ok {
		return file, nil
	}
	// just always replace what's there if O_CREATE is set
	if flag&os.O_CREATE == os.O_CREATE {
		debugPrintf("ramfs: create %s", name)
		ramFSmap[name] = &ramFSNode{}
		return ramFSmap[name], nil
	}

	return nil, errors.New("ramfs: nope")
}

// Mkdir implements os.Mkdir
func (n *ramFSNode) Mkdir(name string, int, perm os.FileMode) error {
	return errors.New("ramfs: mkdir: nope")
}

// Stat returns a FileInfo
func (n *ramFSNode) Stat() (*FileInfo, error) {
	fi := &FileInfo{}
	if debugPrint {
		log.Printf("Stat %v\n", n)
	}
	return fi, nil
}

func (n *ramFSNode) Sys() FileInfo {
	return n.FileInfo
}

// ReadAt implements ReadAt
// but for now we ignore offset.
func (n *ramFSNode) ReadAt(b []byte, Off int64) (int, error) {
	n.b.Reset()
	return n.b.Read(b)
}

// Write implements os.Write
func (n *ramFSNode) Write(data []byte, Off int64) (size int, err error) {
	return n.b.Write(data)
}

// Close implements os.Close
func (n *ramFSNode) Close() (err error) {
	if debugPrint {
		log.Printf("filenode.Close node %v\n", n)
	}

	if err != nil {
		log.Print(err)
	}
	return nil
}

/* we don't even implement opendir because it never
 * made any sense. Just call ReadDir with a node
 * you walked to and we're done.
 * What should we do if there are too many entries? Interesting question.
 */

// ReadDir implements os.Readdir
func (n *ramFSNode) Readdir() ([]Node, error) {
	return nil, fmt.Errorf("not now")

	var fi []FileInfo
	if debugPrint {
		log.Printf("filenode.ReadDir node %v\n", n)
	}

	for v := range ramFSmap {
		fi = append(fi, FileInfo{Name: v})
	}
	return nil, nil
}
