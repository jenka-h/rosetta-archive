package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "rosetta: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer, stderr io.Writer) error {
	request, err := ParseCommand(args)
	if err != nil {
		printUsage(stderr)
		return err
	}
	return execute(request, stdout)
}
