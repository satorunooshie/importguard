package blacklist

import (
	"fmt" // want "prohibited import package: \"fmt\""

	"github.com/satorunooshie/example/targets/blocked" // want "prohibited import package: \"github.com/satorunooshie/example/targets/blocked\""
	"github.com/satorunooshie/example/targets/exact"
)

func Do() {
	_ = fmt.Sprintf("%d", exact.Now())
	blocked.Do()
}
