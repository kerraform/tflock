package main

import (
	"fmt"
	"log"
	"os"

	backendInit "github.com/hashicorp/terraform/backend/init"
	"github.com/hashicorp/terraform/command"
	"github.com/kerraform/tflock/internal/lock"
	"github.com/mitchellh/cli"
)

func main() {
	log.SetOutput(lock.LogOutput())

	backendInit.Init(nil)

	UI := &cli.BasicUi{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	meta := command.Meta{
		Ui: UI,
	}

	commands := map[string]cli.CommandFactory{
		"": func() (cli.Command, error) {
			return &lock.LockCommand{
				StateMeta: command.StateMeta{
					Meta: meta,
				},
			}, nil
		},
	}

	args := os.Args[1:]

	c := &cli.CLI{
		Name:       "tflock",
		Version:    Version,
		Args:       args,
		Commands:   commands,
		HelpWriter: os.Stdout,
	}

	exitStatus, err := c.Run()
	if err != nil {
		UI.Error(fmt.Sprintf("Failed to execute CLI: %s", err))
	}

	os.Exit(exitStatus)
}
