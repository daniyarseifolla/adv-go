package pipeline

import (
	"os"
	"path/filepath"
	"strings"
)

func WalkDir(root, ext string) ([]string, error) {
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
		if strings.HasSuffix(d.Name(), ext) {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
