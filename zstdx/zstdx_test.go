package zstdx_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tier4/x-go/zstdx"
)

func TestUncompress(t *testing.T) {
	// Create a temporary directory to store the output
	tmpDir, err := os.MkdirTemp("", "uncompress-")
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)

	// Uncompress the tarball
	files, err := zstdx.Uncompress("testdata/sample.tar.zst", tmpDir)
	require.NoError(t, err)

	// Check if the files are correctly uncompressed
	assert.Equal(t, 2, len(files), "Expected two files to be uncompressed")
	for _, f := range files {
		info, err := os.Stat(f)
		require.NoError(t, err)
		assert.False(t, info.IsDir(), "Expected a file, not a directory")
	}
}

func TestCompress(t *testing.T) {
	// Create a temporary file to store the output
	tmpFile, err := os.CreateTemp("", "compress.tar.zst")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Compress the directory
	err = zstdx.Compress("testdata/sample", tmpFile.Name())
	require.NoError(t, err)

	// Check if the file is correctly compressed
	info, err := os.Stat(tmpFile.Name())
	require.NoError(t, err)
	assert.True(t, info.Size() > 0, "Expected the compressed file to be larger than 0 bytes")
}

func TestSanitizeExtractPath(t *testing.T) {
	t.Run("valid path", func(t *testing.T) {
		path, err := zstdx.SanitizeExtractPath("valid/path.txt", "/target/directory")
		require.NoError(t, err)
		assert.Equal(t, "/target/directory/valid/path.txt", path)
	})

	t.Run("invalid path", func(t *testing.T) {
		_, err := zstdx.SanitizeExtractPath("../invalid/path.txt", "/target/directory")
		assert.Error(t, err)
	})
}
