# git-chglog

![git-chglog](https://raw.githubusercontent.com/git-chglog/artwork/master/repo-banner%402x.png)

[![godoc.org](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/git-chglog/git-chglog)
[![Travis](https://img.shields.io/travis/git-chglog/git-chglog.svg?style=flat-square)](https://travis-ci.org/git-chglog/git-chglog)
[![AppVeyor](https://img.shields.io/appveyor/ci/tsuyoshiwada/git-chglog/master.svg?style=flat-square)](https://ci.appveyor.com/project/tsuyoshiwada/git-chglog/branch/master)
[![Coverage Status](https://img.shields.io/coveralls/github/git-chglog/git-chglog.svg?style=flat-square)](https://coveralls.io/github/git-chglog/git-chglog?branch=master)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/git-chglog/git-chglog/blob/master/LICENSE)

> CHANGELOG generator implemented in Go (Golang).  
> _Anytime, anywhere, Write your CHANGELOG._




## Table of Contents

* [Features](#features)
* [How it works](#how-it-works)
* [Getting Started](#getting-started)
    - [Installation](#installation)
        + [Homebrew (for macOS users)](#homebrew-for-macos-users)
        + [Go users](#go-users)
    - [Test Installation](#test-installation)
    - [Quick Start](#quick-start)
* [CLI Usage](#cli-usage)
    - [`tag query`](#tag-query)
* [Configuration](#configuration)
* [Templates](#templates)
* [Supported Styles](#supported-styles)
* [FAQ](#faq)
* [TODO](#todo)
* [Thanks](#thanks)
* [Contributing](#contributing)
    - [Development](#development)
    - [Feedback](#feedback)
* [CHANGELOG](#changelog)
* [Related Projects](#related-projects)
* [License](#license)




## Features

* :recycle: High portability
    - It works with single binary. Therefore, any project (environment) can be used.
* :beginner: Simple usability
    - The CLI usage is very simple and has low learning costs.
    - For example, the simplest command is `$ git-chglog`.
* :rocket: High flexibility
    - Commit message format and ...
    - CHANGELOG's style (Template) and ...
    - etc ...




## How it works

`git-chglog` internally uses the `git` command to get data to include in CHANGELOG.  
The basic steps are as follows.

1. Get all the tags.
1. Get the commit contained between `tagA` and `tagB`.
1. Execute with all tags corresponding to [tag query](#tag-query) that specified Step 1 to 2.




## Getting Started

We will start with installation and introduce the steps up to the automatic generation of the configuration file and template.


### Installation

Please install `git-chglog` in a way that matches your environment.

#### Homebrew (for macOS users)

```bash
$ brew tap git-chglog/git-chglog
$ brew install git-chglog
```

#### [Scoop](https://scoop.sh) (for Windows users)

```
$ scoop install git-chglog
```

#### Go users

```bash
$ go get -u github.com/git-chglog/git-chglog/cmd/git-chglog
```

---

If you are in another platform, you can download binary from [release page](https://github.com/git-chglog/git-chglog/releases) and place it in `$PATH` directory.


### Test Installation

You can check with the following command whether the `git-chglog` command was included in a valid `$PATH`.

```bash
$ git-chglog --version
# output the git-chglog version
```


### Quick Start

`git-chglog` requires configuration files and templates to generate CHANGELOG.  
However, it is a waste of time to create configuration files and templates with scratch.

Therefore we recommend using the `--init` option which can create them interactively :+1:

```bash
$ git-chglog --init
```

![init option demo](./docs/assets/init.gif)

---

You are now ready for configuration files and templates!

Let's immediately generate CHANGELOG of your project.  
By doing the following simple command, Markdown of CHANGELOG is displayed on stdout.

```bash
$ git-chglog
```

Use `-o` (`--output`) option if you want to output to file instead of stdout.

```bash
$ git-chglog -o CHANGELOG.md
```

---

This is how basic usage is over!  
In order to make better CHANGELOG, please refer to the following document and customize it.




## CLI Usage

```bash
$ git-chglog --help

USAGE:
  git-chglog [options] <tag query>

    There are the following specification methods for <tag query>.

    1. <old>..<new> - Commit contained in <old> tags from <new>.
    2. <name>..     - Commit from the <name> to the latest tag.
    3. ..<name>     - Commit from the oldest tag to <name>.
    4. <name>       - Commit contained in <name>.

OPTIONS:
  --init                    generate the git-chglog configuration file in interactive
  --config value, -c value  specifies a different configuration file to pick up (default: ".chglog/config.yml")
  --output value, -o value  output path and filename for the changelogs. If not specified, output to stdout
  --next-tag value          treat unreleased commits as specified tags (EXPERIMENTAL)
  --silent                  disable stdout output
  --no-color                disable color output [$NO_COLOR]
  --no-emoji                disable emoji output [$NO_EMOJI]
  --help, -h                show help
  --version, -v             print the version

EXAMPLE:

  $ git-chglog

    If <tag query> is not specified, it corresponds to all tags.
    This is the simplest example.

  $ git-chglog 1.0.0..2.0.0

    The above is a command to generate CHANGELOG including commit of 1.0.0 to 2.0.0.

  $ git-chglog 1.0.0

    The above is a command to generate CHANGELOG including commit of only 1.0.0.

  $ git-chglog $(git describe --tags $(git rev-list --tags --max-count=1))

    The above is a command to generate CHANGELOG with the commit included in the latest tag.

  $ git-chglog --output CHANGELOG.md

    The above is a command to output to CHANGELOG.md instead of standard output.

  $ git-chglog --config custom/dir/config.yml

    The above is a command that uses a configuration file placed other than ".chglog/config.yml".
```


### `tag query`

You can specify a commit to include in the generation of CHANGELOG using `<tag query>`.  
The table below shows Query patterns and summaries, and Query examples.

| Query          | Description                                    | Example                     |
|:---------------|:-----------------------------------------------|:----------------------------|
| `<old>..<new>` | Commit contained in `<new>` tags from `<old>`. | `$ git-chglog 1.0.0..2.0.0` |
| `<name>..`     | Commit from the `<name>` to the latest tag.    | `$ git-chglog 1.0.0..`      |
| `..<name>`     | Commit from the oldest tag to `<name>`.        | `$ git-chglog ..2.0.0`      |
| `<name>`       | Commit contained in `<name>`.                  | `$ git-chglog 1.0.0`        |




## Configuration

The `git-chglog` configuration is write with the yaml file. Default location is `.chglog/config.yml`.

Below is a complete list that you can use with `git-chglog`.

```yaml
bin: git
style: ""
template: CHANGELOG.tpl.md
info:
  title: CHANGELOG
  repository_url: https://github.com/git-chglog/git-chglog

options:
  commits:
    filters:
      Type:
        - feat
    sort_by: Scope

  commit_groups:
    group_by: Type
    sort_by: Title
    title_maps:
      feat: Features

  header:
    pattern: "<regexp>"
    pattern_maps:
      - PropName

  issues:
    prefix:
      - #

  refs:
    actions:
      - Closes
      - Fixes

  merges:
    pattern: "^Merge branch '(\\w+)'$"
    pattern_maps:
      - Source

  reverts:
    pattern: "^Revert \"([\\s\\S]*)\"$"
    pattern_maps:
      - Header

  notes:
    keywords:
      - BREAKING CHANGE
```


### `bin`

Git execution command.

| Required | Type   | Default | Description |
|:---------|:-------|:--------|:------------|
| N        | String | `"git"` | -           |


### `style`

CHANGELOG style. Automatic linking of issues and notices, initial value setting such as merges etc. are done automatically.

| Required | Type   | Default  | Description                                            |
|:---------|:-------|:---------|:-------------------------------------------------------|
| N        | String | `"none"` | Should be `"github"` `"gitlab"` `"bitbucket"` `"none"` |


### `template`

Path for template file. It is specified by a relative path from the setting file. Absolute pass is ok.

| Required | Type   | Default              | Description |
|:---------|:-------|:---------------------|:------------|
| N        | String | `"CHANGELOG.tpl.md"` | -           |


### `info`

Metadata for CHANGELOG. Depending on Style, it is sometimes used in processing, so it is recommended to specify it.

| Key              | Required | Type   | Default       | Description            |
|:-----------------|:---------|:-------|:--------------|:-----------------------|
| `title`          | N        | String | `"CHANGELOG"` | Title of CHANGELOG.    |
| `repository_url` | N        | String | none          | URL of git repository. |


### `options`

Options used to process commits.

#### `options.commits`

Option concerning acquisition and sort of commit.

| Key       | Required | Type        | Default   | Description                                                                                                         |
|:----------|:---------|:------------|:----------|:--------------------------------------------------------------------------------------------------------------------|
| `filters` | N        | Map in List | none      | Filter by using `Commit` properties and values. Filtering is not done by specifying an empty value.                 |
| `sort_by` | N        | String      | `"Scope"` | Property name to use for sorting `Commit`. See [Commit](https://godoc.org/github.com/git-chglog/git-chglog#Commit). |

#### `options.commit_groups`

Option for groups of commits.

| Key          | Required | Type        | Default   | Description                                                                                |
|:-------------|:---------|:------------|:----------|:-------------------------------------------------------------------------------------------|
| `group_by`   | N        | String      | `"Type"`  | Property name of `Commit` to be grouped into `CommitGroup`. See [CommitGroup][doc-commit]. |
| `sort_by`    | N        | String      | `"Title"` | Property name to use for sorting `CommitGroup`. See [CommitGroup][doc-commit-group].       |
| `title_maps` | N        | Map in List | none      | Map for `CommitGroup` title conversion.                                                    |

#### `options.header`

This option is used for parsing the commit header.

| Key            | Required | Type   | Default | Description                                                                                             |
|:---------------|:---------|:-------|:--------|:--------------------------------------------------------------------------------------------------------|
| `pattern`      | Y        | String | none    | A regular expression to use for parsing the commit header.                                              |
| `pattern_maps` | Y        | List   | none    | A rule for mapping the result of `HeaderPattern` to the property of `Commit`. See [Commit][doc-commit]. |

#### `options.issues`

This option is used to detect issues.

| Key      | Required | Type | Default | Description                                |
|:---------|:---------|:-----|:--------|:-------------------------------------------|
| `prefix` | N        | List | none    | Prefix used for issues. (e.g. `#`, `#gh-`) |

#### `options.refs`

This option is for parsing references.

| Key       | Required | Type | Default | Description                                    |
|:----------|:---------|:-----|:--------|:-----------------------------------------------|
| `actions` | N        | List | none    | Word list of `Ref.Action`. See [Ref][doc-ref]. |

#### `options.merges`

Option to detect and parse merge commit.

| Key            | Required | Type   | Default | Description                               |
|:---------------|:---------|:-------|:--------|:------------------------------------------|
| `pattern`      | N        | String | none    | Similar to `options.header.pattern`.      |
| `pattern_maps` | N        | List   | none    | Similar to `options.header.pattern_maps`. |

#### `options.reverts`

Option to detect and parse revert commit.

| Key            | Required | Type   | Default | Description                               |
|:---------------|:---------|:-------|:--------|:------------------------------------------|
| `pattern`      | N        | String | none    | Similar to `options.header.pattern`.      |
| `pattern_maps` | N        | List   | none    | Similar to `options.header.pattern_maps`. |

#### `options.notes`

Option to detect notes contained in commit body.

| Key        | Required | Type | Default | Description                                                                                          |
|:-----------|:---------|:-----|:--------|:-----------------------------------------------------------------------------------------------------|
| `keywords` | N        | List | none    | Keyword list to find `Note`. A semicolon is a separator, like `<keyword>:` (e.g. `BREAKING CHANGE`). |




## Templates

The `git-chglog` template uses the `text/template` package. For basic usage please refer to the following.

> [text/template](https://golang.org/pkg/text/template/)

If you are not satisfied with the prepared template please try customizing.

---

The basic templates are as follows.

**Example:**

```markdown
{{ if .Versions -}}
<a name="unreleased"></a>
## [Unreleased]

{{ if .Unreleased.CommitGroups -}}
{{ range .Unreleased.CommitGroups -}}
### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}

{{ range .Versions }}
<a name="{{ .Tag.Name }}"></a>
## {{ if .Tag.Previous }}[{{ .Tag.Name }}]{{ else }}{{ .Tag.Name }}{{ end }} - {{ datetime "2006-01-02" .Tag.Date }}
{{ range .CommitGroups -}}
### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end }}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts
{{ range .RevertCommits -}}
- {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .MergeCommits -}}
### Pull Requests
{{ range .MergeCommits -}}
- {{ .Header }}
{{ end }}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}
{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}

{{- if .Versions }}
[Unreleased]: {{ .Info.RepositoryURL }}/compare/{{ $latest := index .Versions 0 }}{{ $latest.Tag.Name }}...HEAD
{{ range .Versions -}}
{{ if .Tag.Previous -}}
[{{ .Tag.Name }}]: {{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}
{{ end -}}
{{ end -}}
{{ end -}}
```

See godoc [RenderData][doc-render-data] for available variables.



## Supported Styles

| Name                                       | Status             | Features                                               |
|:-------------------------------------------|:-------------------|:-------------------------------------------------------|
| [GitHub](https://github.com/)              | :white_check_mark: | Mentions automatic link. Automatic link to references. |
| [GitLab](https://about.gitlab.com/)        | :white_check_mark: | Mentions automatic link. Automatic link to references. |
| [Bitbucket](https://bitbucket.org/product) | :white_check_mark: | Mentions automatic link. Automatic link to references. |

> :memo: Even with styles that are not yet supported, it is possible to make ordinary CHANGELOG.




## FAQ

<details>
  <summary>Why do not you output files by default?</summary>
  This is not for the purpose of completely automating the generation of CHANGELOG, because it is only for the purpose of assisting generation.

  It is ideal to describe everything included in CHANGELOG in commit. But actually it is very difficult to do it perfectly.

  There are times when you need your hands to write a great CHANGELOG.  
  By displaying it on the standard output, it makes it easy to change the contents.
</details>

<details>
  <summary>Can I commit CHANGELOG changes before creating tags?</summary>

  Yes, it can be solved by using the `--next-tag` flag.

  For example, let's say you want to upgrade your project to `2.0.0`.  
  You can create CHANGELOG containing `2.0.0` as follows.

  ```bash
  $ git-chglog --next-tag 2.0.0 -o CHANGELOG.md
  $ git commit -am "release 2.0.0"
  $ git tag 2.0.0
  ```

  The point to notice is that before actually creating a tag with git, it is conveying the next version with `--next-tag` :+1:

  This is a step that is necessary for project operation in many cases.
</details>




## TODO

* [x] Windows Support
* [x] More styles (GitHub, GitLab, Bitbucket :tada:)
* [ ] Snippetization of configuration files (improvement of reusability)
* [ ] More test test test ... (and example)




## Thanks

`git-chglog` is inspired by [conventional-changelog](https://github.com/conventional-changelog/conventional-changelog). Thank you!




## Contributing

We are always welcoming your contribution :clap:


### Development

1. Fork (https://github.com/git-chglog/git-chglog) :tada:
1. Create a feature branch :coffee:
1. Run test suite with the `$ make test` command and confirm that it passes :zap:
1. Commit your changes :memo:
1. Rebase your local changes against the `master` branch :bulb:
1. Create new Pull Request :love_letter:

Bugs, feature requests and comments are more than welcome in the [issues](https://github.com/git-chglog/git-chglog/issues).


### Feedback

I would like to make `git-chglog` a better tool.  
The goal is to be able to use in various projects.

Therefore, your feedback is very useful.  
I am very happy to tell you your opinions on Issues and PR :heart:




## CHANGELOG

See [CHANGELOG.md](./CHANGELOG.md)




## Related Projects

* [git-chglog/artwork](https://github.com/git-chglog/artwork) - Assets for `git-chglog`.




## License

[MIT © tsuyoshiwada](./LICENSE)




[doc-commit]: https://godoc.org/github.com/git-chglog/git-chglog#Commit
[doc-commit-group]: https://godoc.org/github.com/git-chglog/git-chglog#Commit
[doc-ref]: https://godoc.org/github.com/git-chglog/git-chglog#Ref
[doc-render-data]: https://godoc.org/github.com/git-chglog/git-chglog#RenderData
