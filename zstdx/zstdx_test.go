package zstdx_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
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

func TestUncompressWithCustomSizeLimit(t *testing.T) {
	// Build a tarball containing a single file of a known size so the size
	// limit can be exercised precisely.
	const contentSize = 2048
	srcDir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(srcDir, "data.txt"),
		bytes.Repeat([]byte("a"), contentSize),
		0o600,
	))
	tarball := filepath.Join(t.TempDir(), "archive.tar.zst")
	require.NoError(t, zstdx.Compress(srcDir, tarball))

	t.Run("succeeds", func(t *testing.T) {
		tests := map[string]struct {
			maxFileSize int64
		}{
			"limit equal to file size":    {maxFileSize: contentSize},
			"limit larger than file size": {maxFileSize: contentSize * 2},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				outDir := t.TempDir()
				files, err := zstdx.UncompressWithCustomSizeLimit(tarball, outDir, tt.maxFileSize)
				require.NoError(t, err)
				require.NotEmpty(t, files)
				// The content must be fully extracted, not silently truncated.
				for _, f := range files {
					info, err := os.Stat(f)
					require.NoError(t, err)
					assert.Equal(t, int64(contentSize), info.Size())
				}
			})
		}
	})

	t.Run("fails", func(t *testing.T) {
		tests := map[string]struct {
			maxFileSize int64
		}{
			"limit one byte below file size": {maxFileSize: contentSize - 1},
			"limit far below file size":      {maxFileSize: 10},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				outDir := t.TempDir()
				_, err := zstdx.UncompressWithCustomSizeLimit(tarball, outDir, tt.maxFileSize)
				require.Error(t, err)
			})
		}
	})
}

func TestUncompressWithLimits(t *testing.T) {
	// buildTarball writes fileCount files of fileSize bytes each into a fresh
	// source directory and compresses it into a tarball, returning the tarball
	// path. The resulting archive has fileCount file entries plus one directory
	// entry for the root.
	buildTarball := func(t *testing.T, fileCount int, fileSize int) string {
		t.Helper()
		srcDir := t.TempDir()
		for i := 0; i < fileCount; i++ {
			require.NoError(t, os.WriteFile(
				filepath.Join(srcDir, fmt.Sprintf("data%d.txt", i)),
				bytes.Repeat([]byte("a"), fileSize),
				0o600,
			))
		}
		tarball := filepath.Join(t.TempDir(), "archive.tar.zst")
		require.NoError(t, zstdx.Compress(srcDir, tarball))
		return tarball
	}

	t.Run("succeeds", func(t *testing.T) {
		const fileCount = 4
		const fileSize = 1024
		tarball := buildTarball(t, fileCount, fileSize)

		tests := map[string]zstdx.Limits{
			"generous limits": {
				MaxFileSize:  fileSize,
				MaxTotalSize: fileCount * fileSize,
				MaxEntries:   fileCount + 1, // files + root directory
			},
			"zero-value limits fall back to defaults": {},
		}
		for name, limits := range tests {
			t.Run(name, func(t *testing.T) {
				outDir := t.TempDir()
				files, err := zstdx.UncompressWithLimits(tarball, outDir, limits)
				require.NoError(t, err)
				assert.Len(t, files, fileCount)
				for _, f := range files {
					info, err := os.Stat(f)
					require.NoError(t, err)
					assert.Equal(t, int64(fileSize), info.Size())
				}
			})
		}
	})

	t.Run("fails", func(t *testing.T) {
		const fileCount = 4
		const fileSize = 1024
		tarball := buildTarball(t, fileCount, fileSize)

		tests := map[string]zstdx.Limits{
			"total size below cumulative size": {
				MaxFileSize:  fileSize,               // each file is within the per-file cap
				MaxTotalSize: fileCount*fileSize - 1, // but the sum exceeds the total cap
				MaxEntries:   fileCount + 1,
			},
			"entry count below number of entries": {
				MaxFileSize:  fileSize,
				MaxTotalSize: fileCount * fileSize,
				MaxEntries:   fileCount, // one short: files + root directory exceed this
			},
		}
		for name, limits := range tests {
			t.Run(name, func(t *testing.T) {
				outDir := t.TempDir()
				_, err := zstdx.UncompressWithLimits(tarball, outDir, limits)
				require.Error(t, err)
			})
		}
	})
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
