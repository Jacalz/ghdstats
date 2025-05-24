# Changelog

## v1.2.0
- Improved error handling for failures during processing.
- Show an helpful error when GitHub API rate limit is exceeded.
- Add instrumentation for collecting pprof profiles for use with PGO.
- Fetching a single repository no longer spins up an unnecessary goroutine.
- Release binaries for x86_64 are now built with GOAMD64=v3 and PGO enabled.
- Added GitHub workflows for static analysis.
- Various minor code cleanups.

## v1.1.0
- Fix indentation issues between columns (issue #3).
- Print more helpful error messages when processing fails.
- Exit cleanly instead of calling `panic()` on invalid data.
