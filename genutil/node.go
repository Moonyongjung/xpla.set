package genutil

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/Moonyongjung/xpla.set/types"
	"github.com/Moonyongjung/xpla.set/util"

	"github.com/Moonyongjung/xpla.go/client"
	"github.com/Moonyongjung/xpla.go/key"
	xtypes "github.com/Moonyongjung/xpla.go/types"
	xutil "github.com/Moonyongjung/xpla.go/util"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
)

// Set validator or sentry node info.
// Create not only node key but also priv validator key that is only in the validator's config directory.
func Node(validator Validator, valPath string, valNumber int, keyring keyring.Keyring, xplac *client.XplaClient) error {
	util.LogInfo(util.B("→ start make the validator node info..."))
	serverCtx := server.NewDefaultContext()
	nodeConfig := serverCtx.Config

	nodeConfig.SetRoot(valPath)
	nodeConfig.RPC.ListenAddress = "tcp://0.0.0.0:26657"
	nodeConfig.Moniker = validator.Moniker

	// create node key and private validator key
	nodeId, valPubkey, err := genutil.InitializeNodeValidatorFiles(nodeConfig)
	if err != nil {
		return util.LogErr(types.ErrParse, err)
	}

	util.LogKV("validator node ID", nodeId)
	util.LogKV("validator public key", valPubkey.String())

	types.GenFiles = append(types.GenFiles, nodeConfig.GenesisFile())
	types.NodeIds = append(types.NodeIds, nodeId)
	types.ValPubkeys = append(types.ValPubkeys, valPubkey)

	nodeKeyPath := path.Join(valPath, "config", "node_key.json")
	nodeKeyBytes, err := os.ReadFile(nodeKeyPath)
	if err != nil {
		return util.LogErr(types.ErrNotFound, err)
	}

	privValidatorKeyPath := path.Join(valPath, "config", "priv_validator_key.json")
	privValidatorKeyBytes, err := os.ReadFile(privValidatorKeyPath)
	if err != nil {
		return util.LogErr(types.ErrNotFound, err)
	}

	util.LogKV("node key", nodeKeyPath)
	util.LogKV("priv validator key", privValidatorKeyPath)

	if len(validator.Sentries.IpAddress) != 0 {
		for i, sentry := range validator.Sentries.IpAddress {
			sentryPath := path.Join(valPath, "sentry"+xutil.FromIntToString(i))
			nodeConfig.SetRoot(sentryPath)
			util.LogKV("sentry path", sentryPath)

			err = util.GenFilePath(sentryPath)
			if err != nil {
				return util.LogErr(types.ErrInvalidRequest, err)
			}

			sentryNodeId, _, err := genutil.InitializeNodeValidatorFiles(nodeConfig)
			if err != nil {
				return util.LogErr(types.ErrParse, err)
			}

			util.LogKV("sentry node ID", sentryNodeId)

			sentryIp, err := util.GetIP(sentry)
			if err != nil {
				return util.LogErr(types.ErrParse, err)
			}

			types.SentryInfos[nodeId] = append(types.SentryInfos[nodeId], types.SentryInfo{
				NodeId: sentryNodeId,
				Ip:     sentryIp,
			})
			types.SentryPeersList = append(types.SentryPeersList, fmt.Sprintf("%s@%s:26656", sentryNodeId, sentryIp))

		}
		util.LogInfo("sentry nodes of the " + validator.Moniker + " are ready to set")
	}

	util.LogInfo(util.B("← sucess make the " + validator.Moniker + "'s node info"))

	err = createValidatorMsg(string(nodeKeyBytes), string(privValidatorKeyBytes), validator, keyring, valPath, nodeId, xplac)
	if err != nil {
		return err
	}

	return nil
}

// Generate 'createValidatorMsg' of staking module for gentx.
func createValidatorMsg(nodeKey string, privValidatorKey string, validator Validator, k keyring.Keyring, valPath string, nodeId string, xplac *client.XplaClient) error {
	util.LogInfo(util.B("→ start make the gentx of the " + validator.Moniker + "..."))
	valKeyName := validator.Keys[0].Name
	armored, err := k.ExportPrivKeyArmor(valKeyName, "")
	if err != nil {
		return util.LogErr(types.ErrInvalidRequest, err)
	}

	privKey, _, err := key.UnarmorDecryptPrivKey(armored, "")
	if err != nil {
		return util.LogErr(types.ErrParse, err)
	}

	valAddr := sdk.ValAddress(privKey.PubKey().Address())

	gentxDirPath := path.Join(valPath, "config", "gentx")
	err = util.GenFilePath(gentxDirPath)
	if err != nil {
		return util.LogErr(types.ErrInvalidRequest, err)
	}

	gentxPath := path.Join(gentxDirPath, "gentx-"+nodeId+".json")
	xplac.WithOptions(
		client.Options{
			PrivateKey:     privKey,
			OutputDocument: gentxPath,
		},
	)

	if validator.DelAmount == "" {
		return util.LogErr(types.ErrInvalidRequest, "self-delegation amount must be set")
	}

	createValidatorMsg := xtypes.CreateValidatorMsg{
		NodeKey:                 nodeKey,
		PrivValidatorKey:        privValidatorKey,
		ValidatorAddress:        valAddr.String(),
		Moniker:                 validator.Moniker,
		Identity:                validator.ValidatorOption.Identity,
		Website:                 validator.ValidatorOption.Website,
		SecurityContact:         validator.ValidatorOption.SecurityContact,
		Details:                 validator.ValidatorOption.Details,
		Amount:                  xutil.DenomAdd(validator.DelAmount),
		CommissionRate:          validator.CommissionOption.Rate,
		CommissionMaxRate:       validator.CommissionOption.MaxRate,
		CommissionMaxChangeRate: validator.CommissionOption.MaxChangeRate,
		MinSelfDelegation:       validator.MinSelfDelegation,
		ServerIp:                validator.IpAddress,
	}

	// create and sign transaction
	_, err = xplac.CreateValidator(createValidatorMsg).CreateAndSignTx()
	if err != nil {
		return util.LogErr(types.ErrParse, err)
	}

	util.GuideValidatorInfo(createValidatorMsg)

	file, err := os.Open(gentxPath)
	if err != nil {
		return util.LogErr(types.ErrParse, err)
	}
	defer file.Close()

	// copy to gentxs folder
	copyToGentxs := path.Join(types.GentxsDirPath, validator.Moniker+".json")
	copy, err := os.Create(copyToGentxs)
	if err != nil {
		return util.LogErr(types.ErrParse, err)
	}
	defer copy.Close()

	_, err = io.Copy(copy, file)
	if err != nil {
		return util.LogErr(types.ErrParse, err)
	}

	util.LogKV("gentx path", gentxPath)
	util.LogInfo(util.B("← success to create gen tx"))

	return nil
}
