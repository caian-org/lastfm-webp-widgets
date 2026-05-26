---
name: repo-commit-and-push
description: Use when committing or pushing changes in this repository. Create focused commits that match local history, avoid unrelated user changes, and push only when explicitly requested.
---

# Repo Commit And Push

Apply this skill when a task includes committing changes or pushing to a remote.

## Policy

- Never collapse unrelated work into one commit.
- Stage only files that belong to the current delivery slice.
- Do not stage unrelated untracked files such as local environment helpers.
- Push only when the user explicitly asks for it.
- Do not rewrite history, squash, or bypass hooks unless explicitly requested.

## Commit Format

Recent history uses concise subjects:

```text
build: update go docker and release publishing
docs: simple README
feat: use the chromium binary provided at `CHROMIUM_BROWSER_BINARY_PATH` if set
refactor: export the "API" using a non-dynamic file
```

Use the same style. Prefer one of:

- `build: ...`
- `docs: ...`
- `feat: ...`
- `fix: ...`
- `refactor: ...`
- `test: ...`

Keep the subject short and imperative. Add a body only when the change is not
self-explanatory.

## Workflow

1. Review `git status --short`.
2. Inspect recent history with `git log --oneline -n 8`.
3. Group changes by coherent purpose.
4. Stage explicit paths, not `git add -A`, when unrelated files exist.
5. Commit one coherent slice at a time.
6. If the user asked for a push, push the current branch to the configured
   remote.

## Guardrails

- Do not commit `.env` or secrets.
- Do not commit generated output unless the user explicitly requests it.
- Do not include copied `cmd/*/assets/` output from `just build` unless it is
  intentionally tracked in the future.
