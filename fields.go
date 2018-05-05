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

// Merge info for commit
type Merge struct {
	Ref    string
	Source string
}

// Revert info for commit
type Revert struct {
	Header string
}

// Ref is abstract data related to commit. (e.g. `Issues`, `Pull Request`)
type Ref struct {
	Action string // (e.g. `Closes`)
	Ref    string // (e.g. `123`)
	Source string // (e.g. `owner/repository`)
}

// Note of commit
type Note struct {
	Title string // (e.g. `BREAKING CHANGE`)
	Body  string // `Note` content body
}

// NoteGroup is a collection of `Note` grouped by titles
type NoteGroup struct {
	Title string
	Notes []*Note
}

// Commit data
type Commit struct {
	Hash      *Hash
	Author    *Author
	Committer *Committer
	Merge     *Merge  // If it is not a merge commit, `nil` is assigned
	Revert    *Revert // If it is not a revert commit, `nil` is assigned
	Refs      []*Ref
	Notes     []*Note
	Mentions  []string // Name of the user included in the commit header or body
	Header    string   // (e.g. `feat(core): Add new feature`)
	Type      string   // (e.g. `feat`)
	Scope     string   // (e.g. `core`)
	Subject   string   // (e.g. `Add new feature`)
	Body      string
}

// CommitGroup is a collection of commits grouped according to the `CommitGroupBy` option
type CommitGroup struct {
	RawTitle string // Raw title before conversion (e.g. `build`)
	Title    string // Conversion by `CommitGroupTitleMaps` option, or title converted in title case (e.g. `Build`)
	Commits  []*Commit
}

// RelateTag is sibling tag data of `Tag`.
// If you give `Tag`, the reference hierarchy will be deepened.
// This struct is used to minimize the hierarchy of references
type RelateTag struct {
	Name    string
	Subject string
	Date    time.Time
}

// Tag is data of git-tag
type Tag struct {
	Name     string
	Subject  string
	Date     time.Time
	Next     *RelateTag
	Previous *RelateTag
}

// Version is a tag-separeted datset to be included in CHANGELOG
type Version struct {
	Tag           *Tag
	CommitGroups  []*CommitGroup
	Commits       []*Commit
	MergeCommits  []*Commit
	RevertCommits []*Commit
	NoteGroups    []*NoteGroup
}

// Unreleased is unreleased commit dataset
type Unreleased struct {
	CommitGroups  []*CommitGroup
	Commits       []*Commit
	MergeCommits  []*Commit
	RevertCommits []*Commit
	NoteGroups    []*NoteGroup
}
