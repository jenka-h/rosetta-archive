package main

import (
	"errors"
	"fmt"
	"io"
	"rosetta-archive/internal/archive"
)

type Command string

const (
	CommandCreate  Command = "create"
	CommandList    Command = "list"
	CommandInfo    Command = "info"
	CommandInspect Command = "inspect"
	CommandVerify  Command = "verify"
	CommandExtract Command = "extract"
	CommandHelp    Command = "help"
)

// CommandRequest is the CLI DTO passed from argument parsing to command execution.
type CommandRequest struct {
	Command     Command
	ArchivePath string
	Sources     []string
	Destination string
}

// Parse Commands
func ParseCommand(args []string) (CommandRequest, error) {
	if len(args) == 0 {
		return CommandRequest{}, errors.New("missing command")
	}

	switch Command(args[0]) {
	case CommandCreate:
		if len(args) < 3 {
			return CommandRequest{}, errors.New("usage: rosetta create <archive.rosa> <path> [path...]")
		}
		return CommandRequest{Command: CommandCreate, ArchivePath: args[1], Sources: args[2:]}, nil
	case CommandList:
		if len(args) != 2 {
			return CommandRequest{}, errors.New("usage: rosetta list <archive.rosa>")
		}
		return CommandRequest{Command: CommandList, ArchivePath: args[1]}, nil
	case CommandInfo:
		if len(args) != 2 {
			return CommandRequest{}, errors.New("usage: rosetta info <archive.rosa>")
		}
		return CommandRequest{Command: CommandInfo, ArchivePath: args[1]}, nil
	case CommandInspect:
		if len(args) != 2 {
			return CommandRequest{}, errors.New("usage: rosetta inspect <archive.rosa>")
		}
		return CommandRequest{Command: CommandInspect, ArchivePath: args[1]}, nil
	case CommandVerify:
		if len(args) != 2 {
			return CommandRequest{}, errors.New("usage: rosetta verify <archive.rosa>")
		}
		return CommandRequest{Command: CommandVerify, ArchivePath: args[1]}, nil
	case CommandExtract:
		if len(args) != 3 {
			return CommandRequest{}, errors.New("usage: rosetta extract <archive.rosa> <destination>")
		}
		return CommandRequest{Command: CommandExtract, ArchivePath: args[1], Destination: args[2]}, nil
	case CommandHelp, "-h", "--help":
		return CommandRequest{Command: CommandHelp}, nil
	default:
		return CommandRequest{}, fmt.Errorf("unknown command %q", args[0])
	}
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
		fmt.Fprintf(stdout, "format: %s\n", info.FormatName)
		fmt.Fprintf(stdout, "version: %d.%d\n", info.MajorVersion, info.MinorVersion)
		fmt.Fprintf(stdout, "archive_size: %d\n", info.ArchiveSize)
		fmt.Fprintf(stdout, "entry_count: %d\n", info.EntryCount)
		fmt.Fprintf(stdout, "file_count: %d\n", info.FileCount)
		fmt.Fprintf(stdout, "directory_count: %d\n", info.DirectoryCount)
		fmt.Fprintf(stdout, "stored_size: %d\n", info.StoredSize)
		fmt.Fprintf(stdout, "uncompressed_size: %d\n", info.UncompressedSize)
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
