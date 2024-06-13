package cmd

import "github.com/ncarlier/za/pkg/config"

type Cmd interface {
	Exec(args []string, conf *config.Config) error
	Usage()
}
