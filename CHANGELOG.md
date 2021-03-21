# CHANGELOG
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html)

<a name="unreleased"></a>
## [Unreleased]

<a name="v0.12.0"></a>
## [v0.12.0] - 2021-03-20
### Chore
- bumps version to v0.12.0
- bump golang to 1.16 ([#118](https://github.com/git-chglog/git-chglog/issues/118))
- **ci:** add golangci-lint action and apply linting changes ([#120](https://github.com/git-chglog/git-chglog/issues/120))

### Features
- allow tag sorting by semver ([#124](https://github.com/git-chglog/git-chglog/issues/124))

### BREAKING CHANGE

`JiraIssueId` has been renamed to `JiraIssueID`. This impacts the value for `pattern_maps` in `config.yml`.

<a name="v0.11.2"></a>
## [v0.11.2] - 2021-03-13
### Bug Fixes
- `--template` and `--repository-url` flags not being used ([#119](https://github.com/git-chglog/git-chglog/issues/119))

<a name="v0.11.1"></a>
## [v0.11.1] - 2021-03-12

### Bug Fixes
- **short flags**: correctly define cli flags with shorthands ([#117](https://github.com/git-chglog/git-chglog/issues/117))

<a name="v0.11.0"></a>
## [v0.11.0] - 2021-03-12
### Features
- add Jira integration ([#52](https://github.com/git-chglog/git-chglog/issues/52))

<a name="v0.10.0"></a>
## [v0.10.0] - 2021-01-16

### BREAKING CHANGE
We adapted the naming schema of the uploaded binaries by adding a `v` prefix for tags/versions starting from `v0.10.0`: 
```
# up to 0.9.1
https://github.com/git-chglog/git-chglog/releases/download/<VERSION>/git-chglog_<OS>_<ARCH>

# from v0.10.0
https://github.com/git-chglog/git-chglog/releases/download/v<VERSION>/git-chglog_<OS>_<ARCH>

```

### Bug Fixes
- ignore only git-chglog binary in root and not subfolder

### Features
- Adds 'Custom' sort_type to CommitGroup ([#69](https://github.com/git-chglog/git-chglog/issues/69))
- enable tag_filter_pattern in config options ([#72](https://github.com/git-chglog/git-chglog/issues/72))
- switch from dep to go mod ([#85](https://github.com/git-chglog/git-chglog/issues/85))
- add option to filter commits in a case insensitive way
- add upperFirst template function
- Add emoji format and some formatters in variables


<a name="0.9.1"></a>
## [0.9.1] - 2019-09-23

<a name="0.9.0"></a>
## [0.9.0] - 2019-09-23
### Bug Fixes
- Fixing tests on windows

### Features
- Add --tag-filter-pattern flag.


<a name="0.8.0"></a>
## [0.8.0] - 2019-02-23
### Features
- add the contains, hasPrefix, hasSuffix, replace, lower and upper functions to the template functions map


<a name="0.7.1"></a>
## [0.7.1] - 2018-11-10
### Bug Fixes
- Panic occured when exec --next-tag with HEAD with tag


<a name="0.7.0"></a>
## [0.7.0] - 2018-05-06
### Bug Fixes
- Remove accidentally added `Unreleased.Tag`

### Features
- Add URL of output example for template style
- Add `--next-tag` flag (experimental)


<a name="0.6.0"></a>
## [0.6.0] - 2018-05-04
### Features
- Add tag name header id for keep-a-changelog template


<a name="0.5.0"></a>
## [0.5.0] - 2018-05-04
### Bug Fixes
- Add unreleased commits section to keep-a-changelog template [#15](https://github.com/git-chglog/git-chglog/issues/15)

### Features
- Update template format to human readable
- Add `Unreleased` field to `RenderData`


<a name="0.4.0"></a>
## [0.4.0] - 2018-04-14
### Features
- Add support for Bitbucket :tada:


<a name="0.3.3"></a>
## [0.3.3] - 2018-04-07
### Features
- Change to kindly error message when git-tag does not exist


<a name="0.3.2"></a>
## [0.3.2] - 2018-04-02
### Bug Fixes
- Fix color output bug in windows help command


<a name="0.3.1"></a>
## [0.3.1] - 2018-03-15
### Bug Fixes
- Fix preview string of commit subject ([@kt3k](https://github.com/kt3k))


<a name="0.3.0"></a>
## [0.3.0] - 2018-03-12
### Features
- Add support for GitLab :tada:


<a name="0.2.0"></a>
## [0.2.0] - 2018-03-02
### Features
- Add template for `Keep a changelog` to the `--init` option
- Supports vim like `j/k` keybind with item selection of `--init`

### Bug Fixes
- Support Windows colors :tada: ([@mattn](https://github.com/mattn))
- Fixed several bugs in Windows


<a name="0.1.0"></a>
## [0.1.0] - 2018-02-25
### Bug Fixes
- Fix error message when `Tag` can not be acquired
- Fix `Revert` of template created by Initializer

### Code Refactoring
- Refactor `Initializer` to testable

### Features
- Supports annotated git-tag and adds `Tag.Subject` field [#3](https://github.com/git-chglog/git-chglog/issues/3)
- Remove commit message preview on select format
- Add Git Basic to commit message format
- Add preview to the commit message format of `--init` option


<a name="0.0.2"></a>
## [0.0.2] - 2018-02-18
### Bug Fixes
- Fix a bug that `Commit.Revert.Header` is not converted by `GitHubProcessor`

### Features
- Add preview to the commit message format of `--init` option


<a name="0.0.1"></a>
## 0.0.1 - 2018-02-18
### Bug Fixes
- Fix parsing of revert and body

### Code Refactoring
- Fix typo
- Change to return an error if corresponding commit is empty
- Refactor the main logic

### Features
- Add cli client
- Add commits in commit version struct
- Add config normalize process
- Add Next and Previous in Tag
- Add MergeCommits and RevertCommits
- First implement


[Unreleased]: https://github.com/git-chglog/git-chglog/compare/v0.12.0...HEAD
[v0.12.0]: https://github.com/git-chglog/git-chglog/compare/v0.11.2...v0.12.0
[v0.11.2]: https://github.com/git-chglog/git-chglog/compare/v0.11.1...v0.11.2
[v0.11.1]: https://github.com/git-chglog/git-chglog/compare/v0.11.0...v0.11.1
[v0.11.0]: https://github.com/git-chglog/git-chglog/compare/v0.10.0...v0.11.0
[v0.10.0]: https://github.com/git-chglog/git-chglog/compare/0.9.1...v0.10.0
[0.9.1]: https://github.com/git-chglog/git-chglog/compare/0.9.0...0.9.1
[0.9.0]: https://github.com/git-chglog/git-chglog/compare/0.8.0...0.9.0
[0.8.0]: https://github.com/git-chglog/git-chglog/compare/0.7.1...0.8.0
[0.7.1]: https://github.com/git-chglog/git-chglog/compare/0.7.0...0.7.1
[0.7.0]: https://github.com/git-chglog/git-chglog/compare/0.6.0...0.7.0
[0.6.0]: https://github.com/git-chglog/git-chglog/compare/0.5.0...0.6.0
[0.5.0]: https://github.com/git-chglog/git-chglog/compare/0.4.0...0.5.0
[0.4.0]: https://github.com/git-chglog/git-chglog/compare/0.3.3...0.4.0
[0.3.3]: https://github.com/git-chglog/git-chglog/compare/0.3.2...0.3.3
[0.3.2]: https://github.com/git-chglog/git-chglog/compare/0.3.1...0.3.2
[0.3.1]: https://github.com/git-chglog/git-chglog/compare/0.3.0...0.3.1
[0.3.0]: https://github.com/git-chglog/git-chglog/compare/0.2.0...0.3.0
[0.2.0]: https://github.com/git-chglog/git-chglog/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/git-chglog/git-chglog/compare/0.0.2...0.1.0
[0.0.2]: https://github.com/git-chglog/git-chglog/compare/0.0.1...0.0.2
