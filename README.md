[![CI][ci-shield]][ci-url]
[![Release][rel-shield]][rel-url]
[![GitHub tag][tag-shield]][tag-url]

# `lastfm-webp-widgets`

> animated WebP "now playing" widgets powered by Last.fm

This repo renders Last.fm "now playing" widgets in headless Chrome, captures
per-frame screenshots, and encodes animated WebP files for light and dark
themes. The output can be saved locally or uploaded to S3.

[ci-shield]: https://img.shields.io/github/actions/workflow/status/caian-org/lastfm-webp-widgets/ci.yml?label=ci&logo=github&style=flat-square
[ci-url]: https://github.com/caian-org/lastfm-webp-widgets/actions/workflows/ci.yml
[rel-shield]: https://img.shields.io/github/actions/workflow/status/caian-org/lastfm-webp-widgets/release.yml?label=release&logo=github&style=flat-square
[rel-url]: https://github.com/caian-org/lastfm-webp-widgets/actions/workflows/release.yml
[tag-shield]: https://img.shields.io/github/tag/caian-org/lastfm-webp-widgets.svg?logo=git&logoColor=FFF&style=flat-square
[tag-url]: https://github.com/caian-org/lastfm-webp-widgets/releases


## Widgets

### Now Playing

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://caian-org.s3.amazonaws.com/lastfm-now-playing-dark.webp" width="600px">
  <img src="https://caian-org.s3.amazonaws.com/lastfm-now-playing-light.webp" width="500px">
</picture>


## Quick Start

### Run with Docker

```bash
docker run \
    --rm \
    -e LASTFM_USERNAME='your-lastfm-username' \
    -e LASTFM_API_KEY='your-lastfm-api-key' \
    -e S3_BUCKET_NAME='an-s3-bucket-you-own' \
    -e AWS_ACCESS_KEY_ID='access-key-to-upload-to-the-bucket' \
    -e AWS_SECRET_ACCESS_KEY='secret-key-to-upload-to-the-bucket' \
    -e AWS_REGION='us-east-1' \
    ghcr.io/caian-org/lastfm-widget-now-playing:latest
```

### Build from source

```bash
just run lastfm-now-playing
```

Required env vars:
- `LASTFM_USERNAME`, `LASTFM_API_KEY` — for all runs.
- `S3_BUCKET_NAME`, `AWS_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` — when built with `save_s3`.
- `CHROMIUM_BROWSER_BINARY_PATH` (optional) — bypass go-rod's Chromium auto-download for local runs.


## Development

```bash
devbox shell        # enter pinned dev environment
just list           # list available targets under cmd/
just build lastfm-now-playing
just test           # go test -tags exec_local,save_s3 ./...
just cover          # coverage profile
just lint           # go vet with build tags
```

See `AGENTS.md` for the build-tag matrix (`exec_local`/`exec_lambda`,
`save_disk`/`save_s3`) and vendor boundaries.


## License

To the extent possible under law, [Caian Ertl][me] has waived __all copyright
and related or neighboring rights to this work__. In the spirit of _freedom of
information_, I encourage you to fork, modify, change, share, or do whatever
you like with this project! [`^C ^V`][kopimi]

[![License][cc-shield]][cc-url]

[me]: https://github.com/upsetbit
[cc-shield]: https://forthebadge.com/images/badges/cc-0.svg
[cc-url]: http://creativecommons.org/publicdomain/zero/1.0
[kopimi]: https://kopimi.com
