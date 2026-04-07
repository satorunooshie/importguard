package regexdeny

import "github.com/satorunooshie/example/targets/blocked" // want "prohibited import package: \"github.com/satorunooshie/example/targets/blocked\""

func Do() {
	blocked.Do()
}
