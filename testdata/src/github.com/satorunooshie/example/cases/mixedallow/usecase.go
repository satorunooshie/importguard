package mixedallow

import (
	"github.com/satorunooshie/example/targets/exact"
	"github.com/satorunooshie/example/targets/family/nested"
)

func Do() {
	nested.Do()
	exact.Now()
}
