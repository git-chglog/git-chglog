{{range .Versions}}
<a name="{{urlquery .Tag.Name}}"></a>
## {{.Tag.Name}} ({{datetime "2006-01-02" .Tag.Date}})
{{range .CommitGroups}}
### {{.Title}}
{{range .Commits}}
* {{if ne .Scope ""}}**{{.Scope}}:** {{end}}{{.Subject}}{{end}}
{{end}}{{if .RevertCommits}}
### Reverts
{{range .RevertCommits}}
* {{.Revert.Header}}{{end}}
{{end}}{{if .MergeCommits}}
### Pull Requests
{{range .MergeCommits}}
* {{.Header}}{{end}}
{{end}}{{range .NoteGroups}}
### {{.Title}}
{{range .Notes}}
{{.Body}}
{{end}}
{{end}}
{{end}}
