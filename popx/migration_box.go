package popx

import (
	"bytes"
	"embed"
	"io/fs"
	"strings"

	"github.com/gobuffalo/packd"
)

type migrationsFS struct {
	dir embed.FS
}

func NewMigrationBox(fs embed.FS) packd.Walkable {
	return &migrationsFS{dir: fs}
}

func (m *migrationsFS) Walk(wf packd.WalkFunc) error {
	return fs.WalkDir(m.dir, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		content, err := m.dir.ReadFile(path)
		if err != nil {
			return err
		}
		f, err := packd.NewFile(info.Name(), bytes.NewReader(content))
		if err != nil {
			return err
		}

		return wf(path, f)
	})
}

func (m *migrationsFS) WalkPrefix(prefix string, wf packd.WalkFunc) error {
	return m.Walk(func(path string, file packd.File) error {
		if strings.HasPrefix(path, prefix) {
			return wf(path, file)
		}
		return nil
	})
}
