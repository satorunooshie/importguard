package internal

import (
	"fmt" // want "prohibited import package: \"fmt\""
	"time"

	"github.com/satorunooshie/repo/libs/collection" // want "prohibited import package: \"github.com/satorunooshie/repo/libs/collection\""
)

func Do() {
	_ = time.Now()
	_ = fmt.Sprintf("%x", 42)
	collection.Do()
}
