package genutil

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Moonyongjung/xpla-set/types"
	"github.com/Moonyongjung/xpla-set/util"
	xtypes "github.com/Moonyongjung/xpla.go/types"

	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	evmcfg "github.com/evmos/ethermint/server/config"
	tmcfg "github.com/tendermint/tendermint/config"
	"github.com/xpladev/xpla/app/params"
)

// Read config.yaml file.
func ReadConfig(configFilePath string) error {
	util.GuideReadConfig()
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return util.LogErr(types.ErrParse, err)
	}

	err := ConfigFile().Read(configFilePath)
	if err != nil {
		return util.LogErr(types.ErrParse, err)
	}

	return nil
}

// Create config files of the validator.
// These are included in the config directory of the validator.
func SetInit(rootDir string) error {

	xplaTemplate, xplaAppConfig := initAppConfig()

	configPath := filepath.Join(rootDir, "config")
	tmCfgFile := filepath.Join(configPath, "config.toml")
	util.LogKV("config path", configPath)

	conf := tmcfg.DefaultConfig()

	switch _, err := os.Stat(tmCfgFile); {
	case os.IsNotExist(err):
		tmcfg.EnsureRoot(rootDir)

		if err = conf.ValidateBasic(); err != nil {
			return util.LogErr(types.ErrInvalidRequest, "error in config file:", err)
		}

		conf.RPC.PprofListenAddress = "localhost:6060"
		conf.P2P.RecvRate = 5120000
		conf.P2P.SendRate = 5120000
		conf.Consensus.TimeoutCommit = 5 * time.Second
		tmcfg.WriteConfigFile(tmCfgFile, conf)

	case err != nil:
		return util.LogErr(types.ErrInvalidRequest, err)

	default:
		return util.LogErr(types.ErrInvalidRequest, "cannot set config file")
	}

	conf.SetRoot(rootDir)

	appCfgFilePath := filepath.Join(configPath, "app.toml")
	if _, err := os.Stat(appCfgFilePath); os.IsNotExist(err) {
		if xplaTemplate != "" {
			config.SetConfigTemplate(xplaTemplate)

			config.WriteConfigFile(appCfgFilePath, xplaAppConfig)
		}
	}
	return nil
}

// Initialize xpla app configuration file.
func initAppConfig() (string, interface{}) {
	customAppTemplate, customAppConfig := evmcfg.AppConfig(xtypes.XplaDenom)
	srvCfg, ok := customAppConfig.(evmcfg.Config)
	if !ok {
		panic(fmt.Errorf("unknown app config type %T", customAppConfig))
	}

	srvCfg.StateSync.SnapshotInterval = 1000
	srvCfg.StateSync.SnapshotKeepRecent = 10

	return customAppTemplate + params.CustomConfigTemplate, params.CustomAppConfig{
		Config: srvCfg,
		BypassMinFeeMsgTypes: []string{
			sdk.MsgTypeURL(&ibcchanneltypes.MsgRecvPacket{}),
			sdk.MsgTypeURL(&ibcchanneltypes.MsgAcknowledgement{}),
			sdk.MsgTypeURL(&ibcclienttypes.MsgUpdateClient{}),
		},
	}
}
