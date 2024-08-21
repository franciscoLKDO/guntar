package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const ExtractFolder = "guntar_extracted"

func scan[T any](r io.Reader, OnNodeCreation func(*Node[T]) error, readData bool) (*Node[T], error) {
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

		nf, err := root.addChildFromHeader(header)
		if _, ok := err.(NodeExistError); ok { // If path already exist, stop and continue
			continue
		}
		if readData {
			// Read data (Read returns directly 0,io.EOF if not TypeReg)
			if _, err := tr.Read(nf.data); err != nil && err != io.EOF {
				return nil, fmt.Errorf("on reading file: %s", err)
			}
		}

		if err := OnNodeCreation(nf); err != nil {
			return nil, fmt.Errorf("on node creation: %s", err)
		}
	}
}

// Scan through a reader (file,string,etc...) with a tar archive and return the root directory node of the archive,
// Node is a generic type, you can implement it with the callback Node type eg: func(n *Node[struct{}])
// The type T is used to add additionnal data into each nodes on creation. It let the possibility to initialize each node.
func Scan[T any](r io.Reader, OnNodeCreation func(*Node[T]) error) (*Node[T], error) {
	return scan(r, OnNodeCreation, true)
}

// List through archive to extract all headers name
func List(r io.Reader) ([]string, error) {
	root, err := scan(r, func(*Node[struct{}]) error { return nil }, false)
	if err != nil {
		return nil, err
	}

	list := []string{}
	root.OnNestedChildren(func(nd *Node[struct{}]) error {
		list = append(list, nd.GetPath())
		return nil
	})

	return list, nil
}

// Extract all nodes to the output file.
// isSkipped callback can be used to add logic (skip current node if true) on nodes extraction
func Extract[T any](node *Node[T], outputPath string, isSkipped func(*Node[T]) bool) error {
	if len(outputPath) == 0 {
		var err error
		outputPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("error on get current directory: %s", err)
		}
	}
	outputPath = filepath.Join(outputPath, ExtractFolder)
	if err := os.Mkdir(outputPath, 0777); os.IsExist(err) {
		return fmt.Errorf("error on create extract directory %s: %s", outputPath, err)
	}
	return node.OnNestedChildren(func(nd *Node[T]) error {
		if isSkipped(nd) {
			return nil
		}

		dirPath := filepath.Join(outputPath, nd.GetParent().GetPath())
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
