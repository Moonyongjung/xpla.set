package util

import (
	xtypes "github.com/Moonyongjung/xpla.go/types"
)

func GuideReadConfig() {
	LogInfo(B("======================================================="))
	LogInfo(B("Read the config file"))
	LogInfo(B("======================================================="))
}

func GuideCreateValidator(validator string) {
	LogInfo(B("======================================================="))
	LogInfo(B("Create"), G(validator))
	LogInfo(B("======================================================="))
}

func GuideValidatorInfo(v xtypes.CreateValidatorMsg) {
	LogInfo(BB("=================== validator infos ==================="))
	LogKV("moniker", v.Moniker)
	LogKV("node IP", v.ServerIp)
	LogKV("identity", v.Identity)
	LogKV("website", v.Website)
	LogKV("security contact", v.SecurityContact)
	LogKV("details", v.Details)
	LogKV("self-delelgation amount", v.Amount)
	LogInfo(BB("======================================================="))
}

func GuideMakeGenesis() {
	LogInfo(B("======================================================="))
	LogInfo(B("Make the genesis file"))
	LogInfo(B("======================================================="))
}

func GuideCollectGetxs() {
	LogInfo(B("======================================================="))
	LogInfo(B("Collect gentxs"))
	LogInfo(B("======================================================="))
}
