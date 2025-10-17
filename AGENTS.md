# Repository Guidelines

## Project Structure & Module Organization
The Go module lives at the repository root (`go.mod` targets Go 1.25.3). Add each pattern as its own package under `patterns/<pattern-name>`, keeping the package name short and all lowercase. Runnable examples belong in `cmd/<pattern-name>/main.go`, while shared helpers can sit in `internal/shared` and reusable fixtures in `testdata/`. Place brief notes alongside the code in `README.md` files so future readers can tie implementations back to the source material.

## Build, Test, and Development Commands
- `go test ./...` — run every unit test across pattern packages; use this before opening a pull request.
- `go run ./cmd/<pattern-name>` — execute an example or demo for the selected pattern.
- `go build ./...` — verify that new packages compile cleanly.
- `go fmt ./...` followed by `go vet ./...` — format the codebase and catch suspicious constructs; both must succeed before you push.

## Coding Style & Naming Conventions
Format code with `gofmt` (tabs for indentation). Keep exported identifiers in PascalCase, unexported ones in camelCase, and package names concise (e.g., `circuitbreaker`). Group related files by pattern and mirror Go's standard layout: implementation in `pattern.go`, tests in `pattern_test.go`, and example code in `example_test.go`. Avoid introducing third-party dependencies unless they reinforce the pattern being documented.

## Testing Guidelines
Use the Go `testing` package with table-driven tests named `Test<Pattern><Behavior>`. Document behavior in comments when it differs from the canonical pattern description. Place integration-style checks under `cmd/<pattern-name>/main_test.go` or in separate `_example_test.go` files. Aim for deterministic tests that run in under a second; log slow or flaky cases and gate merges on a clean `go test ./...`.

## Commit & Pull Request Guidelines
Commit subjects should be short, present-tense summaries (recent history mixes Japanese and English; follow whichever you prefer). Include a single pattern or improvement per change set and mention the associated chapter when relevant. Pull requests need: what changed, why it matters, how to validate (`go test ./...` output), and screenshots or logs if the pattern exposes a CLI. Rebase against `main`, fix merge conflicts locally, and ensure CI (when added) is green before requesting review.

## Communication Norms
Repository discussions, issues, and reviews must be in Japanese; include English only when referencing external sources.

## Pattern Implementation Checklist
Before submitting, confirm you have: a focused package under `patterns/`, runnable demo code when applicable, fmt/vet/test results, and documentation that references the pattern's source or expected use.
