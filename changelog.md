# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [unreleased]

- Use specific files or glob files after expression
- Allow combine short flags, e.g. -cv over -c -v

## [0.7.0] - 2023-03-05

- Consider files without extension, but with x perm set, as binary
- Replace --exclude-extension flag with --exclude using regexp

## [0.6.0] - 2022-12-25

- Update dependencies
- Add option --verbose
- Add option --exclude-extensions, $IFIND_EXCLUDE_EXT
- Add option --alias-prefix
- Add option --write-aliases

## [0.5.0] - 2022-04-12

- Skip binary files by default
- Add flag -i, --include-binary
- Improve usage and examples

## [0.4.0] - 2021-10-19

- Support VSCode editor

## [0.3.0] - 2021-10-17

- Add cmd/ifind for grepping and opening results

## [0.2.1] - 2021-06-10

- Update dependencies

## [0.2.0] - 2019-12-28

- Add type Result with a Map function

## [0.1.1] - 2019-04-11

- Minimize code complexity

## [0.1.0] - 2018-03-28

- File searching funcs By and ByName
- InFile and InStream for doing pattern search in files
