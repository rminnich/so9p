package so9p

import (
	"bytes"
	"errors"
	"log"
	"os"
)

var ramFSmap = make(map[string]*ramFSNode)

type ramFS struct {
	ramFSNode
}

type ramFSNode struct {
	b bytes.Buffer
}

// Attach implements a server attach for local file nodes
func (n *ramFSNode) Attach(Args *AttachArgs, Resp *Attachresp) (err error) {
	Resp.FI, err = n.FI(Args.Name)
	if err != nil {
		log.Printf("FI fails for %v\n", Args.Name)
	}
	return err
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
		ramFSmap[name] = &ramFSNode{}
		return ramFSmap[name], nil
	}

	return nil, errors.New("ramfs: nope")
}

// Mkdir implements os.Mkdir
func (n *ramFSNode) Mkdir(name string, int, perm os.FileMode) error {
	return errors.New("ramfs: mkdir: nope")
}

// FI returns an os.FileInnfo
func (n *ramFSNode) FI(name string) (FileInfo, error) {
	var fi FileInfo
	if debugPrint {
		log.Printf("FI %v\n", n)
	}
	return fi, nil
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
func (n *ramFSNode) ReadDir(name string) ([]FileInfo, error) {
	var fi []FileInfo
	if debugPrint {
		log.Printf("filenode.ReadDir node %v\n", n)
	}

	for v := range ramFSmap {
		fi = append(fi, FileInfo{Name: v})
	}
	return fi, nil
}
