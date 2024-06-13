package serve

import (
	"flag"
	"fmt"

	"github.com/ncarlier/za/cmd"
	"github.com/ncarlier/za/pkg/config"
)

const cmdName = "serve"

type ServeCmd struct {
	flagSet *flag.FlagSet
}

func (c *ServeCmd) Exec(args []string, conf *config.Config) error {
	// no args
	startServer(conf)
	return nil
}

func (c *ServeCmd) Usage() {
	fmt.Fprintf(c.flagSet.Output(), "  %s\t\tStart ZerØ Analytics server\n", cmdName)
}

func newServeCmd() cmd.Cmd {
	c := &ServeCmd{
		flagSet: flag.NewFlagSet(cmdName, flag.ExitOnError),
	}
	return c
}

func init() {
	cmd.Add("serve", newServeCmd)
}
