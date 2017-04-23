package main

import (
	"github.com/Focinfi/sqs/config"
	"github.com/Focinfi/sqs/master"
)

func main() {
	master.NewService(config.Config.DefaultMasterAddress).Start()
}
