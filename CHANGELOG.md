# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Types of changes

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.

## [0.3.0]

- `Added` persistance feature with `--persist` flag.
- `Added` computation on disk instead of memory with `--diskstorage` flag.
- `Added` analysis of deep nested structures (arrays and objects).
- `Added` validated constraints use a different shade of green in HTML report.
- `Added` possibility to use template string to generate coherent source with `coherentSource` parameter.

## [0.2.1]

- `Fixed` HTML report fail when rating result is NaN.

## [0.2.0]

- `Added` configuration file with `metrics[].exclude`, `metrics[]coherentWith` and `metrics[]constraints` parameters.

## [0.1.0]

- `Added` first official version.
