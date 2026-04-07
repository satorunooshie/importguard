package importguard

import (
	"encoding/json"
	"errors"
	"go/ast"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type config struct {
	Allow map[string]map[string]struct{} `json:"allow"`
	Deny  map[string]map[string]struct{} `json:"deny"`
}

type compiledConfig struct {
	Allow map[string]matcherSet
	Deny  map[string]matcherSet
}

type matcherSet struct {
	Exact       map[string]struct{}
	HasWildcard bool
	Regex       []*regexp.Regexp
}

const configFileName = ".importguard.json"
const regexPrefix = "re^"

var Analyzer = &analysis.Analyzer{
	Name: "importguard",
	Doc:  "importguard reports prohibited imports",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func parseConfig(pass *analysis.Pass) (compiledConfig, error) {
	fp, err := resolveConfigPath(pass)
	if err != nil {
		return compiledConfig{}, err
	}
	if fp == "" {
		return compiledConfig{}, nil
	}
	return loadConfig(fp)
}

func loadConfig(fp string) (compiledConfig, error) {
	b, err := os.ReadFile(fp)
	if err != nil {
		return compiledConfig{}, err
	}
	var raw config
	err = json.Unmarshal(b, &raw)
	if err != nil {
		return compiledConfig{}, err
	}
	return compileConfig(raw)
}

func compileConfig(raw config) (compiledConfig, error) {
	allow, err := compileMatcherMap(raw.Allow)
	if err != nil {
		return compiledConfig{}, err
	}
	deny, err := compileMatcherMap(raw.Deny)
	if err != nil {
		return compiledConfig{}, err
	}
	return compiledConfig{
		Allow: allow,
		Deny:  deny,
	}, nil
}

func compileMatcherMap(raw map[string]map[string]struct{}) (map[string]matcherSet, error) {
	compiled := make(map[string]matcherSet, len(raw))
	for pkgPath, patterns := range raw {
		ms, err := compileMatcherSet(patterns)
		if err != nil {
			return nil, err
		}
		compiled[pkgPath] = ms
	}
	return compiled, nil
}

func compileMatcherSet(patterns map[string]struct{}) (matcherSet, error) {
	ms := matcherSet{
		Exact: make(map[string]struct{}, len(patterns)),
	}
	for pattern := range patterns {
		switch {
		case pattern == "*":
			ms.HasWildcard = true
		case strings.HasPrefix(pattern, regexPrefix):
			re, err := regexp.Compile("^" + strings.TrimPrefix(pattern, regexPrefix))
			if err != nil {
				return matcherSet{}, err
			}
			ms.Regex = append(ms.Regex, re)
		default:
			ms.Exact[pattern] = struct{}{}
		}
	}
	return ms, nil
}

func resolveConfigPath(pass *analysis.Pass) (string, error) {
	var filenames []string
	for _, f := range pass.Files {
		filename := pass.Fset.PositionFor(f.Pos(), false).Filename
		if filename == "" {
			continue
		}
		filenames = append(filenames, filename)
	}
	return findConfigFile(filenames...)
}

func findConfigFile(filenames ...string) (string, error) {
	type candidate struct {
		path  string
		depth int
	}

	var candidates []candidate
	seen := map[string]struct{}{}

	for _, filename := range filenames {
		if filename == "" {
			continue
		}

		dir, err := filepath.Abs(filepath.Dir(filename))
		if err != nil {
			return "", err
		}

		for depth := 0; ; depth++ {
			fp := filepath.Join(dir, configFileName)
			if _, err := os.Stat(fp); err == nil {
				if _, exists := seen[fp]; !exists {
					seen[fp] = struct{}{}
					candidates = append(candidates, candidate{path: fp, depth: depth})
				}
			} else if !errors.Is(err, os.ErrNotExist) {
				return "", err
			}

			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	if len(candidates) == 0 {
		return "", nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].depth != candidates[j].depth {
			return candidates[i].depth < candidates[j].depth
		}
		return candidates[i].path < candidates[j].path
	})
	return candidates[0].path, nil
}

func run(pass *analysis.Pass) (any, error) {
	conf, err := parseConfig(pass)
	if err != nil {
		return nil, err
	}

	allowlist, aTarget := conf.Allow[pass.Pkg.Path()]
	denylist, dTarget := conf.Deny[pass.Pkg.Path()]
	if !aTarget && !dTarget {
		return nil, nil
	}
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ImportSpec)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		s := n.(*ast.ImportSpec)
		path, _ := strconv.Unquote(s.Path.Value)
		if matchesList(path, denylist) {
			pass.Reportf(s.Pos(), "prohibited import package: %s", s.Path.Value)
			return
		}
		if !aTarget {
			return
		}
		if matchesList(path, allowlist) {
			return
		}
		if !isStandardImportPath(path) {
			pass.Reportf(s.Pos(), "prohibited import package: %s", s.Path.Value)
		}
	})
	return nil, nil
}

func matchesList(path string, patterns matcherSet) bool {
	if patterns.HasWildcard {
		return true
	}
	if _, exists := patterns.Exact[path]; exists {
		return true
	}
	for _, re := range patterns.Regex {
		if re.MatchString(path) {
			return true
		}
	}
	return false
}

// copied from https://pkg.go.dev/cmd/go/internal/search#IsStandardImportPath
func isStandardImportPath(path string) bool {
	i := strings.Index(path, "/")
	if i < 0 {
		i = len(path)
	}
	elem := path[:i]
	return !strings.Contains(elem, ".")
}
