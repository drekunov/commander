## Why

The repository has no release pipeline: the only workflow lints, builds, and tests on pull requests to `develop`, so publishing a version tag produces no downloadable binaries. The commander is a cross-platform Go TUI, so tagged releases need per-OS/arch artifacts users can download.

## What Changes

- Add a GitHub Actions release workflow that runs when a version tag matching `v*` is pushed.
- Cross-compile the `commander` binary for five targets: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, and `windows/amd64`.
- Inject the tag-derived version metadata through the existing `-ldflags` variables in `internal/app/version.go` (version, commit, branch, build timestamp), reusing the Makefile's approach.
- Package each target as a `.tar.gz` (Linux/macOS) or `.zip` (Windows) archive named with the version and target, and generate a `SHA256SUMS` file covering all archives.
- Create a GitHub Release for the tag and attach every archive plus the checksums file.
- Follow the existing CI conventions: self-hosted runner and `actions/setup-go` with the project's Go version.

## Capabilities

### New Capabilities

None. This change adds CI/CD tooling and does not change application behavior, so it opts out of specs via `skip_specs: true` in `.openspec.yaml`.

### Modified Capabilities

None.

## Impact

- **Files**: new `.github/workflows/release.yml`; no Go source changes.
- **Version metadata**: relies on the ldflags variable names and package path in `internal/app/version.go` staying stable, matching `make build`.
- **Permissions**: the release job needs `contents: write` to create a GitHub Release.
- **Runners**: self-hosted, consistent with `.github/workflows/go.yml`; cross-compilation happens from a single runner, so no matrix of runner OSes is required.
- **Dependencies**: uses GitHub-provided actions only; no new Go dependencies.
