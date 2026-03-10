# Contribution Guidelines

Thanks for taking the time to contribute. You're the best! Around! Nothing's gonna ever keep you down.

✨ 🎉 ✨

## Code of conduct

This project is governed by the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). We want to cultivate a respectful and constructive environment for everyone involved, and we expect all other contributors to do the same. Please use Github's reporting features if you have any concerns.

## Asking general questions

We try to answer questions in a timely manner. If you have a bug report or feature request, feel free to open an issue with as much detail as possible.

## Ways to contribute

### Getting started

Before working on an issue, fork the project and create a feature branch. Once you're done working, rebase your changes against the master branch and submit a pull request.

All pull requests must be fully tested before they're submitted for review, and please make sure your work didn't break any existing tests or reduce test coverage by running `make test` and `make test-coverage`, respectively.

We use the [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) pattern for our commit messages, which allows us to generate our [changelog](CHANGELOG.md). Commit messages should include a "type" and a short description (e.g. `fix: add helper methods to show correct trace in sentry (#85)`). For new features, it sometimes makes sense to add a longer description as well. Acceptable commit message types include:

* feat (a new feature)
* fix (a bug fix)
* perf (a performance improvement)
* refactor (a change that improves upon existing code)
* docs (a change to the docs only)
* test (a change that only includes tests)
* chore (a change to CI, etc)

Lastly, make sure your commit is signed before submitting the pull request. If you're not sure how to do this, Github has helpful [instructions](https://help.github.com/en/github/authenticating-to-github/signing-commits).

### Fixing bugs

Open a pull request with the patch, and please include a description that clearly describes the problem and solution and references a relevant issue number if applicable. Please also include a test case that covers the previous bug so that the issue won't occur again.

### Adding new features

If you have an idea for a new feature, please first discuss it with us by submitting a new issue tagged with the `enhancement` label. If we decide to adopt the feature, you're also welcome to submit a pull request if you have the time and inclination for it. Currently, we have no strict requirements around the design of new features. However, we ask that you try to follow the conventions of the current codebase, and feel free to discuss any design issues before submitting.

### Fixing whitespace, formatting code, or making purely cosmetic patches

Changes that are cosmetic in nature and do not add anything substantial to the stability, functionality, or testability of this package will generally not be accepted.

## References

These guidelines are based on examples from [Atom](https://github.com/atom/atom/blob/master/CONTRIBUTING.md) and [Ruby on Rails](https://github.com/rails/rails/blob/master/CONTRIBUTING.md).

## Thanks again for contributing!
