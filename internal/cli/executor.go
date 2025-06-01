package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"golang.org/x/term"
)

type Executor struct {
	cli CLI
}

func NewExecutor(cli CLI) *Executor {
	return &Executor{cli: cli}
}

// ParseCommand splits a command string into arguments, respecting quotes
func ParseCommand(command string) []string {
	var args []string
	var currentArg strings.Builder
	inQuotes := false
	var quoteChar rune

	for _, char := range command {
		switch char {
		case '"', '\'':
			if inQuotes && char == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else {
				currentArg.WriteRune(char)
			}
		case ' ':
			if !inQuotes {
				if currentArg.Len() > 0 {
					args = append(args, currentArg.String())
					currentArg.Reset()
				}
			} else {
				currentArg.WriteRune(char)
			}
		default:
			currentArg.WriteRune(char)
		}
	}

	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return args
}

func (e *Executor) ExecuteInteractive(command string) error {
	args := ParseCommand(command)
	cmd := exec.Command(e.cli.(*BaseCLI).command, args...)

	// Start the command with a pty
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed to start pty: %w", err)
	}
	defer ptmx.Close()

	// Make sure raw mode is restored
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to set terminal raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Copy stdin to the pty
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ptmx.Write(scanner.Bytes())
			ptmx.Write([]byte("\n"))
		}
	}()

	// Copy pty to stdout
	scanner := bufio.NewScanner(ptmx)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return cmd.Wait()
}

func (e *Executor) ExecuteWithOutput(command string) (string, error) {
	return e.cli.Execute(command)
}

func IsInteractiveCommand(command string) bool {
	parts := strings.Split(command, " ")
	if len(parts) == 0 {
		return false
	}

	interactiveCommands := map[string]bool{
		"rsh":          true,
		"exec":         true,
		"debug":        true,
		"attach":       true,
		"logs":         true,
		"port-forward": true,
	}

	return interactiveCommands[parts[0]]
}
