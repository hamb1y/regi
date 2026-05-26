package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	appName         = "regi"
	defaultRegister = "default"
	fileExt         = ".regi"
	registerDir     = ".config/regi/registers"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return listRegister(defaultRegister)
	}

	switch args[0] {
	case "help", "--help", "-h":
		printGeneralHelp()
		return nil
	case "add":
		return addCommand(args[1:])
	case "del", "delete", "rm", "remove":
		return delCommand(args[1:])
	default:
		if len(args) > 1 {
			return fmt.Errorf("unknown command or too many arguments: %s", strings.Join(args, " "))
		}
		return listRegister(args[0])
	}
}

func addCommand(args []string) error {
	if wantsHelp(args) {
		printAddHelp()
		return nil
	}
	if len(args) == 0 {
		return errors.New("add requires text, or a register and text")
	}

	register, text := registerAndText(args)
	text = strings.TrimRight(text, "\r\n")
	if text == "" {
		return errors.New("cannot add an empty item")
	}

	path, err := registerPath(register)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create register dir: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	if _, err := fmt.Fprintln(f, text); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	return nil
}

func delCommand(args []string) error {
	if wantsHelp(args) {
		printDelHelp()
		return nil
	}
	if len(args) == 0 {
		return errors.New("del requires a regex, or a register and regex")
	}

	register, pattern := registerAndText(args)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex: %w", err)
	}

	path, err := registerPath(register)
	if err != nil {
		return err
	}

	items, err := readItems(path)
	if err != nil {
		return err
	}

	kept := items[:0]
	removed := 0
	for _, item := range items {
		if re.MatchString(item) {
			removed++
			continue
		}
		kept = append(kept, item)
	}

	return writeItems(path, kept, removed > 0)
}

func listRegister(register string) error {
	path, err := registerPath(register)
	if err != nil {
		return err
	}

	items, err := readItems(path)
	if err != nil {
		return err
	}

	for _, item := range items {
		fmt.Println(item)
	}
	return nil
}

func readItems(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	if len(data) == 0 {
		return nil, nil
	}

	text := strings.TrimSuffix(string(data), "\n")
	if text == "" {
		return nil, nil
	}
	return strings.Split(text, "\n"), nil
}

func writeItems(path string, items []string, changed bool) error {
	if !changed {
		return nil
	}

	text := ""
	if len(items) > 0 {
		text = strings.Join(items, "\n") + "\n"
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create register dir: %w", err)
	}
	if err := os.WriteFile(path, []byte(text), 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func registerPath(register string) (string, error) {
	name := strings.TrimSpace(register)
	if name == "" {
		name = defaultRegister
	}
	name = strings.TrimSuffix(name, fileExt)

	if name == "" {
		name = defaultRegister
	}
	if name == "." || name == ".." || filepath.Base(name) != name {
		return "", fmt.Errorf("invalid register name %q", register)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("find home dir: %w", err)
	}
	return filepath.Join(home, registerDir, name+fileExt), nil
}

func registerAndText(args []string) (string, string) {
	if len(args) == 1 {
		return defaultRegister, args[0]
	}
	return args[0], strings.Join(args[1:], " ")
}

func wantsHelp(args []string) bool {
	return len(args) == 1 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help")
}

func printGeneralHelp() {
	fmt.Printf(`%s stores newline-delimited plaintext registers in ~/.config/regi/registers/*.regi.

Usage:
  regi [register]
  regi add [register] <text>
  regi del [register] <regex>
  regi help
  regi <subcommand> -h

Examples:
  regi
  regi work
  regi add "buy milk"
  regi add work call Sam
  regi del work "^done:"

`, appName)
}

func printAddHelp() {
	fmt.Println(`Usage:
  regi add <text>
  regi add <register> <text>

Adds one plaintext line to a register. With no register, uses default.regi.`)
}

func printDelHelp() {
	fmt.Println(`Usage:
  regi del <regex>
  regi del <register> <regex>

Removes lines matching the Go regular expression. With no register, uses default.regi.`)
}
