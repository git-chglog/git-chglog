package main

import (
	"io"

	chglog "github.com/git-chglog/git-chglog"
)

// Generator ...
type Generator interface {
	Generate(*chglog.Logger, io.Writer, string, *chglog.Config) error
}

type generatorImpl struct{}

// NewGenerator ...
func NewGenerator() Generator {
	return &generatorImpl{}
}

// Generate ...
func (*generatorImpl) Generate(logger *chglog.Logger, w io.Writer, query string, config *chglog.Config) error {
	return chglog.NewGenerator(logger, config).Generate(w, query)
}
