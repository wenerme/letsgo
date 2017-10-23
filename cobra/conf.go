package wcobra

import (
	"context"
	"github.com/spf13/pflag"
	"reflect"
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
	}
	return conf
}

func (self *commandConf) InstallPersistent(cmd *Command, conf ...CommandConf) error {
	return self.install(cmd, true, conf...)
}
func (self *commandConf) Install(cmd *Command, conf ...CommandConf) error {
	return self.install(cmd, false, conf...)
}

func (self *commandConf) install(cmd *Command, persists bool, conf ...CommandConf) error {
	fs := cmd.Flags()
	if persists {
		fs = cmd.PersistentFlags()
	}
	for _, v := range conf {
		if self.confs[reflect.TypeOf(v)] == nil {
			if err := v.Install(fs, cmd); err != nil {
				return err
			}
			self.confs[reflect.TypeOf(v)] = v

			// Only post install if command run
			if p, ok := v.(CommandConfPost); ok {
				if persists {
					cmd.PersistentPreRunE = AppendRunE(cmd.PersistentPreRunE, func(cmd *Command, args []string) error {
						return p.PostInstall(cmd)
					})
				} else {
					cmd.PreRunE = AppendRunE(cmd.PreRunE, func(cmd *Command, args []string) error {
						return p.PostInstall(cmd)
					})
				}
			}
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
