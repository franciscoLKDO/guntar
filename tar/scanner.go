package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"path/filepath"
)

// Find Node by path starting from n
func findNodeByPath[T any](n *Node[T], path string) *Node[T] {
	queue := []*Node[T]{n}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current.GetPath() == path {
			return current
		}
		queue = append(queue, current.GetChildren()...)
	}
	return nil
}

func getParentName[T any](n *Node[T]) string {
	return filepath.Dir(n.GetPath())
}

// Scan through the io reader and return the root directory node of the archive,
// Node is a generic type, you can implement it with the callback Node type eg: func(n *Node[struct{}])
// The type T is used to add additionnal data into each nodes on creation. It let the possibility to initialize each node.
func Scan[T any](r io.Reader, cb func(*Node[T]) error) (*Node[T], error) {
	tr := tar.NewReader(r)
	root := newRootNode[T]()
	if err := cb(root); err != nil {
		return nil, fmt.Errorf("on node creation: %s", err)
	}
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return root, nil
			}
			return nil, fmt.Errorf("on reading archive %s", err)
		}

		nf := newNode[T](header)

		// Read data (Read returns directly 0,io.EOF if not TypeReg)
		if _, err := tr.Read(nf.data); err != nil && err != io.EOF {
			return nil, fmt.Errorf("on reading file: %s", err)
		}
		if err := cb(nf); err != nil {
			return nil, fmt.Errorf("on node creation: %s", err)
		}

		// If node have same path as root, set it as root and continue
		if nf.isEqual(root) {
			root = nf
			continue
		}

		// Save Node in parent.children if parent exist, else parent must be root
		if pfile := findNodeByPath(root, getParentName(nf)); pfile != nil {
			pfile.addChild(nf)
		} else {
			root.addChild(nf)
		}
	}
}
