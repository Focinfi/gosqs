package main

import (
	"github.com/Focinfi/gosqs/config"
	"github.com/Focinfi/gosqs/master"
)

func main() {
	master.NewService(config.Config.DefaultMasterAddress).Start()
}
