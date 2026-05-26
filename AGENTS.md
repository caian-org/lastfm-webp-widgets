# Repository Guidelines

## Project Structure & Module Organization

`lastfm-webp-widgets` is a Go module that renders a Last.fm "now playing" widget
inside headless Chrome, captures per-frame screenshots, and encodes two animated
WebP files: `lastfm-now-playing-light.webp` and
`lastfm-now-playing-dark.webp`. The binary can run locally, in a container, or as
an AWS Lambda container.

- `cmd/lastfm-now-playing/` contains the executable, orchestration pipeline,
  build-tagged entrypoints, browser setup, DOM injection, and animation flow.
- `internal/lastfm/`, `internal/widget_browser/`, `internal/storage/`,
  `internal/util/`, `internal/logger/`, and `internal/constants/` hold repo-local
  implementation packages.
- `pkg/client/lastfm/` is the in-repo Last.fm HTTP client.
- `pkg/webpanimation/` is the vendored cgo/libwebp animation encoder.
- `pkg/log15-2.16.0/` is vendored `log15`.
- `assets/` contains the widget HTML, JavaScript components, CSS, and fonts.
- `.github/workflows/release.yml`, `.goreleaser.yaml`, `Dockerfile`,
  `docker-compose.yml`, `entry_script.sh`, `devbox.json`, and `.justfile`
  define build, release, and runtime workflows.

Treat `pkg/log15-2.16.0/` and `pkg/webpanimation/` as upstream/vendor code. Do
not reformat, broad-refactor, or modernize those trees unless the task is
explicitly about vendor maintenance.

## Build, Test, and Development Commands

Use `devbox shell` when possible. It pins Go, `just`, and GoReleaser.

- `just list` lists binaries under `cmd/`.
- `just build <target>` copies `assets/` into `cmd/<target>/assets` and builds
  `cmd/<target>` into `bin/<target>`.
- `just build-all` builds every target under `cmd/`.
- `just run <target>` builds and runs `bin/<target>`.
- `go test -tags exec_local,save_s3 ./...` is the default validation command.
- `go test ./...` fails by design because build-tag dispatch shims reference
  implementations selected only when an `exec_*` tag and a `save_*` tag are set.

The only current target is `lastfm-now-playing`. The `.justfile` always builds
with `-tags exec_local,save_s3`; use direct `go build` commands for other
variants, for example `-tags exec_local,save_disk`.

## Runtime Configuration

The repo loads `.env` automatically via `github.com/joho/godotenv/autoload`.
Copy `.env.model` for local setup.

- Required for all runs: `LASTFM_USERNAME`, `LASTFM_API_KEY`.
- Required for `save_s3`: `S3_BUCKET_NAME`, `AWS_REGION`,
  `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`.
- Optional for local browser runs: `CHROMIUM_BROWSER_BINARY_PATH`, used by
  `exec_local` to avoid go-rod auto-downloading Chromium.

Do not commit real secrets. `.env` is ignored.

## Build Tags And Variants

Two orthogonal build-tag pairs control behavior:

- `exec_local` vs. `exec_lambda` selects the executable entrypoint and browser
  launcher. Local mode uses go-rod's launcher or `$CHROMIUM_BROWSER_BINARY_PATH`;
  Lambda mode uses `/opt/google/chrome/chrome` with Lambda-safe Chrome flags.
- `save_disk` vs. `save_s3` selects persistence. Disk mode writes WebP bytes to
  files; S3 mode uploads public `image/webp` objects to `$S3_BUCKET_NAME`.

`internal/widget_browser/main.go` and `internal/storage/main.go` are stable
dispatch shims. Change browser flags or storage behavior in the tagged sibling
files, not by bypassing the shims.

## Rendering Pipeline

The main flow is `cmd/lastfm-now-playing/routine.go::doRoutine`:

1. Fetch Last.fm user info and recent tracks.
2. Open `assets/widget-now-playing.html` in a headless browser at 1130 x 348
   with `WIDGET_PIXEL_RATIO=2`.
3. Inject user, track, album art, scrobble count, listening state, and theme
   data through JavaScript globals exposed by `assets/js/`.
4. Detect oversized track titles with
   `music.getTrackTitleSizeInPixelsRounded()` and switch to scrolling text when
   needed.
5. Capture light and dark frame PNGs under a temporary `/tmp/frames-*`
   directory.
6. Encode frames through `pkg/webpanimation` as lossless animated WebP at
   200ms/frame with KMin=9 and KMax=17.
7. Persist `lastfm-now-playing-light.webp` and
   `lastfm-now-playing-dark.webp`.

Keep HTML/JS/CSS changes coordinated with Go DOM calls. The Go code assumes the
page exposes `music`, `user`, `waves`, and `theme` globals.

## CGO, Docker, And Release

`pkg/webpanimation` compiles a vendored libwebp source tree through cgo, so builds
require `CGO_ENABLED=1` and a C toolchain. Static Linux release builds use
`CC=musl-gcc` with external linking.

The Dockerfile contains runtime-specific stages:

- `build` builds the static Linux binary.
- `lambda-runtime-base` installs Chrome in the AWS Lambda AL2 base.
- `lambda-runtime-local` adds the Lambda Runtime Interface Emulator.
- `lambda-runtime-aws` is the production Lambda image.
- `local-runtime` is an Alpine + Chromium local container runtime.
- `goreleaser-local-runtime` is used by GoReleaser's Docker publishing path.

Releases are tag-driven. Pushing a `v*` tag runs `.github/workflows/release.yml`,
which invokes GoReleaser and publishes GHCR images.

## Coding Style & Safety

- Prefer small, direct Go changes that match existing package boundaries.
- Keep build-tag variants symmetrical and compile-check the tag combination you
  changed.
- Avoid broad rewrites in `assets/` unless the task is about visual behavior.
- Do not edit generated or local-runtime output by hand: `.devbox/`, `bin/`,
  copied `cmd/*/assets/`, temporary frames, release archives, and coverage data.
- If a task requires committing, use focused commits that match repo history,
  such as `build: ...`, `docs: ...`, `feat: ...`, `fix: ...`, or `refactor: ...`.

## Agent-Specific Instructions

`AGENTS.md` is the canonical instruction layer for this repository. `CLAUDE.md`
is a Claude Code pointer and quick reference. Local Codex skills live in
`.codex/skills/`; Claude Code should read those skills from `.codex/skills/`
rather than duplicating them under `.claude/skills/`. Specialist subagents are
mirrored under `.codex/agents/` and `.claude/agents/`.

When changing repository guidance, keep `AGENTS.md`, `CLAUDE.md`, skills, and
subagents consistent in the same change set.

## BRAIN - Obsidian Vault

Vault root: `/Users/upsetbit/Projects/_me/upsetbit/BRAIN`.

Use this personal Obsidian vault when the user asks to save to BRAIN, put in
BRAIN, "guarde no brain", or clearly means storing a note, capture, or memo in
the vault rather than in this project tree.

- Default generic captures go under `inbox/`.
- Investigation work goes under `projects/<slug>/` with frontmatter.
- `knowledge/` is read-only. To propose durable knowledge, write a proposal
  under `inbox/` with `type: proposal`.
- Daily notes under `daily/` are append-only.
- Prefer Markdown. If no filename is specified, choose a sensible slug.
- Include YAML frontmatter with `type`, `created`, `agent`, and `source`.
- Never relocate the vault or assume a different root without explicit user
  instruction.
- Do not overwrite important notes without confirmation if the target path
  already exists and has content.
