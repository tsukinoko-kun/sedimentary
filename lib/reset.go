package lib

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func (sdmt *Sedimentary) Reset() error {
	err := filepath.WalkDir(sdmt.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		relPath := sdmt.relativePath(path)
		content, err := sdmt.ReadTextFile(relPath)
		if err != nil {
			return errors.Join(fmt.Errorf("failed to read text file content"), err)
		}
		if err := os.WriteFile(path, []byte(content), fs.ModeExclusive); err != nil {
			return errors.Join(fmt.Errorf("failed to write %q", path), err)
		}
		return nil
	})
	if err != nil {
		return errors.Join(fmt.Errorf("failed to walk dir %q", sdmt.root), err)
	}
	return nil
}
