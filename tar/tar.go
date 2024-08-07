package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const DefaultExtractFolder = "./extract"

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

func scan[T any](r io.Reader, OnNodeCreation func(*Node[T]) error, skipData bool) (*Node[T], error) {
	tr := tar.NewReader(r)
	root := newRootNode[T]()
	if err := OnNodeCreation(root); err != nil {
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
		if !skipData {
			// Read data (Read returns directly 0,io.EOF if not TypeReg)
			if _, err := tr.Read(nf.data); err != nil && err != io.EOF {
				return nil, fmt.Errorf("on reading file: %s", err)
			}
		}

		if err := OnNodeCreation(nf); err != nil {
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

// Scan through the io reader and return the root directory node of the archive,
// Node is a generic type, you can implement it with the callback Node type eg: func(n *Node[struct{}])
// The type T is used to add additionnal data into each nodes on creation. It let the possibility to initialize each node.
func Scan[T any](r io.Reader, OnNodeCreation func(*Node[T]) error) (*Node[T], error) {
	return scan(r, OnNodeCreation, false)
}

// List through archive to extract all headers name
func List(r io.Reader) ([]string, error) {
	root, err := scan(r, func(*Node[struct{}]) error { return nil }, true)
	if err != nil {
		return nil, err
	}
	res := []string{root.GetPath()}

	root.ForAllChildren(func(nd *Node[struct{}]) error {
		res = append(res, nd.header.Name)
		return nil
	})

	return res, nil
}

// Extract all nodes to the output file.
// isSkipped callback can be used to add logic (skip current node if true) on nodes extraction
func Extract[T any](node *Node[T], output string, isSkipped func(*Node[T]) bool) error {
	if len(output) == 0 {
		output = DefaultExtractFolder
	}
	return node.ForAllChildren(func(nd *Node[T]) error {
		if isSkipped(nd) {
			return nil
		}

		dirPath := filepath.Join(output, nd.GetParent().GetPath())
		if !nd.IsDir() && nd.Mode().IsRegular() {
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				err := os.MkdirAll(dirPath, 0777) //TODO change me to use permissions from archive?
				if err != nil {
					return fmt.Errorf("error on create directory %s: %s", dirPath, err)
				}
			}
			filePath := filepath.Join(dirPath, nd.Name())
			if err := os.WriteFile(filePath, nd.GetData(), nd.Mode().Perm()); err != nil {
				return fmt.Errorf("error on create file %s: %s", filePath, err)
			}
		}
		return nil
	})
}
