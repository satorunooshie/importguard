package importguard

import (
	"encoding/json"
	"go/ast"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type config struct {
	Allow map[string]map[string]struct{} `json:"allow"`
	Deny  map[string]map[string]struct{} `json:"deny"`
}

var (
	conf config
	once sync.Once
)

var Analyzer = &analysis.Analyzer{
	Name: "importguard",
	Doc:  "importguard reports prohibited imports",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func parseConfig() error {
	fp := os.Getenv("IMPORTGUARD_CONFIG")
	if fp == "" {
		return nil
	}
	b, err := os.ReadFile(fp)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &conf)
}

func run(pass *analysis.Pass) (any, error) {
	once.Do(func() {
		if err := parseConfig(); err != nil {
			panic(err)
		}
	})

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
