<a name="unreleased"></a>
## [Unreleased]


<a name="v0.14.2"></a>
## [v0.14.2] - 2021-04-14
### Bug Fixes
- add CommitGroupTitleOrder back to Options ([#143](https://github.com/git-chglog/git-chglog/issues/143))

### Chores
- **deps:** update alpine docker tag to v3.13.5 ([#144](https://github.com/git-chglog/git-chglog/issues/144))


<a name="v0.14.1"></a>
## [v0.14.1] - 2021-04-13
### Bug Fixes
- **template:** address regression in string functions for template engine ([#142](https://github.com/git-chglog/git-chglog/issues/142))

### Chores
- update changelog for v0.14.1
- add docker target to Makefile ([#138](https://github.com/git-chglog/git-chglog/issues/138))
- add make release target ([#130](https://github.com/git-chglog/git-chglog/issues/130))
- **deps:** update alpine docker tag to v3.13.4 ([#136](https://github.com/git-chglog/git-chglog/issues/136))

### Features
- add docker image on release and master ([#135](https://github.com/git-chglog/git-chglog/issues/135))


<a name="v0.14.0"></a>
## [v0.14.0] - 2021-03-28
### Chores
- update changelog for v0.14.0
- **CHANGELOG:** regenerate CHANGELOG with type-scope and KAC template ([#129](https://github.com/git-chglog/git-chglog/issues/129))

### Features
- add sprig template functions support ([#131](https://github.com/git-chglog/git-chglog/issues/131))
- add `--sort [TYPE]` flag  ([#78](https://github.com/git-chglog/git-chglog/issues/78))


<a name="v0.13.0"></a>
## [v0.13.0] - 2021-03-23
### Chores
- update changelog for v0.13.0
- use ldflags to pass version to build process ([#127](https://github.com/git-chglog/git-chglog/issues/127))

### Features
- add support for rendering .Body after .Subject as part of list ([#121](https://github.com/git-chglog/git-chglog/issues/121))


<a name="v0.12.0"></a>
## [v0.12.0] - 2021-03-20
### Chores
- update changelog for v0.12.0
- bumps version to v0.12.0
- bump golang to 1.16 ([#118](https://github.com/git-chglog/git-chglog/issues/118))
- **ci:** add golangci-lint action and apply linting changes ([#120](https://github.com/git-chglog/git-chglog/issues/120))

### Features
- allow tag sorting by semver ([#124](https://github.com/git-chglog/git-chglog/issues/124))

### BREAKING CHANGE

`JiraIssueId` has been renamed to `JiraIssueID`. This impacts the value for `pattern_maps` in `config.yml`.

* chore(ci): add golangci-lint action

* chore(lint): address errcheck lint failures

* chore(lint): address misspell lint failures

* chore(lint): address gocritic lint failures

* chore(lint): address golint lint failures

* chore(lint): address structcheck lint failures

* chore(lint): address gosimple lint failures

* chore(lint): address gofmt lint failures

* chore(ci): port to official golangci-lint github action

* Update golangci configuration for better coverage


<a name="v0.11.2"></a>
## [v0.11.2] - 2021-03-13
### Bug Fixes
- `--template` and `--repository-url` flags not being used ([#119](https://github.com/git-chglog/git-chglog/issues/119))

### Chores
- update changelog for v0.11.2
- bumps version to v0.11.2


<a name="v0.11.1"></a>
## [v0.11.1] - 2021-03-12
### Bug Fixes
- **short flags:** correctly define cli flags with shorthands ([#117](https://github.com/git-chglog/git-chglog/issues/117))

### Chores
- update readme and changelog for v0.11.1
- bumps version to v0.11.1


<a name="v0.11.0"></a>
## [v0.11.0] - 2021-03-12
### Bug Fixes
- **deps:** update all non-major dependencies ([#115](https://github.com/git-chglog/git-chglog/issues/115))
- **deps:** update module gopkg.in/kyokomi/emoji.v1 to github.com/kyokomi/emoji/v2 ([#109](https://github.com/git-chglog/git-chglog/issues/109))
- **deps:** update module github.com/urfave/cli to v2 ([#107](https://github.com/git-chglog/git-chglog/issues/107))
- **deps:** update module github.com/stretchr/testify to v1.7.0 ([#103](https://github.com/git-chglog/git-chglog/issues/103))
- **deps:** update module gopkg.in/alecaivazis/survey.v1 to github.com/AlecAivazis/survey/v2 ([#108](https://github.com/git-chglog/git-chglog/issues/108))
- **init:** support OptionAnswer form in survey/v2 ([#113](https://github.com/git-chglog/git-chglog/issues/113))

### Chores
- update changelog for v0.11.0
- bumps version to v0.11.0
- **deps:** add initial renovatebot configuration ([#102](https://github.com/git-chglog/git-chglog/issues/102))

### Features
- add Jira integration ([#52](https://github.com/git-chglog/git-chglog/issues/52))
- **flag:** --path filtering - refs ([#62](https://github.com/git-chglog/git-chglog/issues/62)). Closes [#35](https://github.com/git-chglog/git-chglog/issues/35)


<a name="v0.10.0"></a>
## [v0.10.0] - 2021-01-16
### Bug Fixes
- ignore only git-chglog binary in root and not subfolder

### Chores
- update changelog for v0.10.0
- bumps version to v0.10.0
- sorts changelog desc and excludes Merge commits
- fix Makefile typo ([#82](https://github.com/git-chglog/git-chglog/issues/82))
- **asdf:** add asdf install support to README ([#79](https://github.com/git-chglog/git-chglog/issues/79))

### Features
- Adds 'Custom' sort_type to CommitGroup ([#69](https://github.com/git-chglog/git-chglog/issues/69))
- enable tag_filter_pattern in config options ([#72](https://github.com/git-chglog/git-chglog/issues/72))
- switch from dep to go mod ([#85](https://github.com/git-chglog/git-chglog/issues/85))
- add option to filter commits in a case insensitive way
- add upperFirst template function
- Add emoji format and some formatters in variables

### Reverts
- Revert "ci: switches to personal GH Token for brew cross repo releases"
- ci: switches to personal GH Token for brew cross repo releases

### Pull Requests
- Merge pull request [#65](https://github.com/git-chglog/git-chglog/issues/65) from barryib/case-sensitive-option
- Merge pull request [#59](https://github.com/git-chglog/git-chglog/issues/59) from momotaro98/feature/add-emoji-template-in-init
- Merge pull request [#66](https://github.com/git-chglog/git-chglog/issues/66) from barryib/add-upper-first-func
- Merge pull request [#68](https://github.com/git-chglog/git-chglog/issues/68) from unixorn/tweak-readme


<a name="0.9.1"></a>
## [0.9.1] - 2019-09-23

<a name="0.9.0"></a>
## [0.9.0] - 2019-09-23
### Bug Fixes
- Fixing tests on windows

### Features
- Add --tag-filter-pattern flag.

### Pull Requests
- Merge pull request [#44](https://github.com/git-chglog/git-chglog/issues/44) from evanchaoli/tag-filter
- Merge pull request [#41](https://github.com/git-chglog/git-chglog/issues/41) from StanleyGoldman/fixing-tests-windows
- Merge pull request [#37](https://github.com/git-chglog/git-chglog/issues/37) from ForkingSyndrome/master


<a name="0.8.0"></a>
## [0.8.0] - 2019-02-23
### Features
- add the contains, hasPrefix, hasSuffix, replace, lower and upper functions to the template functions map

### Pull Requests
- Merge pull request [#34](https://github.com/git-chglog/git-chglog/issues/34) from atosatto/template-functions


<a name="0.7.1"></a>
## [0.7.1] - 2018-11-10
### Bug Fixes
- Panic occured when exec --next-tag with HEAD with tag

### Pull Requests
- Merge pull request [#31](https://github.com/git-chglog/git-chglog/issues/31) from drubin/patch-1
- Merge pull request [#30](https://github.com/git-chglog/git-chglog/issues/30) from vvakame/fix-panic


<a name="0.7.0"></a>
## [0.7.0] - 2018-05-06
### Bug Fixes
- Remove accidentally added `Unreleased.Tag`

### Chores
- Update `changelog` task in Makefile

### Features
- Add URL of output example for template style
- Add `--next-tag` flag (experimental)

### Pull Requests
- Merge pull request [#22](https://github.com/git-chglog/git-chglog/issues/22) from git-chglog/feat/add-preview-style-link
- Merge pull request [#21](https://github.com/git-chglog/git-chglog/issues/21) from git-chglog/feat/next-tag


<a name="0.6.0"></a>
## [0.6.0] - 2018-05-04
### Chores
- Update CHANGELOG template format

### Features
- Add tag name header id for keep-a-changelog template

### Pull Requests
- Merge pull request [#20](https://github.com/git-chglog/git-chglog/issues/20) from git-chglog/feat/kac-template-title-id


<a name="0.5.0"></a>
## [0.5.0] - 2018-05-04
### Bug Fixes
- Add unreleased commits section to keep-a-changelog template [#15](https://github.com/git-chglog/git-chglog/issues/15)

### Chores
- Update CHANGELOG template format

### Features
- Update template format to human readable
- Add `Unreleased` field to `RenderData`

### Pull Requests
- Merge pull request [#19](https://github.com/git-chglog/git-chglog/issues/19) from git-chglog/fix/unreleased-commits
- Merge pull request [#18](https://github.com/git-chglog/git-chglog/issues/18) from ringohub/master


<a name="0.4.0"></a>
## [0.4.0] - 2018-04-14
### Features
- Add support for Bitbucket

### Pull Requests
- Merge pull request [#17](https://github.com/git-chglog/git-chglog/issues/17) from git-chglog/feat/bitbucket


<a name="0.3.3"></a>
## [0.3.3] - 2018-04-07
### Features
- Change to kindly error message when git-tag does not exist

### Pull Requests
- Merge pull request [#16](https://github.com/git-chglog/git-chglog/issues/16) from git-chglog/fix/empty-tag-handling


<a name="0.3.2"></a>
## [0.3.2] - 2018-04-02
### Bug Fixes
- Fix color output bug in windows help command

### Pull Requests
- Merge pull request [#14](https://github.com/git-chglog/git-chglog/issues/14) from git-chglog/fix/windows-help-color


<a name="0.3.1"></a>
## [0.3.1] - 2018-03-15
### Bug Fixes
- fix preview string of commit subject

### Pull Requests
- Merge pull request [#13](https://github.com/git-chglog/git-chglog/issues/13) from kt3k/feature/fix-preview


<a name="0.3.0"></a>
## [0.3.0] - 2018-03-12
### Chores
- Add helper task for generate CHANGELOG

### Features
- Add support for GitLab

### Pull Requests
- Merge pull request [#12](https://github.com/git-chglog/git-chglog/issues/12) from git-chglog/feat/gitlab


<a name="0.2.0"></a>
## [0.2.0] - 2018-03-02
### Chores
- Fix release flow (retry)
- Add AppVeyor config

### Features
- Add template for `Keep a changelog` to the `--init` option
- Supports vim like `j/k` keybind with item selection of `--init`

### Pull Requests
- Merge pull request [#11](https://github.com/git-chglog/git-chglog/issues/11) from git-chglog/develop
- Merge pull request [#10](https://github.com/git-chglog/git-chglog/issues/10) from mattn/fix-test
- Merge pull request [#9](https://github.com/git-chglog/git-chglog/issues/9) from mattn/windows-color


<a name="0.1.0"></a>
## [0.1.0] - 2018-02-25
### Bug Fixes
- Fix a bug that `Commit.Revert.Header` is not converted by `GitHubProcessor`
- Fix error message when `Tag` can not be acquired
- Fix `Revert` of template created by Initializer

### Chores
- Fix release scripts
- Remove unnecessary task
- Add coverage measurement task for local confirmation
- Change release method of git tag on TravisCI

### Code Refactoring
- Refactor `Initializer` to testable

### Features
- Supports annotated git-tag and adds `Tag.Subject` field [#3](https://github.com/git-chglog/git-chglog/issues/3)
- Remove commit message preview on select format
- Add Git Basic to commit message format
- Add preview to the commit message format of `--init` option

### Pull Requests
- Merge pull request [#8](https://github.com/git-chglog/git-chglog/issues/8) from git-chglog/feat/0.0.3
- Merge pull request [#6](https://github.com/git-chglog/git-chglog/issues/6) from git-chglog/chore/coverage
- Merge pull request [#4](https://github.com/git-chglog/git-chglog/issues/4) from paralax/patch-1
- Merge pull request [#5](https://github.com/git-chglog/git-chglog/issues/5) from git-chglog/develop
- Merge pull request [#1](https://github.com/git-chglog/git-chglog/issues/1) from git-chglog/develop


<a name="0.0.2"></a>
## [0.0.2] - 2018-02-18
### Chores
- Fix release script
- Add release process


<a name="0.0.1"></a>
## 0.0.1 - 2018-02-18
### Bug Fixes
- Fix parsing of revert and body

### Chores
- Fix timezone in TravisCI
- Add travis configuration
- Add Makefile for task management
- Fix testcase depending on datetime
- Update vendor packages
- Add e2e tests
- Setup gitignore
- Initial commit
- **editor:** Add Editorconfig

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


[Unreleased]: https://github.com/git-chglog/git-chglog/compare/v0.14.2...HEAD
[v0.14.2]: https://github.com/git-chglog/git-chglog/compare/v0.14.1...v0.14.2
[v0.14.1]: https://github.com/git-chglog/git-chglog/compare/v0.14.0...v0.14.1
[v0.14.0]: https://github.com/git-chglog/git-chglog/compare/v0.13.0...v0.14.0
[v0.13.0]: https://github.com/git-chglog/git-chglog/compare/v0.12.0...v0.13.0
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
