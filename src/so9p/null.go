package so9p

// the null server does nothing, badly.

type nullFS struct {
}

type nullNode struct {
}

var (
	null = &nullNode{}
)

func (node *nullNode) FI(name string) (FileInfo, error) {
	return FileInfo{}, nil
}
