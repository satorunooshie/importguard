package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/satorunooshie/importguard/v3"
)

func main() { singlechecker.Main(importguard.Analyzer) }
