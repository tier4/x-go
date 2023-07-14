package zstdx

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
	"github.com/klauspost/compress/zstd"
	"github.com/pkg/errors"
)

// For protection from decompression bomb
const maxFileSize int64 = 16 * 1024 * 1024 * 1024

// The code at https://go.dev/play/p/A2GXsDFWx9m is used as a refference
func Uncompress(tarball, targetDir string) ([]string, error) {
	file, err := os.Open(filepath.Clean(tarball))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "cannot read file information")
	}
	header := make([]byte, min(262, stat.Size()))
	if _, err := io.ReadFull(file, header); err != nil {
		return nil, errors.Wrap(err, "cannot determine type by reading file header")
	}
	if _, err := file.Seek(0, 0); err != nil {
		return nil, errors.Wrap(err, "cannot seek file")
	}

	if !filetype.Is(header, "zst") {
		return nil, fmt.Errorf("Unknown file format when trying to uncompress %s", tarball)
	}

	reader, err := zstd.NewReader(file)
	if err != nil {
		return nil, err
	}
	return untar(reader, targetDir)
}

func min(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func untar(reader io.Reader, targetDir string) ([]string, error) {
	var extractedFiles []string
	tarReader := tar.NewReader(reader)

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
			if err := untarFile(tarReader, header, path); err != nil {
				return nil, err
			}
			extractedFiles = append(extractedFiles, path)
		}
	}
}

func untarFile(tarReader *tar.Reader, header *tar.Header, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, header.FileInfo().Mode()|0444)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.CopyN(file, tarReader, maxFileSize); err != nil && err != io.EOF {
		return err
	}

	return nil
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
			data, err := os.Open(file)
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
