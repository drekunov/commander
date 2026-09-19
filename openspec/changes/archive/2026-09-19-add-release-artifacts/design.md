## Context

See proposal.md - Why. Constraints that shape the approach:

- `.github/workflows/go.yml` already runs lint, build, and test on `pull_request` to `develop`, on a `self-hosted` runner, with `actions/setup-go` and Go 1.25.
- `make build` cross-compiles nothing; it builds the host binary with `-ldflags` injecting `version`, `commit`, `branch`, and `buildUnixTimestamp` into `github.com/drekunov/gc/internal/app`. The release pipeline must use the same variable names and package path.
- `cmd/commander/main.go` is the single binary entrypoint; the app is pure Go with no cgo and no platform-specific build tags, so cross-compilation from one runner works with `GOOS`/`GOARCH`.
- Git tags are the release version source; the workflow gets the tag via `github.ref_name`.

## Goals / Non-Goals

**Goals:**

- One workflow that turns a pushed `v*` tag into a GitHub Release with downloadable per-platform archives and checksums.
- Version metadata in the released binaries matches the tag and the commit it points at.
- No changes to application code; the pipeline reuses the existing ldflags contract.

**Non-Goals:**

- No installers, package-manager manifests (deb/rpm/brew), containers, or code signing/notarization.
- No release on branch pushes or pull requests; PR CI stays in `go.yml`.
- No matrix of runner operating systems; cross-compilation from the existing self-hosted runner is sufficient.
- No changes to the version variable names, the Makefile, or any Go source.

## Decisions

### D1. Separate `release.yml` workflow

Add `.github/workflows/release.yml` rather than extending `go.yml`. Rejected extending `go.yml`: PR CI and tag releases have different triggers, permissions, and runtime cost, and mixing them makes the release path harder to reason about.

### D2. Trigger on `v*` tag push only

```yaml
on:
  push:
    tags: ["v*"]
```

Tags are immutable release anchors; branch pushes stay in PR CI. Rejected a `workflow_dispatch` input: the user chose tag-driven releases, and the tag is the unambiguous version source.

### D3. Five-target build matrix

A single job with `strategy.matrix.include` builds `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, and `windows/amd64` from one runner. `CGO_ENABLED=0` keeps the binaries static and portable. Rejected a runner-OS matrix: Go cross-compiles these targets cleanly, so extra runners add cost without benefit.

### D4. Version metadata from the tag

Build with the same `-ldflags` contract as `make build`:

- `version` = `github.ref_name` (the tag, e.g. `v1.2.3`)
- `commit` = short `github.sha`
- `branch` = `github.ref_name` (tag context; the Makefile's branch value is not meaningful for a tag)
- `buildUnixTimestamp` = `$(date -u +%s)` for reproducibility-window metadata

The release artifact name uses the tag with any leading `v` kept as-is so it matches the release.

### D5. Package per target and checksum

Each build produces `commander` (or `commander.exe`) and is packed as:

- `commander_<version>_<os>_<arch>.tar.gz` for Linux/macOS
- `commander_<version>_<os>_<arch>.zip` for Windows

The matrix job uploads each archive via `actions/upload-artifact`; a final `release` job (needs the build job) downloads them all, generates a single `SHA256SUMS` file with `sha256sum`, and publishes the GitHub Release. Rejected publishing from each matrix leg: concurrent legs would race to create the same release. Packing runs on the same cross-compile runner: `tar -czf` for the Linux/macOS targets and `zip` for the Windows target.

### D6. Publish with a release action

The release job uses `softprops/action-gh-release@v2` with `files:` pointing at the downloaded archives plus `SHA256SUMS`, and `generate_release_notes: true`. It runs with `permissions: contents: write`. Considered the `gh release create` CLI: it depends on `gh` being installed on the self-hosted runner, while the action is self-contained. Rejected a hand-rolled REST call: more code to maintain for the same result.

### D7. Pin action versions and toolchain

Use pinned action majors (`actions/checkout@v4`, `actions/setup-go@v4`, `actions/upload-artifact@v4`, `actions/download-artifact@v4`, `softprops/action-gh-release@v2`) and `go-version: '1.25'`, matching `go.yml`.

## Risks / Trade-offs

- **ldflags variable names/path drift** → the release workflow uses the exact variable names and package path from `make build`; keep the two in sync or the binaries report `undefined` versions.
- **Self-hosted runner lacks an archiver** → Linux/macOS archives use `tar` and the Windows archive uses `zip`, both run on the single build runner and present on the current self-hosted environment; no runner-OS switch is needed.
- **Matrix legs upload artifacts with duplicate names** → name each upload with the target suffix.
- **Archive name collides with an existing release asset** → the release action replaces assets on the same tag; the tag is immutable so this only happens on a re-run.
- **`permissions` default too narrow** → the release job declares `contents: write` explicitly.

## Migration Plan

No migration. Merge the workflow; the first `v*` tag after merge triggers the first release. Rollback is deleting `.github/workflows/release.yml` (and any draft release), with no effect on the application.

## Open Questions

None.
