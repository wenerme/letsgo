package wcobra

import (
	"context"
	"reflect"
	"github.com/spf13/pflag"
	"github.com/spf13/cobra"
)

type CommandConfPost interface {
	PostInstall(cmd *Command) error
}
type CommandConf interface {
	Install(fs *pflag.FlagSet, cmd *Command) error
}

type commandConf struct {
	confs map[reflect.Type]CommandConf
}

var commandConfKey = reflect.TypeOf(commandConf{})

func confOf(cmd *Command) *commandConf {
	conf, ok := cmd.Value(commandConfKey).(*commandConf)
	if !ok {
		conf = &commandConf{
			confs: make(map[reflect.Type]CommandConf),
		}
		cmd.Context = context.WithValue(cmd.Context, commandConfKey, conf)

		cobra.OnInitialize(func() {
			conf.PostInstall(cmd)
		})
	}
	return conf
}

func (self *commandConf) Install(fs *pflag.FlagSet, cmd *Command, conf ... CommandConf) error {
	for _, v := range conf {
		if self.confs[reflect.TypeOf(v)] == nil {
			if err := v.Install(fs, cmd); err != nil {
				return err
			}
			self.confs[reflect.TypeOf(v)] = v
		}
	}
	return nil
}
func (self *commandConf) PostInstall(cmd *Command) error {
	for _, v := range self.confs {
		if i, ok := v.(CommandConfPost); ok {
			i.PostInstall(cmd)
		}
	}
	return nil
}
