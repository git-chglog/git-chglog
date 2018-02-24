package main

import (
	"io"
	"io/ioutil"
	"os"
)

// FileSystem ...
type FileSystem interface {
	Exists(path string) bool
	MkdirP(path string) error
	Create(name string) (File, error)
	WriteFile(path string, content []byte) error
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

func (*osFileSystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (*osFileSystem) MkdirP(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func (*osFileSystem) Create(name string) (File, error) {
	return os.Create(name)
}

func (*osFileSystem) WriteFile(path string, content []byte) error {
	return ioutil.WriteFile(path, content, os.ModePerm)
}
