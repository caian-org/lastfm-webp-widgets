# Claude Notes

Pointer doc for Claude Code agents working in `lastfm-webp-widgets`.

**Read `AGENTS.md` first**. It is the canonical instruction layer. This file
covers only what is specific to Claude Code or to quick orientation.

## Bootstrap

1. `AGENTS.md` - canonical repository rules, structure, commands, build tags,
   testing, release, and agent-specific instructions.
2. `.codex/skills/lastfm-widget-dev-workflow/SKILL.md` - default local workflow
   guidance for this repo.
3. A subsystem skill when relevant:
   - `.codex/skills/lastfm-widget-render-pipeline/SKILL.md`
   - `.codex/skills/lastfm-widget-webp-cgo/SKILL.md`
   - `.codex/skills/repo-commit-and-push/SKILL.md`

`AGENTS.md` is authoritative. If `CLAUDE.md` disagrees with it, follow
`AGENTS.md` and reconcile this file in the same change set.

## Claude Code Specifics

- Use Claude Code's native subagent feature for delegated lanes when the session
  policy allows it. Specialists live under `.claude/agents/`, mirroring
  `.codex/agents/`.
- Skills are not duplicated for Claude. The canonical home is `.codex/skills/`;
  read skills from there.
- Claude Code settings live under `.claude/settings.json`. Hooks are currently
  empty.
- Use the local checkout as the source of truth. Do not use GitHub API reads as
  a substitute for inspecting files in this repo unless it is a narrow one-off
  check.
- Serialize shared-checkout mutations to one owner: edits, generated output,
  patch application, staging, committing, rebasing, branch switching, and
  pushing.
- After editing `AGENTS.md`, `.codex/skills/`, `.codex/agents/`, or
  `.claude/agents/`, update this file if the Claude-facing guidance changes.

## Quick Reference

Project purpose:

- Render a Last.fm "now playing" HTML widget in headless Chrome.
- Capture per-frame screenshots for light and dark themes.
- Encode animated WebP files via cgo/libwebp.
- Save output locally or upload it to S3 depending on build tags.

Command surface:

- `devbox shell` - enter the pinned development environment.
- `just list` - list binaries under `cmd/`.
- `just build lastfm-now-playing` - build the current target with
  `exec_local,save_s3`.
- `just run lastfm-now-playing` - build and run the target.
- `go test -tags exec_local,save_s3 ./...` - default validation command.

Important caveat:

- `go test ./...` without tags fails by design because `internal/storage` and
  `internal/widget_browser` dispatch to implementations selected by build tags.

Runtime envs:

- Required: `LASTFM_USERNAME`, `LASTFM_API_KEY`.
- Required for S3 builds: `S3_BUCKET_NAME`, `AWS_REGION`,
  `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`.
- Optional local browser override: `CHROMIUM_BROWSER_BINARY_PATH`.

High-traffic paths:

- `cmd/lastfm-now-playing/routine.go` - orchestration pipeline.
- `cmd/lastfm-now-playing/animation.go` - frame capture and WebP encoding.
- `cmd/lastfm-now-playing/browser.go` - page setup and DOM interactions.
- `assets/widget-now-playing.html` and `assets/js/` - browser-side widget API.
- `internal/widget_browser/` - build-tagged browser launch behavior.
- `internal/storage/` - build-tagged disk or S3 persistence.
- `pkg/webpanimation/` - vendored cgo/libwebp encoder.

Vendor boundaries:

- Do not reformat or casually refactor `pkg/log15-2.16.0/`.
- Do not reformat or casually refactor `pkg/webpanimation/`.
- `pkg/client/lastfm/` is repo-owned, not third-party vendor code.
