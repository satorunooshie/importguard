package importguard

import (
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	t.Setenv("IMPORTGUARD_CONFIG", filepath.Join(testdata+"/config.json"))
	analysistest.Run(t, testdata, Analyzer,
		"github.com/satorunooshie/repo/client",
		"github.com/satorunooshie/repo/internal",
		"github.com/satorunooshie/repo/libs/collection",
		"github.com/satorunooshie/repo/libs/crypto",
	)
}
