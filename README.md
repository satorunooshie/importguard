# ImportGuard
ImportGuard is a Go static analysis tool that helps to enforce import rules within your Go projects.

It allows you to specify allowed and denied import paths for specific packages using a configuration file.

This tool is particularly useful for maintaining codebase integrity by preventing unintended or unauthorized dependencies.

## Installation
This module targets Go 1.26 or later.
Projects analyzed by importguard may target older Go versions.

```bash
go install github.com/satorunooshie/importguard/v4/cmd/importguard@latest
```

## Quick Start
A configuration file in JSON format is required.

- `deny` only: blacklist mode
- `allow` present: whitelist mode for non-standard imports
- `deny` always wins over `allow`

- Block a few imports and allow everything else: use `deny` only
- Allow only a few non-standard imports: add `allow`
- Allow everything except a few paths: use `allow: {"*": {}}` together with `deny`

Prefer exact matches when they are sufficient. Use `*` for "allow everything", and use `re^...` only when you need pattern matching. Prefer `*` over `re^.*$`.

- Blacklist example:

```json
{
 "deny": {
  "your/module/restricted": {
   "fmt": {},
   "github.com/satorunooshie/hoge": {}
  }
 }
}
```

- Whitelist example:

```json
{
 "allow": {
  "your/module/app": {
   "your/module/alloweddep": {}
  }
 }
}
```

- Allow everything, then deny one subtree:

```json
{
 "allow": {
  "your/module/usecase": {
   "*": {}
  }
 },
 "deny": {
  "your/module/usecase": {
   "re^your/module/forbidden(/.*)?$": {}
  }
 }
}
```

## Config Discovery
Place `.importguard.json` in your repository root, or in a subdirectory if you want rules that apply only to that subtree.

For example, if your repository looks like this:

```text
repo/
├── .importguard.json
└── service/
    ├── .importguard.json
    └── handler/
```

analyzing files in `repo/service/handler` will load `repo/service/.importguard.json`.

If multiple `.importguard.json` files are visible from the analyzed files, importguard uses the nearest one. Config files are not merged.

If no `.importguard.json` is found, importguard does not report any imports.

## Rule Reference
- `allow`
  Non-standard imports allowed for a package. If `allow` exists for a package, that package uses whitelist behavior for non-standard imports.
- `deny`
  Imports prohibited for a package. If a package has only `deny`, it behaves like a blacklist.
- `*`
  Matches any import path.
- `re^...`
  Treated as a regular expression. Use only when exact matches are not enough.

### Rule Evaluation

Import rules are defined per package path prefix. The keys under `allow` and `deny` select the package being checked, and the nested keys are the import paths matched against each import in that package.

- `deny` takes precedence over `allow`
- standard library packages are allowed by default
- packages with `allow` rules use whitelist behavior for non-standard imports
- packages with only `deny` rules use blacklist behavior
- package rules apply to the exact package and its subpackages
- when multiple package rules match, the longest package path wins

### Nearest Config Example
`example/.importguard.json`

```json
{
 "allow": {
  "github.com/satorunooshie/example/cases/localoverride": {
   "re^github\\.com/satorunooshie/example/targets/family(/.*)?$": {}
  }
 }
}
```

`example/cases/localoverride/.importguard.json`

```json
{
 "allow": {
  "github.com/satorunooshie/example/cases/localoverride": {
   "github.com/satorunooshie/example/targets/exact": {}
  }
 }
}
```

When analyzing `github.com/satorunooshie/example/cases/localoverride`, the nearer config in `cases/localoverride/.importguard.json` is used instead of the root config.
