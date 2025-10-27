package main

import (
	"github.com/DanielAverlon/golangci-linters/internal/rules"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(rules.NoCommas)
}
