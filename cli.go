package main

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/chzyer/readline"
)

func buildCompleter() readline.AutoCompleter {
	cmds := []readline.PrefixCompleterInterface{}
	for k := range CmdMap {
		cmds = append(cmds, readline.PcItem(k))
	}
	return readline.NewPrefixCompleter(cmds...)
}

func getHomeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	}
	return os.Getenv(env)
}

// Function to tokenize the input string with quoted string support
func tokenize(input string) ([]string, error) {
	// Regular expression to match quoted strings or words
	re := regexp.MustCompile(`"([^"]*)"|'([^']*)'|\S+`)

	// Find all matches and return as a slice
	matches := re.FindAllString(input, -1)

	return matches, nil
}

// StartCli starts the repl environment
func StartCli() {
	historyFileDir := filepath.Join(getHomeDir(), ".cache")
	if _, err := os.Stat(historyFileDir); os.IsNotExist(err) {
		// simply ignore error since the history feature is optional.
		os.Mkdir(historyFileDir, 0644)
	}
	l, err := readline.NewEx(&readline.Config{
		AutoComplete:    buildCompleter(),
		Prompt:          DbPath + "> ",
		HistoryFile:     filepath.Join(historyFileDir, "boltclihistory"),
		HistoryLimit:    1000,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		fields, _ := tokenize(strings.TrimSpace(line))
		if len(fields) == 0 {
			continue
		}
		for i, f := range fields {
			fields[i] = strings.Trim(f, "'")
		}

		result := ExecCmdInCli(fields[0], fields[1:]...)
		if result != "" {
			println(result)
		} else {
			println("(empty list or set)")
		}
	}
}
