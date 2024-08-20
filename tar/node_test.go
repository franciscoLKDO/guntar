package tar

import (
	"archive/tar"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newRootNode(t *testing.T) {
	t.Run("Create root Node", func(t *testing.T) {
		assert.Equal(t, "/", newRootNode[struct{}]().GetPath())
	})
}

func Test_Node(t *testing.T) {
	var root *SimpleNode
	name := "test_new_node/"
	t.Run("Create Node", func(t *testing.T) {
		root = newRootNode[struct{}]()
		assert.Equal(t, root.Name(), root.GetPath())
		assert.True(t, root.IsRoot())
	})
	var child *Node[struct{}]
	t.Run("Create child Node", func(t *testing.T) {
		var err error
		child, err = root.addChildFromHeader(&tar.Header{Name: name})
		require.Nil(t, err)
		assert.Len(t, root.GetChildren(), 1)
		assert.False(t, child.IsRoot())
		assert.Equal(t, filepath.Join(root.GetPath(), child.GetPath()), child.GetPath())
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
