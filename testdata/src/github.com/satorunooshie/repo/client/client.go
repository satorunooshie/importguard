package client

import (
	"github.com/satorunooshie/repo/libs/collection"
	"github.com/satorunooshie/repo/libs/crypto" // want "prohibited import package: \"github.com/satorunooshie/repo/libs/crypto\""
)

func Do() {
	crypto.Do()
	collection.Do()
}
