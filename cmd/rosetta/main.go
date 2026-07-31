package main

import (
	"fmt"
	"io"
	"os"

	"rosetta-archive/internal/archive"
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

func execute(request CommandRequest, stdout io.Writer) error {
	switch request.Command {
	case CommandCreate:
		return archive.Create(request.ArchivePath, request.Sources)
	case CommandList:
		entries, err := archive.List(request.ArchivePath)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			fmt.Fprintln(stdout, entry.Path)
		}
		return nil
	case CommandInfo:
		info, err := archive.Info(request.ArchivePath)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "entries: %d\n", info.EntryCount)
		return nil
	case CommandInspect:
		return archive.Inspect(request.ArchivePath, stdout)
	case CommandVerify:
		return archive.Verify(request.ArchivePath)
	case CommandExtract:
		return archive.Extract(request.ArchivePath, request.Destination)
	case CommandHelp:
		printUsage(stdout)
		return nil
	default:
		return fmt.Errorf("unsupported command %q", request.Command)
	}
}

func printUsage(out io.Writer) {
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
