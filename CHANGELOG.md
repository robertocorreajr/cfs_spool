# Changelog

All notable changes to CFS Spool will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security

## [3.1.1] - 2026-05-27

### Added

- Full Creality Print v7 material catalog: 18 new entries covering Generic PA-GF/PP-GF, eight eSUN materials (including eSUN PLA-Basic, PLA+HS, PLA-LW, PETG-CF, ABS-CF, ABS+HS, TPU-95A), and nine Polymaker Fiberon engineering filaments.
- Developer CLI `cmd/check-materials` to diff the bundled material list against the Creality Print v7 `materialList.json` and flag missing, stale, or mismatched entries.

### Changed

- Material dropdown is now sorted alphabetically by name across all vendors.
- Corrected material code `00035`: previously mislabeled as "eSUN PLA-LW", now correctly identified as "Generic PLA-LW" (vendor: Generic). The real eSUN PLA-LW is now code `E1008`. Tags written by previous versions with code `00035` will display "Generic PLA-LW" after this update.

## [3.1.0] - 2026-05-01

### Added

- In-app notification when a new version of CFS Spool is published. The header shows a pulsing icon when an update is available; clicking it opens a modal with the changelog and a download button that picks the right asset for your operating system (macOS, Linux, or Windows).
- "Ignore this version" option that persists locally so the same release stops nagging you on every startup.

### Changed

- The startup update check runs in a background goroutine with a 3 s timeout, so a slow or offline network never delays the app launch.

## [3.0.10] - 2026-05-01

### Fixed

- Custom filament length values were being silently ignored due to a short-circuit in `convertLength` when the lookup table didn't match. The function now always honors decimal input and falls back to a calculated value, fixing custom-grams writes for non-default sizes.

### Changed

- CI workflows updated to Node.js 24.

## [3.0.9] - 2026-05-01

### Changed

- `Tests` workflow split out from the `Auto Tag` workflow so test runs are no longer skipped when a commit is marked with `[skip release]`. Test coverage now runs on every push, regardless of whether a release is being cut.

### Fixed

- README documentation for the development hooks updated to match the new layout.

## [3.0.8] - 2026-05-01

### Changed

- `scard.Context` lifecycle for Open/Close paths is now created through a mockable factory, simplifying tests that need to simulate PC/SC errors without hardware.

[Unreleased]: https://github.com/robertocorreajr/cfs_spool/compare/v3.1.1...HEAD
[3.1.1]: https://github.com/robertocorreajr/cfs_spool/compare/v3.1.0...v3.1.1
[3.1.0]: https://github.com/robertocorreajr/cfs_spool/compare/v3.0.10...v3.1.0
[3.0.10]: https://github.com/robertocorreajr/cfs_spool/compare/v3.0.9...v3.0.10
[3.0.9]: https://github.com/robertocorreajr/cfs_spool/compare/v3.0.8...v3.0.9
[3.0.8]: https://github.com/robertocorreajr/cfs_spool/compare/v3.0.7...v3.0.8
