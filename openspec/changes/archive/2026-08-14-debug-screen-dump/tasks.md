## 1. Dumper implementation (internal/ui/dump.go)

- [x] 1.1 Add `internal/ui/dump.go` defining the `dumper` interface (`Dump(frame string) error`) and `fileDumper` with an injectable `now func() time.Time`
- [x] 1.2 Implement `dumpFilename(now time.Time) string` producing `screen-20060102-150405.ansi` (second-precision local time)
- [x] 1.3 Implement `fileDumper.Dump` writing the raw frame bytes to the working directory, appending `-1`, `-2`, ... to the base name until the target is free so successive dumps never overwrite each other

## 2. Key handler wiring (internal/ui/ui.go)

- [x] 2.1 Add a `dumpScreenKey` constant (F12) and a `dumper` field on `Model`, defaulted to `&fileDumper{}` in `New`
- [x] 2.2 In `Model.Update`, on `tea.KeyMsg` of type `dumpScreenKey`: capture `m.View()`, return a `tea.Cmd` that calls `dumper.Dump(frame)` and swallows the error, leaving F10 quit and all other key handling unchanged
- [x] 2.3 Verify `go build` passes and `app.UI` interface is untouched

## 3. Tests

- [x] 3.1 Unit-test `dumpFilename` (fixed clock → exact name) and collision suffixing (same-second presses produce distinct files, first not overwritten)
- [x] 3.2 Unit-test `fileDumper.Dump` against `t.TempDir()` working directory (bytes written, ANSI preserved)
- [x] 3.3 Unit-test the F12 handler with a fake `dumper`: executing the returned `tea.Cmd` records the frame, and the recorded frame equals `Model.View()` at keypress time
- [x] 3.4 Run `go test -race ./...`, `go vet ./...`, and `make lint`; confirm all pass
- [x] 3.5 Manual smoke with `make run`: press F12, confirm `screen-*.ansi` appears in the working directory and `cat` replays the frame with colors; press F12 again, confirm a distinct file
