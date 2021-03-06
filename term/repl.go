// SPDX-FileCopyrightText: 2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package term

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
)

var (
	rl *readline.Instance
	// PromptDatabaseName contains the used database name when using
	// the prompt.
	PromptDatabaseName string
	promptMultiline    bool
)

// UpdatePrompt updates the displayed prompt in interactive use.
func UpdatePrompt() {
	prompt := "> "

	if promptMultiline {
		prompt = ">>> "
	}

	if PromptDatabaseName != "" {
		prompt = PromptDatabaseName + prompt
	}

	if rl != nil {
		rl.SetPrompt(prompt)
	}
}

// Repl is the interactive interface that reads, evaluates, and prints
// the passed queries.
func Repl(db *sql.DB) error {
	var err error
	rl, err = readline.New("")
	if err != nil {
		return fmt.Errorf("term: failed to initialize readline: %w", err)
	}
	defer rl.Close()

	cmds := []string{}
	for {
		UpdatePrompt()

		exitAfterExecution := false

		line, err := rl.Readline()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// Execute the currently read line and then return
				exitAfterExecution = true
			} else {
				return fmt.Errorf("term: received error from readline: %w", err)
			}
		}

		line = strings.TrimSpace(line)
		if line != "" {
			cmds = append(cmds, line)
		}

		if !strings.HasSuffix(line, ";") && !exitAfterExecution {
			promptMultiline = true
			continue
		}

		// command is finished, reset and execute
		promptMultiline = false

		line = strings.Join(cmds, " ")
		cmds = []string{}

		err = ParseAndExecQueries(db, line)
		if exitAfterExecution {
			return err
		}

		if err != nil {
			log.Println(err)
		}
	}
}
