package pipeline

import (
	"os"
	"path/filepath"
	"strings"
)

func WalkDir(root, ext string, minSize, maxSize int64) ([]string, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return []string{root}, nil
	}

	var files []string

	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ext) {
			return nil
		}

		if minSize > 0 || maxSize > 0 {
			info, err := d.Info()
			if err != nil {
				return nil
			}
			if minSize > 0 && info.Size() < minSize {
				return nil
			}
			if maxSize > 0 && info.Size() > maxSize {
				return nil
			}
		}

		files = append(files, path)
		return nil
	})

	return files, err
}
