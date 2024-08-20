package tar

import (
	"archive/tar"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"time"
)

type SimpleNode = Node[struct{}]

type Node[T any] struct {
	fs.FileInfo
	header   *tar.Header
	path     string     // path is the unique id of the node
	parent   *Node[T]   // parent is the direct parent of actual node, nil if node is root
	children []*Node[T] // children are all childs under this node
	data     []byte     // data is the content of file, empty if not a file
	Spec     T          // Spec is the additionalData that users can set on node creation
}

func (n Node[T]) GetPath() string         { return n.path }               // Id and full path of Node
func (n Node[T]) GetChildren() []*Node[T] { return n.children }           // GetChildren return the node's children
func (n Node[T]) LenChildren() int        { return len(n.GetChildren()) } // Get size children of current Node
func (n Node[T]) IsRoot() bool            { return n.parent == nil }      // Node is root if no parents
func (n Node[T]) GetData() []byte         { return n.data }               // Get data (used for files, other are empty)

// Get Root Node from current node
func (n *Node[T]) GetRoot() *Node[T] {
	if n.IsRoot() {
		return n
	} else {
		return n.parent.GetRoot()
	}
}

// GetParent return the parent node, or itself if it's root
func (n *Node[T]) GetParent() *Node[T] {
	if n.IsRoot() {
		return n
	}
	return n.parent
}

// GetDirChildren return the node's directory children
func (n Node[T]) GetDirChildren() []*Node[T] {
	var res []*Node[T]
	for _, node := range n.GetChildren() {
		if node.IsDir() {
			res = append(res, node)
		}
	}
	return res
}

// ForAllChildren run cb in all nested children
func (n *Node[T]) ForAllChildren(cb func(*Node[T]) error) error {
	for _, node := range n.GetChildren() {
		if err := cb(node); err != nil {
			return err
		}
		if err := node.ForAllChildren(cb); err != nil {
			return fmt.Errorf("error on loop into children for %s: %s", n.GetPath(), err)
		}
	}
	return nil
}

// addChild set this node as parent to the added node and append it to children
func (n *Node[T]) addChild(node *Node[T]) {
	node.parent = n
	n.children = append(n.children, node)
}

func (n *Node[T]) addChildFromHeader(header *tar.Header) (*Node[T], error) {
	path := filepath.Join(n.GetPath(), header.Name)
	nd := newNode[T](header, path)

	queue := []*Node[T]{n}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current.GetPath() == path { // Return an error if path already exist
			return nil, NodeExistError{path: path}
		}
		if current.GetPath() == filepath.Dir(path) { // link node to parent
			current.addChild(nd)
			return nd, nil
		}
		queue = append(queue, current.GetDirChildren()...)
	}
	n.addChild(nd)
	return nd, nil
}

func newNode[T any](header *tar.Header, path string) *Node[T] {
	return &Node[T]{
		header:   header,
		FileInfo: header.FileInfo(),
		path:     path,
		data:     make([]byte, header.Size),
	}
}

type rootFI struct {
	name    string
	modTime time.Time
}

func (d *rootFI) Name() string       { return path.Base(path.Clean(d.name)) }
func (d *rootFI) Size() int64        { return 0 }
func (d *rootFI) Mode() fs.FileMode  { return fs.ModeDir }
func (d *rootFI) ModTime() time.Time { return d.modTime }
func (d *rootFI) IsDir() bool        { return true }
func (d *rootFI) Sys() interface{}   { return nil }

func newRootNode[T any]() *Node[T] {
	return &Node[T]{
		FileInfo: &rootFI{modTime: time.Now(), name: "/"},
		path:     "/",
	}
}

type NodeExistError struct {
	path string
}

func (n NodeExistError) Error() string {
	return fmt.Sprintf("Path already exist in the tree: %s", n.path)
}
