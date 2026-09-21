# Architecture — Add a Go Hello World program

`changes_existing_design: true` — not because the change is large, but because there is no
architecture here to change. The repository currently contains no software: `git ls-files`
lists only `.gitignore`, `README.md` and `docs/REQUIREMENTS-add-a-go-hello-world-program.md`;
there is no `*.go`, no `go.mod`, no `cmd/`, no CI config, no Makefile and no `.claude/`.
This stage therefore establishes the repository's first Go module, its first executable
package and its first executable test contract. The design below is deliberately three
files, because FR-5 and AC-07 of this job's requirements forbid anything more.

## 1. Architecture summary

A single-binary, single-package, no-service design:

```
repository root  (= Go module root)
├── go.mod          module declaration, no dependencies
├── main.go         package main; the only entry point; writes the greeting
└── main_test.go    package main; builds the binary and asserts its exact bytes
```

There is no server, no queue, no database, no configuration file, no background process and
no persistent state. The whole run is one short-lived process that writes thirteen bytes to
file descriptor 1 and exits. Every element above is required by a named criterion; nothing
else is built. Specifically, the design refuses an HTTP/CLI surface (this job's §6 Out of
scope, NFR-3) and refuses any third-party dependency (FR-3, NFR-2).

#### Decisions taken against the requirements' open questions

The requirements document could not inspect the repository, so it left Q1–Q5 open. Resolved
against the checked-out tree:

- **Q1 (existing Go code/module):** none exists, so a fresh module at the repository root is
  the correct choice, and the path is the one FR-3 names: `github.com/aliceyao1/mytestproject`.
- **Q2 (existing program layout):** there is no `cmd/` tree, no `internal/` tree and no other
  convention to follow, so the entry point goes at the module root as `main.go` (FR-4). This
  is the first program in the repository, not a deviation from one.
- **Q3 (pinned greeting):** neither `README.md` (15 bytes, the title line only) nor any code
  establishes a greeting, so the default is pinned: `Hello, World!`.
- **Q4 (stack disagreement):** the tree holds no source in another language — the only legacy
  artefact is a Python-flavoured `.gitignore` from a repository template. There is no
  conflict to report and nothing to migrate.
- **Q5 (README build instructions):** not added. AC-07 forbids touching files outside the
  three, and the requirements class this as a human decision.

One assumption in the requirements document is contradicted by the tree and must not be
acted on: **A4** asserts `main` is the checked-out delivery branch. The working branch is
`flowai/add-a-go-hello-world-program-20260921-115857-d6ebf251`, which matches `AGENTS.md`.
Implementation runs on the current branch and creates, switches or rewrites no branch.

## 2. Primary interfaces

This program exposes no network API. Its interfaces are the module declaration, the process
contract and the test harness; these are the complete set of seams a reviewer or a
validation gate can exercise.

### I-1 Module / build interface

`go.mod` at the module root, exactly:

```
module github.com/aliceyao1/mytestproject

go 1.22
```

No `require`, no `toolchain`, no `replace`, no `vendor/`. `go 1.22` is chosen because the
provisioned toolchain is go1.22.5 (measured: `go version` returns `go version go1.22.5
darwin/arm64`), satisfying AC-04's "no newer than the installed toolchain" while staying
portable across 1.22.x patch releases. `go mod init` would write `go 1.22.5`; that also
satisfies AC-04, but the pinned `1.22` is the recorded choice so the file is reproducible by
hand.

### I-2 Process interface (the external contract of the program)

| Direction | Channel | Contract |
| --- | --- | --- |
| input | argv | ignored; no arguments or flags are parsed (FR-1) |
| input | stdin | not read (FR-1) |
| output | stdout (fd 1) | exactly the bytes `Hello, World!` followed by one newline; nothing else |
| output | stderr (fd 2) | empty |
| output | exit status | `0` |

Errors: the contract defines no error channel, because no input can make the program fail.
The one edge case is a closed downstream reader (`... | head -0`): Go's runtime raises
`SIGPIPE` when a write to fd 1 or fd 2 fails with `EPIPE`, so the process can die by signal
in that case. That is outside the acceptance surface (AC-03 covers a normal run), and it must
not be "fixed" by adding stderr diagnostics, which would break the empty-stderr half of the
contract.

### I-3 Source-level surface

`main.go`, the only executable source file:

```go
package main

import "fmt"

const message = "Hello, World!"

func main() {
	fmt.Println(message)
}
```

The package is `main`, the entry point is `func main()`, and the greeting is an unexported
constant so the printed value has exactly one definition in non-test code. No exported
symbol is introduced, so no other package may become a consumer of this one (AC-01).

### I-4 Test interface

`main_test.go` is in the same package and provides one test; its behaviour is fixed by this
design rather than left to the implementer:

1. An `exec.LookPath("go")` guard. If absent, skip with a message naming the missing
   toolchain; if present, build with `exec.Command("go", "build", "-o", bin, ".")` where
   `bin = filepath.Join(t.TempDir(), "hello")`. Building into a temp directory is a
   requirement, not a preference — see the artefact constraint in §4.
2. Run the built binary with the child's working directory set to a directory other than the
   module root, so working-directory independence (AC-03, AC-08) is exercised.
3. Capture stdout and stderr into separate buffers and read the exit status from
   `cmd.ProcessState.ExitCode()`.
4. Assert, in this order: exit status `0`; stderr empty; `string(stdout) == "Hello, World!\n"`.

The expected string is written as a **literal in the test**, never as a reference to the
`message` constant. A test comparing the constant with itself would still pass after the
greeting was mutated, which the requirements explicitly reject.

### I-5 Verification interface

Command → observable result, all run from the module root:

| Command | Expected | Artefact left behind |
| --- | --- | --- |
| `go build ./...` | exit 0, no output | **`./mytestproject` binary in the module root** |
| `go build -o "$(mktemp -d)/hello" .` | exit 0 | none |
| `go vet ./...` | exit 0, no output | none |
| `gofmt -l .` | no output | none |
| `go test ./...` | `ok ...`, exit 0 | none (the test builds into `t.TempDir()`) |
| `go run .` | prints `Hello, World!`, exit 0 | none |

## 3. Data model and state

There is none, and that is a decision rather than an omission.

| Concept | Representation | Storage | Lifecycle |
| --- | --- | --- | --- |
| greeting | 13-byte UTF-8 string constant | compiled into the binary | process lifetime |
| expected output | byte literal in the test | test binary | test lifetime |
| emitted bytes | 14 bytes on fd 1 (13 + one newline) | none | written once, never retained |

Relationship: the test literal and the production constant must denote the same bytes, and
the test is the only thing enforcing that. No file, socket, environment variable, clock,
locale or random source is read, so output is byte-identical on every run and on any host
(NFR-1).

## 4. Operational constraints

**Sandbox and toolchain.** Sandbox `local-agent:d6ebf251-...` runs in `local_agent` mode with
host network access, dependency installation allowed and system package installation
forbidden. The design needs none of that latitude: it is standard-library-only, so no module
download, no vendored tree and no cgo occur. Verified with `GOPROXY=off`: the module builds
and the subprocess test passes with zero network access, and `GOTOOLCHAIN=auto` will not try
to fetch a toolchain because the `go` directive is not newer than the installed one. An
unbounded time budget is irrelevant to a process whose runtime is milliseconds (NFR-3).

**Build artefacts — the one operational trap.** `go build ./...` in a module whose only
package is `main` writes an executable named after the last element of the module path
(`mytestproject`) into the current directory. Measured here on go1.22.5: the file appears in
the module root after `go build ./...` and the command exits 0. Nothing in this repository
ignores it — `git check-ignore -v mytestproject` exits 1, and the `.gitignore` is the Python
template with no Go section. Because AC-07 audits `git status --short`, the implementer must
remove that one artefact path before the diff audit, or use the artefact-free form
`go build -o "$(mktemp -d)/hello" .` alongside it. `go run .`, `go test ./...` and
`go vet ./...` leave nothing behind (all three measured). Adding a Go section to `.gitignore`
is not a permitted fix, because that would be a fourth file.

**Repository boundary.** The change set is exactly the three files in §1. No existing file is
edited, renamed or reformatted; `README.md`, `AGENTS.md`, `CLAUDE.md` and every file under
`docs/` are untouched. `AGENTS.md` and `CLAUDE.md` currently show as untracked in
`git status --short`; they are FlowAI-owned, they are not part of this change set, and the
implementer must neither stage nor edit nor delete them.

**Authentication and data isolation.** Not applicable, and the architecture must not invent
it. The program reads no data, holds no data, crosses no tenant boundary and calls no
authorisation path, so there is no isolation surface to preserve (NFR-5). The trigger for
revisiting is any future requirement to read configuration, files, environment secrets or
remote resources, or to serve more than one caller; under this design that is a new job, not
an extension of this one.

**Failure handling and recoverability.** There is no runtime state to corrupt and no
partial-completion window: the process either writes its fourteen bytes and exits 0, or it does
not start. Failure routes are: a compile error fails `go build` (fix the source, nothing to
roll back); a mutated greeting fails the test by exact-byte comparison, which is the
detection path AC-05 relies on; a missing toolchain makes the test skip with a recorded
reason, which the implementer must report rather than leave silent. Recovery from the whole
change is deleting the three added files, or reverting the commit FlowAI creates for it.

## 5. Verification plan and stage ownership

| Criterion | Check | Stage |
| --- | --- | --- |
| AC-01 | `main.go` declares `package main` with `func main()` | implementation, then peer review |
| AC-02 | `go build ./...` and `go vet ./...` exit 0 | implementation, then SDET / quality gate |
| AC-03 | `go run .` prints the pinned bytes, empty stderr, exit 0 | implementation, then SDET |
| AC-04 | `go.mod` contents and absence of any `require` | implementation, then structural gate |
| AC-05 | `go test ./...` green with the exact-byte subprocess assertion | implementation, then SDET |
| AC-06 | `gofmt -l .` prints nothing | implementation, then quality gate |
| AC-07 | `git status --short` plus diff review show only the three files, artefact removed | structural gate |
| AC-08 | the built binary run from a temp working directory gives identical bytes | SDET |

## 6. Relationship to the existing codebase

There is no prior application code to extend: this design is additive and touches nothing
that exists. The two conventions it *does* inherit are the language declared in `AGENTS.md`
("Language: golang") and the module path named in FR-3. Where the workspace offers a
convention this design does not adopt, it is because adopting it would break a criterion —
notably the `.gitignore`: adding Go ignore rules would be a fourth file, so the artefact
problem is solved by the verification commands instead.

Anticipated growth, deliberately not built: a second program would move to `cmd/<name>/` with
the greeting in an importable package, and the module would gain a dependency only when a
requirement names one. Nothing in this job justifies either step now.

## 7. Design record consulted

Of the design-record documents this job names, only the requirements document exists so far:
`docs/REQUIREMENTS-add-a-go-hello-world-program.md` (read in full). It is binding, and it set
the change to three files, the module path, the root-level `main.go` default, the pinned
greeting, the exact-byte assertion and the AC-01…AC-08 list, all carried into §1–§5 above.
The same document's answers were incomplete in two places, and this document supplies them:
its assumptions A1–A4 were unverified, so they are resolved here against the tree (Q1–Q5 in
§1, with A4 corrected because the checked-out branch is not `main`), and its AC-02 names
`go build ./...` without noting the binary that command leaves behind, so §4 records that
artefact and the audit consequence. No implementation, peer-review, validation or delivery
document for this job existed to read — `docs/` held the requirements file and one
`flowai-task-deliverables/` envelope from the product-manager stage, whose `manifest.json`
and `task-output.json` confirmed the requirements text was the only upstream artefact.

## 8. Workspace analysis performed

Read: `README.md` (one title line, no build instructions); `AGENTS.md` and `CLAUDE.md` (byte
identical; they declare the Go stack, the FlowAI branch and the design-record list, and they
forbid committing from an implementation stage); `.gitignore` (the GitHub Python template,
no Go section); `docs/REQUIREMENTS-add-a-go-hello-world-program.md`;
`docs/flowai-task-deliverables/product_manager-Product-Requirements-fcbb9a85-bd6/`
(`manifest.json`, `task-output.json`). Ran: `ls -la`, `git log --oneline -10`,
`git status --short`, `git ls-files`, `git branch -a`, `find` for `*.go` and `go.mod` (none),
`git check-ignore -v` for the build-output name (not ignored), `go version` (go1.22.5), and
scratch builds in `mktemp -d` outside the workspace to measure what `go build ./...`,
`go build -o`, `go run .`, `go test ./...` and `GOPROXY=off` actually do. `README.md` is the
repository's own account of itself and is deliberately left unchanged (Q5); no `.claude/`,
`adr/`, `rfc/` or `design/` directory exists, and no project manifest other than the
Python-flavoured `.gitignore` was found.
