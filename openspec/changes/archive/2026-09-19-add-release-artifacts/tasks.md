## 1. Release workflow

- [x] 1.1 Create `.github/workflows/release.yml` with `name: Release` and trigger `on: push: tags: ["v*"]`.
- [x] 1.2 Add the build job on `runs-on: self-hosted` with `actions/checkout@v4` and `actions/setup-go@v4` (`go-version: '1.25'`), matching `go.yml`.
- [x] 1.3 Define the build matrix for `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, and `windows/amd64` using `strategy.matrix.include` with `goos` and `goarch` fields.
- [x] 1.4 Build `./cmd/commander` with `CGO_ENABLED=0`, `GOOS`/`GOARCH` from the matrix, and the ldflags contract from `make build` (`-X` for `version`, `commit`, `branch`, `buildUnixTimestamp` in `github.com/drekunov/gc/internal/app`), using `github.ref_name` for the version and short `github.sha` for the commit.
- [x] 1.5 Name the binary `commander` (`.exe` on Windows) and produce an archive named `commander_<version>_<goos>_<goarch>.tar.gz` for Linux/macOS or `.zip` for Windows.
- [x] 1.6 Upload each archive with `actions/upload-artifact@v4` using a target-suffixed artifact name so matrix legs do not collide.

## 2. Release job

- [x] 2.1 Add a `release` job that `needs` the build job, runs on `self-hosted`, and declares `permissions: contents: write`.
- [x] 2.2 Download all build artifacts with `actions/download-artifact@v4`.
- [x] 2.3 Generate a `SHA256SUMS` file covering every downloaded archive.
- [x] 2.4 Publish the GitHub Release for the tag with `softprops/action-gh-release@v2`, attaching all archives and `SHA256SUMS`, with `generate_release_notes: true`.

## 3. Verification

- [x] 3.1 Validate the workflow YAML syntax (for example with `actionlint`, or by a YAML parse) and confirm the trigger, matrix, and permissions are as specified.
- [x] 3.2 Dry-run the build steps locally for each matrix target (`CGO_ENABLED=0 GOOS=<os> GOARCH=<arch> go build ...`) and confirm every binary is produced.
- [x] 3.3 Confirm the ldflags strings match the variable names and package path in `internal/app/version.go` and `make build`.
- [ ] 3.4 After merge, push a throwaway `v*` tag (or run on a fork) and confirm the GitHub Release contains five archives and `SHA256SUMS`; run a released binary and confirm its startup log line (`getBuildInfo`, logged from `internal/app/app.go`) reports the tag as the version.
