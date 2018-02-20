# Contributing `git-chglog`

Thank you for contributing `git-chglog` :tada:


## Templates

Please use issue/PR templates which are inserted automatically.


## Found a Bug?

If you find a bug in the source code, you can help us by [submitting an issue](https://github.com/git-chglog/git-chglog/issues) to our [GitHub Repository](https://github.com/git-chglog/git-chglog). Even better, you can submit a Pull Request with a fix.


## Commit Message Format

A format influenced by [Angular commit message](https://github.com/angular/angular/blob/master/CONTRIBUTING.md#commit-message-format).

```
<type>: <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```


### Type

Must be one of the following:

* **docs:** Documention only changes
* **ci:** Changes to our CI configuration files and scripts
* **chore:** Updating Makefile etc, no production code changes
* **feat:** A new feature
* **fix:** A bug fix
* **perf:** A code change that improves performance
* **refactor:** A code change that neither fixes a bug nor adds a feature
* **style:** Changes that do not affect the meaning of the code
* **test:** Adding missing tests or correcting existing tests


### Footer

The footer should contain a [closing reference to an issue](https://help.github.com/articles/closing-issues-via-commit-messages/) if any.

The **footer** should contain any information about **Breaking Changes** and is also the place to reference GitHub issues that this commit **Closes**.

**Breaking Changes** should start with the word `BREAKING CHANGE:` with a space or two newlines. The rest of the commit message is then used for this.
