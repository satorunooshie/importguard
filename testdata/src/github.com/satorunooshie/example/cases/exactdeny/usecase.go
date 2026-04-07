package exactdeny

import (
	"github.com/satorunooshie/example/targets/blocked" // want "prohibited import package: \"github.com/satorunooshie/example/targets/blocked\""
	"github.com/satorunooshie/example/targets/family"
)

func Do() {
	family.Do()
	blocked.Do()
}
