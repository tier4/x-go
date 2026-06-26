package zstdx

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
)

// zstdMagic is the magic number for Zstandard compressed data (RFC 8478).
var zstdMagic = []byte{0x28, 0xB5, 0x2F, 0xFD}

// For protection from decompression bomb
const (
	defaultMaxFileSize  int64 = 16 * 1024 * 1024 * 1024 // per-file uncompressed size cap
	defaultMaxTotalSize int64 = 64 * 1024 * 1024 * 1024 // cumulative uncompressed size cap
	defaultMaxEntries   int   = 1_000_000               // max number of tar entries
)

// Limits bounds resource usage during extraction to protect against
// decompression bombs. A non-positive field falls back to its default, so a
// zero-value Limits is equivalent to DefaultLimits().
type Limits struct {
	// MaxFileSize is the maximum uncompressed size of a single entry, in bytes.
	MaxFileSize int64
	// MaxTotalSize is the maximum cumulative uncompressed size across all
	// entries, in bytes.
	MaxTotalSize int64
	// MaxEntries is the maximum number of tar entries (directories and files)
	// that may be processed.
	MaxEntries int
}

// DefaultLimits returns the default extraction limits.
func DefaultLimits() Limits {
	return Limits{
		MaxFileSize:  defaultMaxFileSize,
		MaxTotalSize: defaultMaxTotalSize,
		MaxEntries:   defaultMaxEntries,
	}
}

// withDefaults replaces any non-positive field with its default value.
func (l Limits) withDefaults() Limits {
	if l.MaxFileSize <= 0 {
		l.MaxFileSize = defaultMaxFileSize
	}
	if l.MaxTotalSize <= 0 {
		l.MaxTotalSize = defaultMaxTotalSize
	}
	if l.MaxEntries <= 0 {
		l.MaxEntries = defaultMaxEntries
	}
	return l
}

// Uncompress with the default extraction limits (see DefaultLimits).
func Uncompress(tarball, targetDir string) ([]string, error) {
	return uncompress(tarball, targetDir, DefaultLimits())
}

// UncompressWithCustomSizeLimit with a specified per-file max size limit. The
// total-size and entry-count limits use their defaults.
func UncompressWithCustomSizeLimit(tarball, targetDir string, maxFileSize int64) ([]string, error) {
	return uncompress(tarball, targetDir, Limits{MaxFileSize: maxFileSize})
}

// UncompressWithLimits with fully customized extraction limits. Any
// non-positive field in limits falls back to its default.
func UncompressWithLimits(tarball, targetDir string, limits Limits) ([]string, error) {
	return uncompress(tarball, targetDir, limits)
}

// The code at https://go.dev/play/p/A2GXsDFWx9m is used as a reference
func uncompress(tarball, targetDir string, limits Limits) ([]string, error) {
	limits = limits.withDefaults()

	file, err := os.Open(filepath.Clean(tarball))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("cannot read file information: %w", err)
	}
	header := make([]byte, min(4, stat.Size()))
	if _, err := io.ReadFull(file, header); err != nil {
		return nil, fmt.Errorf("cannot determine type by reading file header: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("cannot seek file: %w", err)
	}

	if !bytes.Equal(header, zstdMagic) {
		return nil, fmt.Errorf("unknown file format when trying to uncompress %s", tarball)
	}

	reader, err := zstd.NewReader(file)
	if err != nil {
		return nil, err
	}
	return untar(reader, targetDir, limits)
}

func untar(reader io.Reader, targetDir string, limits Limits) ([]string, error) {
	var extractedFiles []string
	tarReader := tar.NewReader(reader)

	var totalWritten int64
	var entries int

	for {
		header, err := tarReader.Next()
		switch {
		// if no more files are found return
		case err == io.EOF:
			return extractedFiles, nil

		// return any other error
		case err != nil:
			return extractedFiles, err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// Bound the number of entries to protect against archives with a huge
		// number of (possibly tiny) files or directories.
		entries++
		if entries > limits.MaxEntries {
			return nil, fmt.Errorf("archive exceeds the maximum allowed number of entries (%d)", limits.MaxEntries)
		}

		// the target location where the dir/file should be created
		path, err := SanitizeExtractPath(header.Name, targetDir)
		if err != nil {
			return nil, err
		}

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if err := os.MkdirAll(path, header.FileInfo().Mode()|0444); err != nil {
				return nil, err
			}

		// if it's a file create it
		case tar.TypeReg, tar.TypeGNUSparse:
			// tar.Next() will externally only iterate files, so we might have to create intermediate directories here
			written, err := untarFile(tarReader, header, path, limits.MaxFileSize, limits.MaxTotalSize-totalWritten)
			if err != nil {
				return nil, err
			}
			totalWritten += written
			// Bound the cumulative uncompressed size across all entries.
			if totalWritten > limits.MaxTotalSize {
				return nil, fmt.Errorf("archive exceeds the maximum allowed total uncompressed size of %d bytes", limits.MaxTotalSize)
			}
			extractedFiles = append(extractedFiles, path)
		}
	}
}

func untarFile(tarReader *tar.Reader, header *tar.Header, path string, maxFileSize, remainingTotal int64) (int64, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return 0, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, header.FileInfo().Mode()|0444)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Copy at most one byte beyond the smaller of the per-file limit and the
	// remaining total budget, so an entry that exceeds either limit is detected
	// and rejected instead of being silently truncated or exhausting the disk
	// (decompression-bomb protection). Guard the +1 against int64 overflow: when
	// the limit is math.MaxInt64 (used to effectively disable a cap), adding one
	// would wrap to a negative count, making io.CopyN silently copy nothing.
	copyLimit := min(maxFileSize, remainingTotal)
	if copyLimit < math.MaxInt64 {
		copyLimit++
	}
	written, err := io.CopyN(file, tarReader, copyLimit)
	if err != nil && err != io.EOF {
		return written, err
	}
	if written > maxFileSize {
		return written, fmt.Errorf("file %q exceeds the maximum allowed size of %d bytes", header.Name, maxFileSize)
	}

	return written, nil
}

// cf. https://snyk.io/research/zip-slip-vulnerability
func SanitizeExtractPath(filePath, destination string) (string, error) {
	path := filepath.Join(destination, filePath)
	if !strings.HasPrefix(path, filepath.Clean(destination)+string(os.PathSeparator)) {
		return "", fmt.Errorf("%s: illegal file path", filePath)
	}
	return path, nil
}

func Compress(src, dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	enc, err := zstd.NewWriter(out)
	if err != nil {
		return err
	}
	defer enc.Close()

	tarWriter := tar.NewWriter(enc)
	defer tarWriter.Close()

	root, err := os.OpenRoot(src)
	if err != nil {
		return err
	}
	defer root.Close()

	dir := filepath.Base(src)
	// walk through every file in the folder
	return filepath.Walk(src, func(file string, fi os.FileInfo, e error) error {
		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}
		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = filepath.ToSlash(file)
		header.Name = strings.Replace(header.Name, src, dir, 1)

		if filepath.IsAbs(header.Name) {
			if len(header.Name) <= 1 {
				// similar to what the 'tar' command does
				header.Name = "./"
			}
			// the 'tar' command strips the leading '/' from absolute paths
			header.Name = header.Name[1:]
		}

		// write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			rel, err := filepath.Rel(src, file)
			if err != nil {
				return err
			}
			data, err := root.Open(rel)
			if err != nil {
				return err
			}
			defer data.Close()
			if _, err := io.Copy(tarWriter, data); err != nil {
				return err
			}
		}
		return nil
	})
}
