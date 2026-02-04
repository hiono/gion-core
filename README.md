# gion-core

Shared core library for `gion` and `gionx`.

## Packages

- `planner`: computes workspace/repo diffs from desired vs actual inventory.
- `applyplan`: evaluates apply-time semantics (destructive change detection, safe branch rename checks, change counters).
- `repospec`: normalizes repo spec strings (`git@...`, `https://...`, `file://...`) into canonical repo keys.
- `repostore`: shared bare-repo store path/list utilities.
- `gitparse`: parsers for git command output (e.g., `ls-remote --symref`).
- `gitref`: remote-branch ref parsers.
- `gcplan`: helper logic for manifest-gc fetch planning.
- `importplan`: inventory reconstruction helpers used by `gion import`.
- `manifestlsplan`: computes drift categories/counts for `manifest ls`.
- `workspacerisk`: aggregate repo-level risk into workspace-level risk.

## Development

Run local checks:

- `go test ./...`
- `go vet ./...`
- `go build ./...`
