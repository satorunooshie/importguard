package crypto

import (
	"math/rand/v2"

	"github.com/satorunooshie/repo/libs/collection" // want "prohibited import package: \"github.com/satorunooshie/repo/libs/collection\""
)

func Do() {
	_ = rand.New(rand.NewPCG(1, 2))
	collection.Do()
}
