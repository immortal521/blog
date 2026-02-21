package test

import (
	"path/filepath"
	"runtime"

	"blog-server/config"
)

func LoadConfig() (*config.Config, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return config.Load(filepath.Join(dir, "../config.yml"))
}
