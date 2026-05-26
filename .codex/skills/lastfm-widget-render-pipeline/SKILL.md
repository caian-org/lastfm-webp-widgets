---
name: lastfm-widget-render-pipeline
description: Rendering pipeline guidance for Last.fm fetches, headless Chrome DOM injection, widget assets, frame capture, theme variants, scroll behavior, and output persistence.
---

# Lastfm Widget Render Pipeline

Use this skill when changing the widget generation path: Last.fm data, browser
startup, DOM calls, HTML/JS/CSS assets, screenshots, animation timing, themes,
or storage keys.

## Pipeline

The pipeline starts in `cmd/lastfm-now-playing/routine.go::doRoutine`:

1. `internal/lastfm` fetches `user.getinfo` and `user.getrecenttracks`.
2. `initBrowser` opens `assets/widget-now-playing.html` through go-rod.
3. Go injects user and track state by calling JavaScript globals.
4. Go checks whether the title exceeds `TRACK_TITLE_MAX_SIZE_PIXELS`.
5. `animation.go` captures per-theme PNG frames into `/tmp/frames-*`.
6. `pkg/webpanimation` encodes light and dark WebPs.
7. `internal/storage.Save` persists the WebP bytes.

## Browser-Side Contract

The HTML page must expose these globals:

- `music`
- `user`
- `waves`
- `theme`

Before renaming JS functions or changing layout measurements, inspect
`cmd/lastfm-now-playing/browser.go` and the relevant file under `assets/js/`.

## Scroll And Timing Rules

- Non-scrolling output uses 10 frames.
- Scrolling output advances title ticks every 3 randomized sound-wave frames and
  adds a tail of randomized sound-wave frames.
- Both light and dark themes are captured for each logical frame.
- WebP frame timing is currently 200ms/frame.

Preserve visual parity between themes unless the task explicitly changes one
theme.

## Storage Keys

The public output names are:

- `lastfm-now-playing-light.webp`
- `lastfm-now-playing-dark.webp`

Treat these as user-facing because README examples and S3 consumers may depend
on them.

## Risk Checks

- Empty recent-track responses should exit early.
- Last.fm image arrays are indexed for the largest image; check bounds if
  changing response handling.
- Browser launch behavior differs between `exec_local` and `exec_lambda`.
- CSS or font changes can alter overflow detection and scroll frame counts.

## Validation

Run:

```bash
go test -tags exec_local,save_s3 ./...
```

When changing visual behavior, prefer one manual run with a real `.env` and
inspect the generated light and dark WebP outputs.
