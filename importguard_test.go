package importguard

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestWithAutoDiscoveredConfig(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer,
		"github.com/satorunooshie/example/cases/allimports",
		"github.com/satorunooshie/example/cases/blacklist",
		"github.com/satorunooshie/example/cases/domainregex",
		"github.com/satorunooshie/example/cases/exactdeny",
		"github.com/satorunooshie/example/cases/localoverride",
		"github.com/satorunooshie/example/cases/regexdeny",
		"github.com/satorunooshie/example/cases/mixedallow",
	)
}
