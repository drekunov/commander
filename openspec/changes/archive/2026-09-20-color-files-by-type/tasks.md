## 1. File-kind classification

- [x] 1.1 Add the file-kind token constants (`directory`, `symlink`, `executable`, `image`, `archive`, `source`, `config`) to `internal/app` for the connector and panel to share.
- [x] 1.2 In `internal/connectors/filesystem/connector.go`, add a hidden `Kind` attribute to each header/entry row alongside the existing hidden `IsDir`.
- [x] 1.3 Implement classification with the design D2 precedence (directory → symlink → executable → extension category → empty) using `os.DirEntry` type/mode and the name suffix.
- [x] 1.4 Add the case-insensitive extension sets from design D3 for `image`, `archive`, `source`, and `config`.
- [x] 1.5 Connector test: directories, symlinks, executables, one file per extension category, and an unclassified file each carry the expected hidden `Kind`.

## 2. Theme classes

- [x] 2.1 Add `FileDirectoryStyle`, `FileExecutableStyle`, `FileSymlinkStyle`, `FileImageStyle`, `FileArchiveStyle`, `FileSourceStyle`, and `FileConfigStyle` to `config.Styles` (`internal/config/config.go`).
- [x] 2.2 Load them from `.file-directory`, `.file-executable`, `.file-symlink`, `.file-image`, `.file-archive`, `.file-source`, and `.file-config` in `LoadStyles`.
- [x] 2.3 Add the seven classes with the design D8 foreground values and the directory's bold weight to `internal/config/styles.css`, leaving row backgrounds untouched.
- [x] 2.4 Config test: the seven selectors load into their fields, a working-directory override changes one, and a stylesheet omitting a class leaves that field zero while the rest load.

## 3. Panel Name-cell coloring

- [x] 3.1 In `internal/ui/widgets/panel/panel.go`, map a kind token to the injected field with a switch (design D7), returning a zero style for absent/unknown kinds.
- [x] 3.2 Record each display row's kind when rows are built (including `directory` for the synthetic `..` row), so the render step can align kinds with rendered lines.
- [x] 3.3 In `View()`, after rendering the table and its table-background wrap, splice a foreground/bold-only SGR around the Name column's visible span per visible data row, using `charmbracelet/x/ansi` for ANSI-aware offsets (design D4).
- [x] 3.4 Bound the candidate window from the table's `Cursor()` and `Height()` and match rendered lines to rows by their plain text (design D4), and skip the focused cursor row (design D5); skip the splice entirely for zero styles (design D6).
- [x] 3.5 Close each splice by resetting only the properties the kind style set, never with `\x1b[0m`, so the table background and cursor highlight survive.
- [x] 3.6 Panel test: a focused panel colors each kind's Name cell with its injected style, leaves Size/Date/Time and the header untouched, and keeps column alignment identical to an unstyled render.
- [x] 3.7 Panel test: the focused cursor row is uniform (no kind color) and a blurred panel shows kind colors on every row.
- [x] 3.8 Panel test: an entry with no kind, an unknown kind, or a missing class renders default; a listing with no `Kind` attribute splices nothing.
- [x] 3.9 Panel test: a scrolled listing colors the correct visible rows after the cursor moves below the fold.

## 4. Verification

- [x] 4.1 Run `go build ./...`, `go vet ./...`, and `go test -race ./...` and keep them green.
- [x] 4.2 Run `make lint` and confirm golangci-lint passes.
- [ ] 4.3 Manual smoke via `make run`: confirm directories, executables, symlinks, and category files show distinct Name colors, the focused cursor row stays uniform, a blurred panel shows the colors, and a `styles.css` override restyles a kind.
