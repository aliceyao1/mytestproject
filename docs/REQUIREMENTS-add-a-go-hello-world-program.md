# Product Requirements — Go Hello World in aliceyao1/mytestproject

## 1. Requirements summary

Add the smallest possible Go Hello World program to the existing repository `aliceyao1/mytestproject` on branch `main` (working tree `/Users/yaoalice/Documents/flowai-projects`), together with the minimal module configuration needed to build and run it locally, and one automated test that proves the program prints its message. The complete change is at most three added files: a Go entry point, `go.mod`, and a Go test. Everything else — servers, CLI arguments, containers, pipelines, dependencies, documentation rewrites, refactoring of existing code — is out of scope. The four customer-specified acceptance criteria are carried through unchanged in intent and restated as executable checks.

## 2. Evidence base, assumptions and open questions

**Evidence base.** This document was authored by the product-manager stage, which has no shell, filesystem or network access. It is grounded in the goal statement (repository, branch, checkout path, Go stack, README.md as existing guidance) and the customer acceptance criteria. **No repository inspection was performed by this stage.** The implementer must run workspace analysis (`ls -la`, `git log`, `git status --short`, `README.md`, `.flowai/CODEBASE_CONTEXT.md` if present, any existing `go.mod` or `.go` files) and confirm or correct the assumptions below before writing code.

**Assumptions (must be confirmed):**
- A1. The repository is cloned at the stated path with `main` checked out; working tree is clean before the change.
- A2. The repository does not already contain a `go.mod` or Go source. If it does, reuse the existing module path, `go` directive and directory convention instead of initialising a second module.
- A3. A Go toolchain capable of building a `package main` (Go 1.21 or later) is available in the sandbox and can build offline.
- A4. `main` is the delivery branch; no new branch or commit is created by the implementer.

**Open questions (resolve during workspace analysis, record the answer in the delivery report):**
- Q1. Does the tree already contain Go code or a `go.mod`? If yes, which module path, `go` directive and layout apply?
- Q2. Does the repository already have an established convention for where new programs live (for example a `cmd/` tree)? If so, follow it and record the deviation from a root-level `main.go`; do not create a new layout when one exists.
- Q3. Does `README.md`, or any existing program, already establish a specific Hello World string? If so, pin that string in the test. Otherwise use `Hello, World!`.
- Q4. The stack note says the working tree is authoritative where it disagrees with the stated language. If the tree is an unrelated-language project, that does not remove the Go deliverables: AC-02 and AC-03 are customer-specified and Go-specific. Record the discrepancy rather than substituting another language's program.
- Q5. Should `README.md` gain a one-line build/run instruction? Not a requirement here, because of the no-unrelated-changes criterion. Treat as a human decision, not an implementer decision.

## 3. Functional requirements

- **FR-1 Entry point.** Add a Go program with `package main` and `func main()` that writes the Hello World message to standard output. It takes no arguments and reads no input.
- **FR-2 Output contract.** On success the program writes exactly the pinned message bytes followed by a single trailing newline, writes nothing else to stdout, writes nothing to stderr, and exits with status 0. Default pinned message: `Hello, World!`. If Q3 resolves to a different established string, pin that string instead and record the choice; exactly one string is pinned and asserted.
- **FR-3 Module configuration.** Add `go.mod` at the module root with module path `github.com/aliceyao1/mytestproject` (or the existing path per Q1) and a `go` directive no newer than the installed toolchain. No `require` entries; standard library only.
- **FR-4 Layout.** Place the entry point at the module root as `main.go` unless an established layout exists (Q2), in which case follow it. Test file lives beside the entry point in the same package or the convention the repo already uses.
- **FR-5 Minimal, isolated change set.** Only files required for the entry point, module file, and test are added. No existing file is edited, deleted, renamed, reformatted or re-dependency-upgraded.
- **FR-6 Automated coverage.** Add at least one Go test that fails if the printed output differs from the pinned string, comparing exact bytes rather than a substring. The preferred test executes the built program as a subprocess and compares captured stdout and exit code (this is the end-to-end check for this feature). An additional pure-function unit test is acceptable but optional.
- **FR-7 Verification commands.** The delivery report records the exact commands run and their observed results: `go build ./...`, `go vet ./...`, `gofmt -l .`, `go test ./...`, and `go run .` from the module root.
- **FR-8 Delivery report.** The implementer's report states the Go version used, whether Go files or a module already existed, the layout chosen and why, the pinned greeting string, the diff scope, and any open question left unresolved.

## 4. Acceptance criteria

- **AC-01** The repository contains a Go program entry point: a file declaring `package main` with a `func main()`, at a location consistent with the repository's layout. *(Customer criterion 1)*
- **AC-02** `go build ./...` run from the module root exits 0 with no compilation errors, and `go vet ./...` exits 0. *(Customer criterion 1)*
- **AC-03** Running the program (`go run .` from the module root, or the compiled binary executed from any working directory) writes the pinned Hello World message followed by a single newline to stdout, nothing to stderr, and exits with status 0. *(Customer criterion 2)*
- **AC-04** `go.mod` exists at the module root with a module path and a `go` directive no newer than the installed toolchain, and declares no third-party dependencies. *(Customer criterion 3)*
- **AC-05** `go test ./...` exits 0 and includes at least one test asserting the exact stdout bytes of the built program. *(FR-6, test focus)*
- **AC-06** `gofmt -l .` lists no files, i.e. all added Go source is gofmt-clean.
- **AC-07** The change set, audited with `git status --short` and a diff review, contains only the entry point, the module file and the test file; no pre-existing file is modified, deleted or renamed, and any pre-existing test suite in the repository still passes. *(Customer criterion 4)*
- **AC-08** The program has no side effects beyond writing to stdout: it opens no files or sockets, reads no environment secrets, starts no background process and keeps no state.

## 5. Non-functional requirements

- **NFR-1 Determinism.** Output is byte-identical on every run, independent of locale, timezone, hostname and working directory. No timestamps, randomness or environment-dependent formatting.
- **NFR-2 Sandbox executability.** The program builds and runs offline with the provisioned Go toolchain, using only the standard library: no module downloads, no vendored dependencies, no cgo, no network access at build or run time.
- **NFR-3 Footprint.** A single short-lived process: no listening ports, no daemons, no persistent state, observable runtime well under a second.
- **NFR-4 Source hygiene.** Added Go source is gofmt-clean and passes `go vet`; no suppressions or generated artefacts are introduced.
- **NFR-5 Scope boundary on security surfaces.** The change must introduce no authentication, authorization, tenant-scoped data access or persistence path. The generic tenant-isolation and authorization NFRs from the task template do not apply to a program that touches no data — but if any implementation approach would require them, stop and raise an open question rather than adding them.

## 6. Out of scope

- HTTP/REST/gRPC servers, web frameworks, routing, middleware.
- CLI argument or flag parsing, environment configuration, config files, dotenv handling.
- Dockerfiles, container images, Kubernetes manifests, CI/CD pipelines, Makefiles, linter or formatter configuration files.
- Any third-party Go dependency, vendoring, or module proxy configuration.
- Creating new package trees, `internal/` layouts, or a `cmd/` structure where none exists.
- Internationalisation, custom logging frameworks, coloured or formatted terminal output.
- Benchmarks, load, performance or chaos testing.
- README or documentation rewrites, LICENSE changes, changelog entries.
- Refactoring, reformatting, dependency upgrades, or language migration of any existing code.
- Release artefacts: versioning, installers, packaging, cross-compilation target matrices.

## 7. Test focus areas

- **AC traceability matrix.** Every AC-01…AC-08 maps to at least one automated or scripted check; the delivery report enumerates the mapping.
- **Exact-output assertion (primary).** End-to-end test that builds and executes the program, capturing stdout, stderr and exit code, asserting exact equality against the pinned string plus newline. A `contains` check is not acceptable: a mutation of the message must fail the test.
- **Build and static checks.** `go build ./...`, `go vet ./...`, `gofmt -l .` each exit clean; results recorded.
- **Working-directory independence.** Execute the compiled binary from a directory other than the module root; output must be identical (AC-08 edge case).
- **Negative/scope regression.** Audit the diff: no pre-existing file touched; run any test suite the repository already ships and confirm it still passes; confirm no new dependency appears in `go.mod`.
- **Toolchain precondition.** Record `go version` in the report so a failure to build can be distinguished from a missing toolchain.

## 8. Risks

- **R1 Root-level package collision.** If the module root already contains a `package main` or another `main.go`, adding a second entry point breaks the build. Resolve via Q2 and place the program in its own directory rather than merging packages.
- **R2 Subprocess test flakiness.** A test that shells out to `go build`/`go run` depends on the toolchain being present in the test environment. Mitigate by building once into `t.TempDir()` and executing the resulting binary, and skip with a clear message only if the toolchain is genuinely absent (recorded, not silently suppressed).
- **R3 Scope creep through the README.** The repository guidance doc invites documentation updates; AC-07 forbids unrelated edits. Keep documentation changes out unless a human asks (Q5).
- **R4 Unverified pre-existing state.** Because this stage could not inspect the repository, A1-A4 may be wrong; the implementer must correct them in the delivery report rather than coding around a stale assumption.