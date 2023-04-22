package main

import (
	"fmt"
	"os"

	damonuntil "github.com/clkbug/damonutil"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		os.Exit(1)
	}
}

func run() error {
	input := "damon.data"

	if len(os.Args) > 1 {
		input = os.Args[1]
	}

	damon, err := damonuntil.ParseDamonFile(input)
	if err != nil {
		return err
	}

	fmt.Printf("Version: %d\n", damon.Version)

	return nil
}
