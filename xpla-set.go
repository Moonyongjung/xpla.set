package main

import (
	"github.com/Moonyongjung/xpla.set/genutil"
	"github.com/Moonyongjung/xpla.set/util"
)

const configFilePath = "./config.yaml"

func main() {
	err := genutil.Set(configFilePath)
	if err == nil {
		util.LogInfo(util.BgG("xpla.set successfully complete."))
	}
}
