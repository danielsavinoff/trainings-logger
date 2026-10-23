# GRACE 4 Project Engineering Protocol

## Keywords
training, workout logging, exercises, sets, Go, terminal UI, Bubble Tea, SQLite

## Annotation
A local terminal application for recording exercises and workout sets, including repetitions, weight, rest intervals, and workout history, backed by SQLite.

## GRACE 4 Source of Truth

This project uses the GRACE 4 `.grace` artifact model.

- Product and technical context: `.grace/context/*.xml`
- Current graph projection source: `.grace/graph/index.xml` plus routed graph documents such as `.grace/graph/main.xml`
- Current verification projection source: `.grace/verification/index.xml` plus routed verification documents such as `.grace/verification/main.xml`
- Active work: `.grace/changes/active/C-*/spec.xml` and `.grace/changes/active/C-*/plan.xml`
- Completed or terminal work: `.grace/changes/archive/C-*/*`

Legacy `docs/*.xml` files are not GRACE 4 state. If legacy GRACE 3 docs appear, use `grace-migrate`; do not silently validate, convert, or delete them.

## Workflow Rules

1. Do not implement source behavior before an approved active `GraceChangeSpec` and `GraceChangePlan` exist, unless the user explicitly requests a small direct fix.
2. Treat `spec.xml` as normative. Treat `design-context.xml` as explanatory memory only.
3. Before execution, check `BaselineAssertions`, `TargetAssertions`, `DurableScope`, and `ObservedWriteScope` in the plan.
4. Update durable `.grace` graph and verification state only as part of the approved change lifecycle.
5. Never store transient run state by mutating approved XML statuses. Runtime states are derived from current files, assertions, and scopes.

## Semantic Anchor Rules

- GRACE semantic anchors are XML tags, never attributes: use `<M-EXAMPLE />`, not `<Module ref="M-EXAMPLE" />`.
- Module IDs use `M-*`; data-flow IDs use `DF-*`; graph document wrappers use `GD-*`; verification entries use deterministic `V-M-*`; verification document wrappers use `VD-*`; change bundles use `C-*`.
- Code-level semantic markup remains grep-stable: `START_MODULE_CONTRACT`, `START_MODULE_MAP`, `START_CONTRACT:`, `START_BLOCK_`, and `START_CHANGE_SUMMARY`.

## Grep-First Navigation

1. Locate module ownership through `.grace/graph/index.xml`, then open the routed graph document.
2. Locate verification through `.grace/verification/index.xml`, then open the routed verification document.
3. Locate active work through `.grace/changes/active/C-*`.
4. Use file-local `LINKS:` fields and `START_BLOCK_` anchors to narrow code reads before loading whole files.

## CLI Checks

- `grace lint --path .` validates `.grace` grammar, projections, assertions, lifecycle locations, and scope overlaps.
- `grace status --path .` summarizes durable and operational GRACE 4 health.
- `grace module`, `grace verification`, and `grace file` navigate graph, verification, and file-local anchors.

## File-Local Markup Reference

```ts
// START_MODULE_CONTRACT
//   PURPOSE: [What this module does]
//   SCOPE: [Bounded responsibility]
//   DEPENDS: [M-* dependencies or none]
//   LINKS: [Related M-* and V-M-* anchors]
// END_MODULE_CONTRACT
//
// START_MODULE_MAP
//   exportedSymbol - one-line responsibility
// END_MODULE_MAP
//
// START_CONTRACT: functionName
//   PURPOSE: [What it does]
//   INPUTS: { paramName: Type - description }
//   OUTPUTS: { ReturnType - description }
//   SIDE_EFFECTS: [External state changes or none]
// END_CONTRACT: functionName
//
// START_BLOCK_EXAMPLE
// ... implementation slice ...
// END_BLOCK_EXAMPLE
```
## Repository Context

- Go module: `github.com/danielsavinoff/trainings-logger`; `go.mod` declares Go 1.27.
- `cmd/cli`: process entry point and dependency wiring.
- `internal/domain`: exercises, workout sets, weight values and units, user preferences, and repository interfaces.
- `internal/app`: command and query handlers plus the transaction-runner port.
- `internal/infra/persistence`: SQL transaction support and SQLite repositories and database setup.
- `internal/infra/tui`: application bindings, startup model, record views, reusable components, and TUI tests.
- `migrations`: embedded SQLite schema migrations, applied when the database opens.

These paths describe the observed package layout. The initial GRACE graph and verification documents are empty bootstrap skeletons; they do not yet claim module ownership, file-local markup coverage, or verification evidence. Populate them through subsequent approved changes.

## Project Checks

- `go test ./...`: existing automated tests, including TUI snapshots and interaction tests.
- `go vet ./...`: Go static checks.
- Do not set `UPDATE_GOLDEN=1` during normal validation; review snapshot updates when changing intended rendering.
- Run the application with `go run ./cmd/cli` only when interactive execution is needed; it opens the real user database and applies embedded migrations.

## Go Source Formatting

- In every `.go` source file, including tests, write calls to functions or methods with more than one argument using one argument per line, with a trailing comma after the final argument.
- Separate distinct logical blocks of Go code with a blank line. Keep statements that form one continuous operation together.

## Use Cases

- Start every use case name with a verb, such as `CreateExercise`, `ListSets`, or `SaveRemoteState`.
- Keep exactly one use case operation in each source file. Split files containing multiple operations into separate files named after their use cases.

## Domain Model Construction

- Initialize every domain model through a constructor function whose name starts with `New` and validates all model invariants before returning the model.
- Keep each domain model in its own source file with its constructor in that file; keep each enum in its own source file.
- Do not add `Valid` methods as a substitute for constructor validation. Call the constructor at persistence and application boundaries instead of assembling domain model values directly.
- Define domain errors only in package's errors.go file.

## Keyboard Shortcuts

- When adding keyboard shortcuts, prefer familiar Vim-style bindings where they fit the existing interaction model.
