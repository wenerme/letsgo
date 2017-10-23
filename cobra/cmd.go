package wcobra

import (
	"context"
	"github.com/spf13/cobra"
	"reflect"
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

func (self *Command) InstallPersistent(conf ...CommandConf) *Command {
	confOf(self).InstallPersistent(self, conf...)
	return self
}
func (self *Command) Install(conf ...CommandConf) *Command {
	confOf(self).Install(self, conf...)
	return self
}
