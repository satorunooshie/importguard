package allimports

import (
	"time"

	"github.com/satorunooshie/example/targets/blocked" // want "prohibited import package: \"github.com/satorunooshie/example/targets/blocked\""
	"github.com/satorunooshie/example/targets/family"
)

func Do() {
	_ = time.Now()
	family.Do()
	blocked.Do()
}
