package importguard

import (
	"encoding/json"
	"errors"
	"go/ast"
	"os"
	"path/filepath"
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

const configFileName = ".importguard.json"

var Analyzer = &analysis.Analyzer{
	Name: "importguard",
	Doc:  "importguard reports prohibited imports",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func parseConfig(pass *analysis.Pass) (config, error) {
	fp, err := resolveConfigPath(pass)
	if err != nil {
		return config{}, err
	}
	if fp == "" {
		return config{}, nil
	}
	return loadConfig(fp)
}

func loadConfig(fp string) (config, error) {
	b, err := os.ReadFile(fp)
	if err != nil {
		return config{}, err
	}
	var conf config
	err = json.Unmarshal(b, &conf)
	if err != nil {
		return config{}, err
	}
	return conf, nil
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
		if _, exists := allowlist[path]; exists {
			return
		}
		if _, exists := denylist[path]; exists || !isStandardImportPath(path) {
			pass.Reportf(s.Pos(), "prohibited import package: %s", s.Path.Value)
		}
	})
	return nil, nil
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
