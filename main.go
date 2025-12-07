package main

import (
	"fmt"
	"os"

	"github.com/chickeniq/go-discovery/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
}
