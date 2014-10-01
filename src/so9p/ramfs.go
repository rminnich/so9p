package so9p

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var ramFSmap = make(map[string]*ramFSnode)

type ramFS struct {
	ramFSnode
}

type ramFSnode struct {
	File string
}

// AddRamFS adds RamFS to the set of file systems. Really needed only for primitive testing.
func AddRamFS() {
	node := &ramFSnode{}
	path2Server["/ramfs"] = node
}

func (node *ramFSnode) Create(name string, flag int, perm os.FileMode) (Node, error) {
	if DebugPrint {
		fmt.Printf("filenode.Create, node is %v\n", node)
	}
	if file, ok := ramFSmap[name]; ok {
		return file, nil
	}
	// just always replace what's there if O_CREATE is set
	if flag&os.O_CREATE == os.O_CREATE {
		fmt.Printf("filenode.Create, create node %v\n", name)
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
		fmt.Printf("FI %v\n", node)
	}
	return fi, nil
}

func (node *ramFSnode) Read(Size int, Off int64) ([]byte, error) {
	if DebugPrint {
		fmt.Printf("node %v\n", node)
	}
	b := []byte(node.File)
	if DebugPrint {
		fmt.Printf("file %v\n", node.File)
	}

	if len(b) < Size {
		Size = len(b)
	}
	return b[0:Size], nil
}

func (node *ramFSnode) Write(data []byte, Off int64) (size int, err error) {
	if DebugPrint {
		fmt.Printf("ramfs write node %v\n", node)
	}
	if DebugPrint {
		fmt.Printf("ramfs write file %v\n", node.File)
	}
	node.File = string(data)
	return size, nil
}

func (node *ramFSnode) Close() (err error) {
	if DebugPrint {
		fmt.Printf("filenode.Close node %v\n", node)
	}
	if DebugPrint {
		fmt.Printf("filenode.Close file %v\n", node.File)
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
		fmt.Printf("filenode.ReadDir node %v\n", node)
	}
	if DebugPrint {
		fmt.Printf("filenode.ReadDir file %v\n", node.File)
	}

	for v := range ramFSmap {
		fi = append(fi, FileInfo{Name: v})
	}
	return fi, nil
}
