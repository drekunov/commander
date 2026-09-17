## 1. Connector: classify symlinked directories as directories

- [x] 1.1 In `internal/connectors/filesystem/connector.go`, add a helper that returns whether an entry is a directory, resolving symlinks: return `entry.IsDir()` when true; otherwise, when `entry.Type()&os.ModeSymlink != 0`, `os.Stat(filepath.Join(path, entry.Name()))` and return `info.IsDir()` only when the stat succeeds; return false on stat error or for non-symlink non-directories
- [x] 1.2 Use that helper in `ReadDir` to set the `IsDir` attribute; keep the `Type` attribute from `entry.Type()` unchanged so links still render as `L---------`; add the `path/filepath` import

## 2. Tests

- [x] 2.1 In `internal/connectors/filesystem/connector_test.go`, add a case: a symlink whose target is a directory is reported with `IsDir=true`
- [x] 2.2 Add a case: a symlink whose target is a regular file is reported with `IsDir=false`
- [x] 2.3 Add a case: a broken symlink is reported with `IsDir=false` and `ReadDir` returns no error
- [x] 2.4 Confirm the existing header-contract test still passes with the unchanged column set (`Name`, `IsDir`, `Type`)
- [x] 2.5 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 2.6 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 2.7 Manual pty smoke (`make run`): Enter on a symlink to a directory (e.g. `/bin`) shows the target's listing and the window title shows the link path; Backspace returns to the link's parent; Enter on a symlink to a file and on a broken link does nothing and reports no error; regular directories, the other panel, and F10 quit are unaffected
