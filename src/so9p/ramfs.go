package so9p

import (
	"bytes"
	"errors"
	"log"
	"os"
)

var ramFSmap = make(map[string]*ramFSnode)

type ramFS struct {
	ramFSnode
}

type ramFSnode struct {
	b bytes.Buffer
}

// AddRamFS adds RamFS to the set of file systems. Really needed only for primitive testing.
func AddRamFS() {
	node := &ramFSnode{}
	path2Server["/ramfs"] = node
}

func (node *ramFSnode) Create(name string, flag int, perm os.FileMode) (Node, error) {
	if DebugPrint {
		log.Printf("filenode.Create, node is %v\n", node)
	}
	if file, ok := ramFSmap[name]; ok {
		return file, nil
	}
	// just always replace what's there if O_CREATE is set
	if flag&os.O_CREATE == os.O_CREATE {
		ramFSmap[name] = &ramFSnode{}
		return ramFSmap[name], nil
	}

	return nil, errors.New("ramfs: nope")
}

func (node *ramFSnode) Mkdir(name string, int, perm os.FileMode) error {
	return errors.New("ramfs: mkdir: nope")
}

func (node *ramFSnode) FI(name string) (FileInfo, error) {
	var fi FileInfo
	if DebugPrint {
		log.Printf("FI %v\n", node)
	}
	return fi, nil
}

// but for now we ignore offset.
func (node *ramFSnode) ReadAt(b []byte, Off int64) (int, error) {
	node.b.Reset()
	return node.b.Read(b)
}

func (node *ramFSnode) Write(data []byte, Off int64) (size int, err error) {
	return node.b.Write(data)
}

func (node *ramFSnode) Close() (err error) {
	if DebugPrint {
		log.Printf("filenode.Close node %v\n", node)
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
func (node *ramFSnode) ReadDir(name string) ([]FileInfo, error) {
	var fi []FileInfo
	if DebugPrint {
		log.Printf("filenode.ReadDir node %v\n", node)
	}

	for v := range ramFSmap {
		fi = append(fi, FileInfo{Name: v})
	}
	return fi, nil
}
