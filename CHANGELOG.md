# CHANGELOG


<a name="0.1.0"></a>
## [0.1.0](https://github.com/git-chglog/git-chglog/compare/0.0.2...0.1.0) (2018-02-25)

### Bug Fixes

* Fix a bug that `Commit.Revert.Header` is not converted by `GitHubProcessor`
* Fix error message when `Tag` can not be acquired
* Fix `Revert` of template created by Initializer

### Code Refactoring

* Refactor `Initializer` to testable

### Features

* Supports annotated git-tag and adds `Tag.Subject` field [#3](https://github.com/git-chglog/git-chglog/issues/3)
* Remove commit message preview on select format
* Add Git Basic to commit message format
* Add preview to the commit message format of `--init` option


<a name="0.0.2"></a>
## [0.0.2](https://github.com/git-chglog/git-chglog/compare/0.0.1...0.0.2) (2018-02-19)

### Bug Fixes

* Fix `Revert` of template created by Initializer

### Features

* Add Git Basic to commit message format
* Add preview to the commit message format of `--init` option


<a name="0.0.1"></a>
## 0.0.1 (2018-02-18)

First release :tada:
