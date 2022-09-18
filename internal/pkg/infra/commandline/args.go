package commandline

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/mattn/go-shellwords"
)

func ParseArgs(command string, obg interface{}) error {
	args, err := shellwords.Parse(command)
	if err != nil {
		return fmt.Errorf("ResolveArgs: %w", err)
	}

	if _, err := flags.ParseArgs(obg, args); err != nil {
		return fmt.Errorf("ResolveArgs: %w", err)
	}

	return nil
}
