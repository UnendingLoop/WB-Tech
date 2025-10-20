// Package saver provides physical write of downloaded content to userdrive
package saver

import (
	"os"
	"path/filepath"
)

type Saver struct {
	BaseDir string
}

func (s *Saver) Save(path string, data []byte) (int, error) {
	// Создать все директории
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return 0, err
	}

	// Записать файл
	err := os.WriteFile(path, data, 0o644)
	if err != nil {
		return 0, err
	}

	return len(data), nil
}
