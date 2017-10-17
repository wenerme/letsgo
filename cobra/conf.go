package wcobra

import (
	"context"
	"reflect"
)

type CommandConfPost interface {
	PostInstall(cmd *Command) error
}
type CommandConf interface {
	Install(cmd *Command) error
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

		cmd.AppendPreRun(conf.PostInstall)
	}
	return conf
}

func (self *commandConf) Install(cmd *Command, conf ... CommandConf) error {
	for _, v := range conf {
		if self.confs[reflect.TypeOf(v)] == nil {
			if err := v.Install(cmd); err != nil {
				return err
			}
			self.confs[reflect.TypeOf(v)] = v
		}
	}
	return nil
}
func (self *commandConf) PostInstall(cmd *Command, args []string) error {
	for _, v := range self.confs {
		if i, ok := v.(CommandConfPost); ok {
			i.PostInstall(cmd)
		}
	}
	return nil
}
