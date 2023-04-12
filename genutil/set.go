package genutil

import (
	"os"

	"github.com/Moonyongjung/xpla.set/types"
	"github.com/Moonyongjung/xpla.set/util"

	"github.com/Moonyongjung/xpla.go/client"
)

// Xpla node setting can be used after creating xpla binary through make install.
// This function is the entry point to create node by implementing sub-functions,
// such as creating gentx of each validator, initalizing genesis file, collecting gentx and etc.
func Set(configFilePath string) error {
	util.LogInfo(util.B("start read base configuration..."))
	// read config file
	err := ReadConfig(configFilePath)
	if err != nil {
		return err
	}

	// read chain ID of the config file.
	chainId := ConfigFile().Get().XplaGen.ChainId
	if chainId == "" {
		return util.LogErr(types.ErrInvalidRequest, "chain ID must be set")
	}
	util.LogKV("chain ID", chainId)

	// set xpla client of xpla.go
	xplac := client.NewXplaClient(chainId)

	// set root directory
	home := ConfigFile().Get().XplaGen.Home
	if home == "" {
		home = types.DefaultHome
	}
	util.LogKV("home directory", home)
	util.LogInfo(util.G("success to read base configuration"))

	if util.IsExistPath(home) {
		// select removing home files or not
		util.LogWarning("already exist", home, "directory")
		util.LogWarning("Do you want to reset validator nodes through remove home dir? [y/N]")

		if util.GetConfirm() {
			err = os.RemoveAll(home)
			if err != nil {
				return util.LogErr(types.ErrInvalidRequest, err)
			}
		} else {
			return util.LogErr(types.ErrAlreadyExist, "not reset")
		}
	}

	// set info of the validator nodes
	err = Validators(ConfigFile().Get().XplaGen.Validators, home, xplac)
	if err != nil {
		os.RemoveAll(home)
		return err
	}

	// initialize the genesis file
	genDoc, err := InitGenFile(xplac)
	if err != nil {
		os.RemoveAll(home)
		return err
	}

	// collect gentxs
	err = CollectGenFiles(ConfigFile().Get().XplaGen.Validators, home, xplac, genDoc)
	if err != nil {
		os.RemoveAll(home)
		return err
	}

	return nil
}
