package tar

import (
	"archive/tar"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newRootNode(t *testing.T) {
	t.Run("Create root Node", func(t *testing.T) {
		assert.Equal(t, ".", newRootNode[struct{}]().GetPath())
	})
}

func Test_Node(t *testing.T) {
	var root *Node[struct{}]
	name := "test_new_node/"
	t.Run("Create Node", func(t *testing.T) {
		root = newNode[struct{}](&tar.Header{Name: name})
		assert.Equal(t, root.Name(), root.GetPath())
		assert.True(t, root.IsRoot())
	})
	var child *Node[struct{}]
	t.Run("Create child Node", func(t *testing.T) {
		child = newNode[struct{}](&tar.Header{Name: name + "some_file.txt"})
		root.addChild(child)
		assert.Len(t, root.GetChildren(), 1)
		assert.False(t, child.IsRoot())
	})

	t.Run("Get Root/Parent Node", func(t *testing.T) {
		assert.Equal(t, root, child.GetRoot())
		assert.Equal(t, root, child.GetParent()) // return root
		assert.Equal(t, root, root.GetParent())  // return itself because it's the root node
	})

	t.Run("Get data from Node", func(t *testing.T) {
		td := []byte("hello there")
		root.data = td
		assert.Equal(t, td, root.GetData())
	})
}
