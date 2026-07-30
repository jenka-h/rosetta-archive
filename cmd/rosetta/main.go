package main

import (
	"errors"
	"fmt"
	"os"

	"rosetta-archive/internal/archive"
)

// entry point untuk cli
func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "rosetta: %v\n", err)
		os.Exit(1)
	}
}

// routing command cli
func run(args []string) error {
	if len(args) == 0 {
		printUsage(os.Stderr)
		return errors.New("missing command")
	}

	switch args[0] {
	case "create":
		if len(args) < 3 {
			return errors.New("usage: rosetta create <archive.rosa> <path> [path...]")
		}
		return archive.Create(args[1], args[2:])
	case "list":
		if len(args) != 2 {
			return errors.New("usage: rosetta list <archive.rosa>")
		}
		entries, err := archive.List(args[1])
		if err != nil {
			return err
		}
		for _, entry := range entries {
			fmt.Fprintln(os.Stdout, entry.Path)
		}
		return nil
	case "info":
		if len(args) != 2 {
			return errors.New("usage: rosetta info <archive.rosa>")
		}
		info, err := archive.Info(args[1])
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "entries: %d\n", info.EntryCount)
		return nil
	case "inspect":
		if len(args) != 2 {
			return errors.New("usage: rosetta inspect <archive.rosa>")
		}
		return archive.Inspect(args[1], os.Stdout)
	case "verify":
		if len(args) != 2 {
			return errors.New("usage: rosetta verify <archive.rosa>")
		}
		return archive.Verify(args[1])
	case "extract":
		if len(args) != 3 {
			return errors.New("usage: rosetta extract <archive.rosa> <destination>")
		}
		return archive.Extract(args[1], args[2])
	case "help", "-h", "--help":
		printUsage(os.Stdout)
		return nil
	default:
		printUsage(os.Stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

// help menu
func printUsage(out *os.File) {
	fmt.Fprintln(out, "usage: rosetta <command> [arguments]")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "commands:")
	fmt.Fprintln(out, "  create  create an archive")
	fmt.Fprintln(out, "  list    list archive entries")
	fmt.Fprintln(out, "  info    print archive summary")
	fmt.Fprintln(out, "  inspect inspect archive structure")
	fmt.Fprintln(out, "  verify  verify archive integrity")
	fmt.Fprintln(out, "  extract extract archive contents")
}
