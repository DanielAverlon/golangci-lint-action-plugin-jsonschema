package main

import (
	"github.com/golangci/golangci-lint-action-plugin-example/internal/rules"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(rules.NoDescriptionCommas)
}
