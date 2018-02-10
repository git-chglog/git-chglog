# go-gitcmd

[![godoc.org](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/tsuyoshiwada/go-gitcmd)
[![Travis](https://img.shields.io/travis/tsuyoshiwada/go-gitcmd.svg?style=flat-square)](https://travis-ci.org/tsuyoshiwada/go-gitcmd)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/tsuyoshiwada/go-gitcmd/blob/master/LICENSE)

> Go (golang) package providing for tiny git command wrapper.




## Installation

To install `go-gitcmd`, simply run:

```bash
$ go get -u github.com/tsuyoshiwada/go-gitcmd
```




## Usage

It is the simplest example.

```go
package main

import (
	"log"
	"github.com/tsuyoshiwada/go-gitcmd"
)

func main() {
	git := gitcmd.New(nil) // or `git := gitcmd.New(&Config{Bin: "/your/custom/git/bin"})`

	out, err := git.Exec("rev-parse", "--git-dir")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out) // ".git"
}
```

See [godoc](https://godoc.org/github.com/tsuyoshiwada/go-gitcmd) for API detail :+1:




## How to Mock

Since `Client` is an interface, it can easily be Mocked.

```go
type MockClient struct {
	gitcmd.Client // Interface embedding
	ReturnCanExec        func() error
	ReturnExec           func(string, ...string) (string, error)
	ReturnInsideWorkTree func() error
}

func (m *MockClient) CanExec() error {
	return m.ReturnCanExec()
}

func (m *MockClient) Exec(subcmd string, args ...string) (string, error) {
	return m.ReturnExec(subcmd, args...)
}

func (m *MockClient) InsideWorkTree() error {
	return m.ReturnInsideWorkTree()
}

func main() {
	git := &MockClient{}

	// Set `InsideWorkTree()` mock function
	git.ReturnInsideWorkTree = func() error {
		return errors.New("error...")
	}

	err := git.InsideWorkTree()
	fmt.Println(err.Error()) // "error..."
}
```




## Contribute

1. Fork (https://github.com/tsuyoshiwada/go-gitcmd)
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test` command and confirm that it passes
1. Create new Pull Request :muscle:

Bugs, feature requests and comments are more than welcome in the [issues](https://github.com/tsuyoshiwada/go-gitcmd/issues).




## License

[MIT © tsuyoshiwada](./LICENSE)
