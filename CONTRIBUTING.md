# Contributing to KC SSH PAM

We want to make contributing to this project as easy and transparent as
possible.

## Project structure

- `internal` - Contains internal packages
  - `internal/auth` - Helper functions related to authentication.
  - `internal/conf` - Functions related to configuration file.
  - `internal/flags` - Commandline flags.


## Commits

Commit messages should be well formatted, and to make that "standardized", we
are using Conventional Commits.

```shell

  <type>[<scope>]: <short summary>
     │      │             │
     │      │             └─> Summary in present tense. Not capitalized. No
     |      |                 period at the end. 
     │      │
     │      └─> Scope (optional): eg. common, compiler, authentication, core
     │
     └─> Type: chore, docs, feat, fix, refactor, style, or test.
     
```

You can follow the documentation on
[their website](https://www.conventionalcommits.org).

## Pull Requests

We actively welcome your pull requests.

1. Fork the repo and create your branch from `master`.
2. If you've added code that should be tested, add tests.

## Issues

We use GitHub issues to track public bugs. Please ensure your description is
clear and has sufficient instructions to be able to reproduce the issue.

## License

By contributing to KC SSH PAM, you agree that your contributions will be licensed
under the LICENSE file in the root directory of this source tree.
