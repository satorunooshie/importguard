# ImportGuard
ImportGuard is a Go static analysis tool that helps to enforce import rules within your Go projects.

It allows you to specify allowed and denied import paths for specific packages using a configuration file.

This tool is particularly useful for maintaining codebase integrity by preventing unintended or unauthorized dependencies.

## Installation
```bash
go install github.com/satorunooshie/importguard/cmd/importguard@latest
```

## Configuration
A configuration file in JSON format is required, which specifies the allowed and/or denied import paths.

The configuration file is loaded through an environment variable named `IMPORTGUARD_CONFIG`.

```bash
export IMPORTGUARD_CONFIG=/path/to/config.json
```

## Configuration Details
- Allow: Specifies non-standard import paths that are explicitly allowed for specific packages.
  - Note: Standard library packages do not need to be listed here; they are allowed by default. Only non-standard (external) packages that you explicitly want to allow should be listed.
- Deny: Specifies import paths that are denied for specific packages.
  - Note: This section can include both standard library packages and non-standard packages. Use this to list exceptions that should be denied, even if they are generally acceptable.

### Example
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

In this example, the following rules are applied:
- Allow
  - The `github.com/satorunooshie/repo/client` package is allowed to import only the `github.com/satorunooshie/repo/libs/collection` package and std packages.
  - The `github.com/satorunooshie/repo/crypto` package is allowed to import only std packages.
- Deny
  - The `github.com/satorunooshie/repo/internal` package is denied from importing the `fmt` package and the `github.com/satorunooshie/repo/libs/collection` package, and everything else is allowed.

## Example Output
If an import rule is violated, ImportGuard will output a message similar to the following:

```bash
~/importguard/testdata/src/github.com/satorunooshie/repo/internal/internal.go:4:2: prohibited import package: "fmt"
```

This indicates that the fmt package was imported in github.com/satorunooshie/repo/internal, which violates the rules defined in the configuration file.
