package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/satorunooshie/importguard/v4"
)

func main() { singlechecker.Main(importguard.Analyzer) }
