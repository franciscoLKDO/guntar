package tar

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/franciscolkdo/guntar/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationsOnArchive(t *testing.T) {
	files := []test.File{
		{Name: "./test/", Mode: fs.ModeDir, Body: ""},
		{Name: "./test/readme.txt", Mode: 0600, Body: "This archive contains some text files."},
		{Name: "./test/hello.txt", Mode: 0600, Body: "world"},
		{Name: "gopher.txt", Mode: 0600, Body: "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{Name: "todo.txt", Mode: 0600, Body: "Get animal handling license."},
	}
	type addStruct struct {
		first  int
		second string
	}
	tests := []struct {
		name             string
		files            []test.File
		wantErr          bool
		spec             any
		expectedChildren int
	}{
		{
			name:  "archive with files and directories",
			files: files,
			spec: addStruct{
				first:  1,
				second: "two",
			},
			expectedChildren: 3,
		},
		{
			name: "archive with files only",
			files: []test.File{
				{Name: "gopher.txt", Mode: 0600, Body: "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
				{Name: "todo.txt", Mode: 0600, Body: "Get animal handling license."},
			},
			spec:             "blop",
			expectedChildren: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := Scan(test.CreateArchive(t, tt.files), func(n *Node[any]) error {
				n.Spec = tt.spec
				return nil
			})
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedChildren, root.LenChildren())
			for _, nd := range root.children {
				assert.Equal(t, tt.spec, nd.Spec)
			}
		})
	}

	t.Run("List archive", func(t *testing.T) {
		res, err := List(test.CreateArchive(t, files))
		require.Nil(t, err)
		assert.Len(t, res, len(files))
	})
	t.Run("Extract archive", func(t *testing.T) {
		root, err := Scan(test.CreateArchive(t, files), func(n *SimpleNode) error {
			return nil
		})
		require.Nil(t, err)
		// Create tmpDir to extract archive
		tmpDir, err := os.MkdirTemp("", "*")
		require.Nil(t, err)
		defer os.RemoveAll(tmpDir)

		err = Extract(root, tmpDir, func(n *SimpleNode) bool { return false })
		require.Nil(t, err)

		// Check all files and directory has been created in tmpDir
		root.ForAllChildren(func(n *SimpleNode) error {
			if n.IsDir() {
				assert.DirExists(t, filepath.Join(tmpDir, n.GetPath()))
			} else {
				assert.FileExists(t, filepath.Join(tmpDir, n.GetPath()))
			}
			return nil
		})
	})
}

func TestScanWithArchiveFile(t *testing.T) {
	type addStruct struct {
		first  int
		second string
	}
	tests := []struct {
		name    string
		wantErr bool
		spec    any
	}{
		{
			name:    "simple node with struct",
			wantErr: false,
			spec: addStruct{
				first:  1,
				second: "two",
			},
		},
		{
			name:    "simple node with string",
			wantErr: false,
			spec:    "hello there",
		},
		{
			name:    "simple node with int",
			wantErr: false,
			spec:    1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pwd, err := os.Getwd()
			require.Nil(t, err)
			file, err := os.Open(filepath.Join(pwd, "../test/mytarfolder.tar"))
			require.Nil(t, err)
			root, err := Scan(file, func(n *Node[any]) error {
				n.Spec = tt.spec
				return nil
			})
			assert.Nil(t, err)
			assert.Greater(t, root.LenChildren(), 0)
			for _, nd := range root.children {
				assert.Equal(t, tt.spec, nd.Spec)
			}
			require.Nil(t, file.Close())
		})
	}
}

func TestPathTraversal(t *testing.T) {
	fileName := "exploit_test.txt"
	files := []test.File{
		{Name: "./test/", Mode: fs.ModeDir, Body: ""},
		{Name: "./test/..", Mode: fs.ModeDir, Body: ""},
		{Name: "./test/../..", Mode: fs.ModeDir, Body: ""},
		{Name: "./test/../../..", Mode: fs.ModeDir, Body: ""},
		{Name: "./test/../../../" + fileName, Mode: 0600, Body: "tryin to get out!!!"},
	}
	var root *SimpleNode
	t.Run("Avoid path traversal", func(t *testing.T) {
		expectedPath := "/" + fileName
		var err error
		root, err = Scan(test.CreateArchive(t, files), func(n *SimpleNode) error { return nil })
		require.Nil(t, err)
		// exploit_test is on the root directory of the tree
		assert.Equal(t, expectedPath, root.children[root.LenChildren()-1].GetPath())

		// Create tmpDir to extract archive
		tmpDir, err := os.MkdirTemp("", "*")
		require.Nil(t, err)
		defer os.RemoveAll(tmpDir)
		err = Extract(root, tmpDir, func(n *SimpleNode) bool { return false })
		require.Nil(t, err)
		// Assert fileName is in tmpdir
		assert.FileExists(t, filepath.Join(tmpDir, expectedPath))
	})
}
