package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/satorunooshie/importguard"
)

func main() { singlechecker.Main(importguard.Analyzer) }
