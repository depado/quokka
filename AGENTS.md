# Quokka — Agent Instructions

Quokka is a **boilerplate/template engine** CLI written in Go. It fetches a template (from git or local path), interactively prompts the user for variable values, and renders the template to an output directory.

## Build & Test

```sh
make build      # builds ./qk binary (CGO_ENABLED=0, injects version via ldflags)
make tmp        # builds to /tmp/qk for quick testing
make install    # go install to $GOPATH
make test       # go test ./...
go test ./...   # run all tests
```

The binary entry point is `./cmd/qk/main.go`. Version info is injected at build time via ldflags.

## Architecture

```
CLI (Cobra/Viper) → Renderer → {Provider, Conf}
cmd/              → renderer/ → {provider/, conf/}
```

| Package | Responsibility |
|---------|----------------|
| `cmd/` | Cobra CLI commands (`qk`, `new`, `version`), Viper flag bindings |
| `provider/` | Template source abstraction: git clone or local copy |
| `renderer/` | Entry point (`Render()`), two-phase directory walk with recursive include support (`Analyze()`) |
| `conf/` | Config structs, per-file rendering logic, variable/prompt system, includes |
| `utils/` | Colored output helpers, spinner wrappers |

## Key Conventions

**Template config files** are always named `.quokka.yml`. They define variables, delimiters, copy/ignore flags, and per-file metadata. The root config can also declare `includes` to recursively compose multiple templates.

**Config hierarchy** (highest-precedence last):
1. Root `.quokka.yml`
2. Subdirectory `.quokka.yml` files (apply to their subtree)
3. File-level frontmatter (`---` blocks at top of file)

**Frontmatter** is always stripped from output, even in copy-mode (no-render) files.

**Variable types** are determined by YAML structure in `variables:`:
- Default → text input prompt
- Has `values:` array → selection/dropdown
- Has `confirm: true` → yes/no confirmation
- Has nested `variables:` → sub-variables (only prompted if parent is truthy)

**Nested variable names** use `parent_child` in the template context (e.g., `slack.channel` → `{{.slack_channel}}`).

**Conditional rendering** on files uses `if:`:
- Single word → direct truthiness check on variable
- Expression string → evaluated by [expr-lang](https://github.com/expr-lang/expr)

**Includes** allow the root `.quokka.yml` to pull in and compose external templates recursively. An `Include` supports `source` (URL/path), `path` (inner sub-path within the fetched template, works for both git and local), `dest` (output sub-directory), `if` (condition), `confirm` (yes/no prompt), and `prompt` (custom confirm message). Included templates inherit the parent's accumulated context and can contribute variables back to the parent.

**Template functions** available: `title` (Title Case), `uc` (SCREAMING CASE).

**Default delimiters** are `{{ }}` but can be overridden per-directory or per-file via `delimiters:`.

## Important Files

- [`conf/conf.go`](conf/conf.go) — Core structs: `Config`, `ConfigFile`, `File`
- [`conf/root.go`](conf/root.go) — `Root` struct (embeds `ConfigFile`), `Include` struct for recursive template composition
- [`conf/file.go`](conf/file.go) — `File.Render()`: full per-file rendering pipeline
- [`conf/variables.go`](conf/variables.go) — Variable prompting logic, custom YAML parsing to preserve key order
- [`conf/input.go`](conf/input.go) — `InputCtx` type, `MergeCtx`, `GetInputContext`, `GetSetContext`
- [`conf/survey.go`](conf/survey.go) — Custom survey prompt templates (`init()`)
- [`renderer/analyze.go`](renderer/analyze.go) — Two-phase walk: collect variables across full template tree → render all files
- [`provider/provider.go`](provider/provider.go) — `Provider` interface (Fetch, Name, UsesTmp)
- [`cmd/flags.go`](cmd/flags.go) — Flag definitions

## CLI Flags & Config

All flags are bound to Viper with `QUOKKA_` env var prefix. See [`cmd/flags.go`](cmd/flags.go) for the full list. Notable flags: `--yes` (skip confirmations), `--git.depth` (clone depth, default 1), `--input` (pre-fill variables from YAML file), `--set key=value` (override individual variables).

## Testing

Tests live alongside source files. Use `testify` for assertions. See [`conf/variables_test.go`](conf/variables_test.go) and [`provider/provider_test.go`](provider/provider_test.go) and [`renderer/analyze_test.go`](renderer/analyze_test.go) for patterns.
