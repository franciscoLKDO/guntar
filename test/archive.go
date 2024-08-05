package test

import (
	"archive/tar"
	"bytes"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

type File struct {
	Name string
	Mode fs.FileMode
	Body string
}

// CreateArchive for tests, this funcion will return a tar archive buffer based on given files
func CreateArchive(t *testing.T, files []File) *bytes.Buffer {
	// Create and add some files to the archive.
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: int64(file.Mode),
			Size: int64(len(file.Body)),
		}
		require.Nil(t, tw.WriteHeader(hdr))
		_, err := tw.Write([]byte(file.Body))
		require.Nil(t, err)
	}
	require.Nil(t, tw.Close())
	return &buf
}
