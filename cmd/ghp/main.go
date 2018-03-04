package main

import (
	"os"

	"github.com/Code-Hex/ghp"
)

func main() {
	os.Exit(ghp.New().Run())
}
