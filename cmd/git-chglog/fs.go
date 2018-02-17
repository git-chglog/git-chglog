package main

import (
	"io"
	"os"
)

// FileSystem ...
type FileSystem interface {
	MkdirP(path string) error
	Create(name string) (File, error)
}

// File ...
type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	Stat() (os.FileInfo, error)
}

var fs = &osFileSystem{}

type osFileSystem struct{}

func (*osFileSystem) MkdirP(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func (*osFileSystem) Create(name string) (File, error) {
	return os.Create(name)
}
