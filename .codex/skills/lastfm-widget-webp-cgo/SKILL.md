---
name: lastfm-widget-webp-cgo
description: CGO and libwebp guidance for the vendored WebP animation encoder, static Linux builds, musl toolchain, encoder settings, and Docker/GoReleaser release paths.
---

# Lastfm Widget Webp CGO

Use this skill when changing `pkg/webpanimation`, cgo settings, static linking,
Docker build stages, GoReleaser configuration, or animated WebP encoding
parameters.

## Encoder Boundary

`pkg/webpanimation/` contains Go bindings plus a vendored libwebp source tree.
It is not ordinary application code. Avoid broad formatting, file moves, or
style refactors in this tree.

The application uses:

- `webpanimation.NewWebpAnimation`
- `WebPAnimEncoderOptions.SetKmin(9)`
- `WebPAnimEncoderOptions.SetKmax(17)`
- `webpConfig.SetLossless(1)`
- 200ms frame intervals

Changing these values affects output quality, size, and animation behavior.

## Toolchain Rules

CGO is required. Linux release builds use:

```bash
CGO_ENABLED=1 CC=musl-gcc go build \
  --ldflags '-linkmode=external -extldflags="-static"' \
  -tags exec_local,save_s3
```

Docker and GoReleaser both rely on musl for static Linux binaries. Do not remove
`CGO_ENABLED=1`, `CC=musl-gcc`, or external static linking without replacing the
release validation path.

## Docker And Release

- `Dockerfile` stage `build` creates the static binary.
- `goreleaser-local-runtime` is the Docker target used by `.goreleaser.yaml`.
- `.github/workflows/release.yml` runs GoReleaser for `v*` tags.

When changing release behavior, inspect all three files together.

## Validation

For encoder or cgo changes, run:

```bash
go test -tags exec_local,save_s3 ./...
```

For release or Docker changes, add one targeted Linux build or Docker build
check. Prefer the relevant Docker target over running a full release.

## Risk Checks

- Pure-Go replacements are not drop-in unless they preserve animated WebP
  quality and timing requirements.
- macOS local compile success does not prove Linux static release success.
- `pkg/webpanimation` C sources may produce noisy diffs if auto-formatters touch
  them. Keep patches tight.
