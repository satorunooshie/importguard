# ImportGuard
ImportGuard is a Go static analysis tool that helps to enforce import rules within your Go projects.

It allows you to specify allowed and denied import paths for specific packages using a configuration file.

This tool is particularly useful for maintaining codebase integrity by preventing unintended or unauthorized dependencies.

## Installation
This module targets Go 1.26 or later.
Projects analyzed by importguard may target older Go versions.

```bash
go install github.com/satorunooshie/importguard/v2/cmd/importguard@latest
```

## Configuration
A configuration file in JSON format is required, which specifies the allowed and/or denied import paths.

ImportGuard loads `.importguard.json` by searching upward from the analyzed Go files.

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

## Configuration Details
- Allow: Specifies non-standard import paths that are explicitly allowed for specific packages.
  - Note: Standard library packages do not need to be listed here; they are allowed by default. Only non-standard (external) packages that you explicitly want to allow should be listed.
- Deny: Specifies import paths that are denied for specific packages.
  - Note: This section can include both standard library packages and non-standard packages. Use this to list exceptions that should be denied, even if they are generally acceptable.

### Example
`repo/.importguard.json`

```json
{
 "allow": {
  "github.com/satorunooshie/repo/client": {
   "github.com/satorunooshie/repo/libs/collection": {}
  },
  "github.com/satorunooshie/repo/libs/crypto": {}
 },
 "deny": {
   "github.com/satorunooshie/repo/internal": {
     "fmt": {},
     "github.com/satorunooshie/repo/libs/collection": {}
   }
 }
}
```

`repo/libs/.importguard.json`

```json
{
 "allow": {
  "github.com/satorunooshie/repo/libs/crypto": {
   "github.com/satorunooshie/repo/libs/collection": {}
  }
 }
}
```

In this example, the following rules are applied:
- `github.com/satorunooshie/repo/client` can import `github.com/satorunooshie/repo/libs/collection` and standard library packages.
- `github.com/satorunooshie/repo/internal` cannot import `fmt` or `github.com/satorunooshie/repo/libs/collection`.
- `github.com/satorunooshie/repo/libs/crypto` matches both config files, so the nearer `repo/libs/.importguard.json` is used and `github.com/satorunooshie/repo/libs/collection` is allowed there.

## Example Output
If an import rule is violated, ImportGuard will output a message similar to the following:

```bash
~/importguard/testdata/src/github.com/satorunooshie/repo/internal/internal.go:4:2: prohibited import package: "fmt"
```

This indicates that the fmt package was imported in github.com/satorunooshie/repo/internal, which violates the rules defined in the configuration file.
