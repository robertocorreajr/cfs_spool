# Release Asset Retention Policy

This document defines how CFS Spool retains the binary assets attached to
GitHub Releases. It is the canonical source linked from
[`CLAUDE.md`](../CLAUDE.md) and [`.claude/rules/release.md`](../.claude/rules/release.md).

> **TL;DR** — We keep downloads only for recent stable releases plus a
> short list of historic milestones. We **never** delete tags, release
> entries, or release notes. Anything older can be rebuilt from source by
> checking out the tag.

## Goals

1. Keep the repository light and the Releases page easy to browse.
2. Always offer downloads for current and recent versions.
3. Preserve the project's full historical record (tags + release notes).
4. Make removals intentional, scoped, and reversible (rebuild from tag).

## What we keep

### Stable releases

A release is considered **stable** when its tag matches the strict pattern
`vX.Y.Z` with no pre-release suffix (no `-beta`, `-rc`, `-test`, `-fix`,
`-restore`, etc.).

| Rule | Action |
|------|--------|
| Latest stable release | **Always keep all assets**. |
| Previous 4 stable releases (5 most recent total, including latest) | Keep all assets. |
| Older stable releases | Assets may be retired. Tag and notes are preserved. |

The 5-release window is intentional: we currently ship patches in fast
batches, so 5 stables typically span the active maintenance horizon. The
window may be revisited (and this document updated) when release cadence
changes materially.

### Milestone releases (allowlist)

The following releases are considered project milestones. Their assets
**must always be preserved**, regardless of age, in addition to the
rolling 5-stable window:

| Tag | Why it's a milestone |
|-----|----------------------|
| `v1.0.0` | First public release of CFS Spool. |
| `v2.1.0` | First stable on the 2.x line. |
| `v3.0.0` | Full rewrite to Wails v2 + React + shadcn/ui (native desktop app). |

Future maintainers add to this table when cutting a release that
qualifies (major version bump, signing/notarization milestone, first
release on a new platform, etc.). Edits to this list **must be made in
the same PR** that introduces the milestone.

### Pre-release / experimental tags

Tags that do **not** match strict `vX.Y.Z` are considered pre-releases or
experimental:

- `v1.1.0-beta6`, `v1.1.0-beta7`
- `v2.1.3-test`, `v2.1.8-restore-working`
- `v2.1.10-color-black-fix`, `v2.1.11-black-color-fix`,
  `v2.1.13-validation-fix`
- any future `-beta*`, `-rc*`, `-test*`, `-fix*`, `-restore*` tag

Their assets may be retired as soon as the work they exercise has been
folded into a stable `vX.Y.Z` release. The tag and release entry stay.

## What we never delete

- **Tags** — these are immutable historical anchors and the only way to
  rebuild old binaries from source.
- **Release entries** — the GitHub Release page itself is part of the
  project's audit trail.
- **Release notes** — the body of every release stays intact (and is
  amended with the retired-asset note when applicable).
- **Anything from the `main` branch** — this policy applies only to
  uploaded artifacts attached to Releases.

## Standard note for releases whose assets were retired

When a release has its binaries removed under this policy, append the
following block to the bottom of the release body (one blank line above
it). The build workflow does not regenerate this section, so it is safe
to add manually or via the cleanup workflow:

```markdown
---

## Downloads retired

The binaries originally attached to this release were removed as part of
the project's [asset retention policy](../../blob/main/.github/RELEASE_RETENTION.md).
Please download the latest release at
https://github.com/robertocorreajr/cfs_spool/releases/latest. The tag for
this version is still available if you need to rebuild from source:

```bash
git checkout vX.Y.Z
wails build
```
```

Replace `vX.Y.Z` with the actual tag. The relative `../../blob/main/...`
link resolves correctly from any release page on GitHub.

## Cleanup process

### Automatic on every new release

Cleanup runs automatically on every `release:published` event. Each new
release that the Build & Release workflow publishes immediately triggers
the retention workflow in `apply` mode, which retires assets from any
release that falls outside the protection rules above.

This is the right default because:

1. The trigger is the act of publishing a new version — a deliberate,
   audited event — not an arbitrary clock tick.
2. The protection rules already prevent the latest stable, the recent
   stable window, and milestones from ever being touched, so an
   automatic run cannot retire something that should still be available.
3. There is no manual toil: the Releases page stays in compliance with
   the policy without anyone remembering to run anything.

### Manual override

The same workflow also accepts `workflow_dispatch` for ad-hoc runs (e.g.
auditing the plan, retroactive cleanup, testing changes to the policy):

- `dry_run` (default `true`) — when `true`, prints the plan without
  modifying anything.
- `confirm` — must be set to the literal string `apply` to actually
  delete assets when `dry_run=false`.

When triggered by `release:published`, those inputs are ignored and the
workflow always runs in `apply` mode.

### Workflow behavior

The workflow ([`.github/workflows/release-retention.yml`](workflows/release-retention.yml)):

1. Lists every release via `gh release list --limit 200`.
2. Classifies each release as `stable`, `pre-release`, or `milestone`.
3. Picks the 5 most recent stable releases (by published date) and
   protects them.
4. Protects every tag listed in the milestone allowlist.
5. For all other releases, lists their assets, deletes each via
   `gh release delete-asset`, and appends the standard retired-assets
   note to the release body if not already present.
6. Never invokes `gh release delete` or `git tag -d`.

### First run

The very first run after merging this policy will retire roughly 25
releases worth of binaries. That is the intended one-time correction to
align history with the policy. If you want to audit the plan before that
happens, run the workflow once via `workflow_dispatch` with
`dry_run=true` **before merging** (or before the next release publishes).

### Manual restoration

If an asset is needed after retirement, rebuild from the tag:

```bash
git fetch --tags
git checkout vX.Y.Z
wails build
```

The signed/notarized macOS DMG cannot be reproduced bit-for-bit without
the Apple credentials, but a functional unsigned build is always
recoverable from source.

## Future work

- Generate a small badge / status comment on releases whose assets were
  retired (current implementation appends the standard note and stops
  there).
- Add a Slack/Discord notification when the workflow retires assets so
  the project owner has an external paper trail.

## Change history

- **2026-05-01** — Policy created (issue
  [#63](https://github.com/robertocorreajr/cfs_spool/issues/63)).
