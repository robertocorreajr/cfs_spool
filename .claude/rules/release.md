# Release Notes / Changelog Rule (mandatory)

This rule is **non-optional**. Every change that produces a new tag/release
must follow the format below. The release automation (`.github/workflows/build.yml`)
fails the build if the CHANGELOG.md entry is missing.

## Source of truth

- **`CHANGELOG.md`** is the canonical, versioned history. Format:
  [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
- **`.github/RELEASE_TEMPLATE.md`** is the canonical layout for the body
  of a GitHub Release. The workflow extracts the matching version section
  from `CHANGELOG.md` and prepends the static blocks (downloads,
  installation, requirements).

## When to update the changelog

Update `CHANGELOG.md` in the **same PR** that ships a behavior-affecting
change (anything not flagged with `[skip release]` in the commit
convention). Add the entry under the `## [Unreleased]` heading using the
appropriate Keep a Changelog category:

- **Added** — new features.
- **Changed** — changes in existing functionality.
- **Deprecated** — soon-to-be removed features.
- **Removed** — now removed features.
- **Fixed** — bug fixes.
- **Security** — vulnerability fixes.

When cutting a release, promote `## [Unreleased]` to `## [X.Y.Z] - YYYY-MM-DD`
and add a new empty `## [Unreleased]` block at the top. Update the link
references at the bottom of the file.

## Writing style

- **English**, user-facing prose. Not raw commit messages.
- One bullet per visible change. Group by category, not by PR.
- Avoid Conventional Commit prefixes (`feat:`, `fix:`, `chore:`) in the
  changelog text — those belong to commit messages, not release notes.
- Mention the user-visible effect first, then (optionally) the technical
  detail. Example:
  - Bad: `feat: notifica usuário quando há nova versão (closes #21)`
  - Good: `In-app notification when a new version of CFS Spool is published.`

## When closing an issue / merging a PR

If the merged work bumps the version (no `[skip release]`):

1. Confirm `CHANGELOG.md` has an `## [Unreleased]` entry covering the change.
2. If you are also tagging the release in this PR, promote `[Unreleased]` to
   `[X.Y.Z] - YYYY-MM-DD` and update the compare links at the bottom.
3. Do not let GitHub Squash inherit `[skip release]` from internal commits.
   When merging multi-commit PRs that should produce a release, edit the
   squash commit message in the GitHub UI to remove every `[skip release]`
   marker before confirming.

## Editing past releases

When a release was published with a non-conforming body (e.g. raw commit
messages copied from the squash commit), rewrite it via:

```bash
gh release edit vX.Y.Z --notes-file path/to/notes.md
```

Use the same structure as `.github/RELEASE_TEMPLATE.md`.

## Reminder

If you (Claude Code) are about to merge a PR that introduces user-visible
behavior, **stop and update `CHANGELOG.md` before tagging**. The CI will
fail otherwise.
