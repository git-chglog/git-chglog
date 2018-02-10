package chglog

import "time"

// Hash of commit
type Hash struct {
	Long  string
	Short string
}

// Author of commit
type Author struct {
	Name  string
	Email string
	Date  time.Time
}

// Committer of commit
type Committer struct {
	Name  string
	Email string
	Date  time.Time
}

// Merge ...
type Merge struct {
	Ref    string
	Source string
}

// Revert ...
type Revert struct {
	Raw     string
	Subject string
	Hash    string
}

// Ref ...
type Ref struct {
	Action string
	Ref    string
	Source string
}

// Note ...
type Note struct {
	Title string
	Body  string
}

// NoteGroup ...
type NoteGroup struct {
	Title string
	Notes []*Note
}

// Commit data
type Commit struct {
	Hash      *Hash
	Author    *Author
	Committer *Committer
	Merge     *Merge
	Revert    *Revert
	Refs      []*Ref
	Notes     []*Note
	Mentions  []string
	Header    string
	Type      string
	Scope     string
	Subject   string
	Body      string
}

// CommitGroup is group of commit
type CommitGroup struct {
	RawTitle string
	Title    string
	Commits  []*Commit
}

// Tag ...
type Tag struct {
	Name string
	Date time.Time
}

// Version ...
type Version struct {
	Tag           *Tag
	CommitGroups  []*CommitGroup
	MergeCommits  []*Commit
	RevertCommits []*Commit
	NoteGroups    []*NoteGroup
}
