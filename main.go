package main

import (
	"goMedia/core"
	"goMedia/flag"
	"goMedia/global"
	"goMedia/initialize"
)

func main() {
	initialize.InitLogo()
	global.Config = core.InitConfig()
	global.Log = core.InitLogger()
	initialize.OtherInit()

	global.DB = initialize.InitGorm()

	flag.InitFlag()

	core.RunServer()
}
