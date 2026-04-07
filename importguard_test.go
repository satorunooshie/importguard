package importguard

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestWithAutoDiscoveredConfig(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer,
		"github.com/satorunooshie/repo/client",
		"github.com/satorunooshie/repo/internal",
		"github.com/satorunooshie/repo/libs/collection",
		"github.com/satorunooshie/repo/libs/crypto",
	)
}
