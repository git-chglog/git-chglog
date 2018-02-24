package main

// Builder ...
type Builder interface {
	Build(*Answer) (string, error)
}
