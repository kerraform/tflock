package lock

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/logutils"
	"github.com/hashicorp/terraform/command"
	"github.com/hashicorp/terraform/command/arguments"
	"github.com/hashicorp/terraform/command/clistate"
	"github.com/hashicorp/terraform/command/views"
	"github.com/mitchellh/cli"
)

var _ cli.Command = (*LockCommand)(nil)

// LockCommand is a Command implementation that lock a Terraform state.
type LockCommand struct {
	command.StateMeta
}

// Run runs the procedure of this command.
func (c *LockCommand) Run(args []string) int {
	// Read the from state
	stateFromMgr, err := c.State()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error loading the state: %s", err))
		return 1
	}

	view := views.NewStateLocker(arguments.ViewHuman, c.View)
	stateLocker := clistate.NewLocker((0 * time.Second), view)
	if err := stateLocker.Lock(stateFromMgr, "tflock"); err != nil {
		c.Ui.Error(fmt.Sprintf("Error locking source state: %s", err))
		return 1
	}

	return 0
}

// Help returns long-form help text.
func (*LockCommand) Help() string {
	helpText := `
Usage: tflock
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns one-line help text.
func (c *LockCommand) Synopsis() string {
	return "Lock your Terraform state"
}

func LogOutput() io.Writer {
	levels := []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"}
	minLevel := os.Getenv("TFLOCK_LOG")

	// default log writer is null device.
	writer := ioutil.Discard
	if minLevel != "" {
		writer = os.Stderr
	}

	filter := &logutils.LevelFilter{
		Levels:   levels,
		MinLevel: logutils.LogLevel(minLevel),
		Writer:   writer,
	}

	return filter
}
