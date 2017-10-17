package wcobra

import (
	"github.com/spf13/cobra"
	"reflect"
	"context"
)

var commends = make(map[*cobra.Command]*Command)

func Wrap(cmd *cobra.Command) *Command {
	c := commends[cmd]
	if c == nil {
		c = &Command{
			Command: cmd,
			Context: context.WithValue(context.Background(), reflect.TypeOf(cmd), c),
		}
		commends[cmd] = c
	}
	return c
}

type Command struct {
	*cobra.Command
	context.Context
}

func (self *Command) AppendPreRun(f func(cmd *Command, args []string) error) *Command {
	self.PreRunE = AppendRunE(self.PreRunE, f)
	return self
}

func (self *Command) Install(conf ... CommandConf) *Command {
	confOf(self).Install(self, conf...)
	return self
}
