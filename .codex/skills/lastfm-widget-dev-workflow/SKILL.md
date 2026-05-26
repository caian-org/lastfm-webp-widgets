---
name: lastfm-widget-dev-workflow
description: Local development workflow for lastfm-webp-widgets. Use when orienting in the workspace, choosing build tags, running just/go/Docker/GoReleaser commands, or validating general changes.
---

# Lastfm Widget Dev Workflow

Use this skill for repo orientation, command selection, and everyday
implementation hygiene across the workspace.

## Repo Shape

`lastfm-webp-widgets` is a single-module Go project. The main binary is
`cmd/lastfm-now-playing`, supported by repo-local packages in `internal/` and
reusable or vendored packages in `pkg/`.

- `cmd/lastfm-now-playing/` - orchestration, build-tagged entrypoints, browser
  page setup, DOM injection, screenshot capture, and animation encoding.
- `internal/lastfm/` - Last.fm service wrapper around `pkg/client/lastfm`.
- `internal/widget_browser/` - go-rod wrapper and build-tagged launcher
  implementations.
- `internal/storage/` - build-tagged disk or S3 persistence.
- `assets/` - widget HTML, JS components, CSS, and fonts.
- `pkg/client/lastfm/` - repo-owned Last.fm HTTP client.
- `pkg/webpanimation/` and `pkg/log15-2.16.0/` - vendor-style trees.

## Commands

Prefer `devbox shell` when available.

```bash
just list
just build lastfm-now-playing
just run lastfm-now-playing
go test -tags exec_local,save_s3 ./...
```

`just build` copies `assets/` into `cmd/<target>/assets` before compiling.
Avoid using it as a read-only check when the copied assets would be distracting.

## Build Tags

Always choose one tag from each pair:

- `exec_local` or `exec_lambda`
- `save_disk` or `save_s3`

The default repo command path is `exec_local,save_s3`. `go test ./...` without
tags fails because the shims in `internal/storage` and `internal/widget_browser`
need tagged implementations.

## Runtime Env

`.env` is auto-loaded by `godotenv/autoload`.

- Required: `LASTFM_USERNAME`, `LASTFM_API_KEY`.
- Required for `save_s3`: `S3_BUCKET_NAME`, `AWS_REGION`,
  `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`.
- Optional local browser override: `CHROMIUM_BROWSER_BINARY_PATH`.

Do not commit real secrets.

## Implementation Rules

- Prefer existing package boundaries over new abstractions.
- Keep build-tag variants symmetrical and compile-check the affected tag set.
- Treat `pkg/webpanimation/` and `pkg/log15-2.16.0/` as upstream/vendor code.
- Keep browser-side APIs in `assets/js/` coordinated with Go calls in
  `cmd/lastfm-now-playing/browser.go`.
- Do not hand-edit generated or local output such as `.devbox/`, `bin/`,
  copied `cmd/*/assets/`, temporary frames, release archives, or coverage data.

## Validation

Run the smallest useful checks first. Before finishing a meaningful change,
prefer:

```bash
go test -tags exec_local,save_s3 ./...
```

For runtime or Docker changes, add a targeted build check such as:

```bash
docker build --target local-runtime .
```

Only run GoReleaser commands when the task is explicitly about release behavior.

## Commit Hygiene

Recent history uses concise subjects such as `build: ...`, `docs: ...`,
`feat: ...`, and `refactor: ...`. Keep commits focused and do not stage
unrelated user changes.
