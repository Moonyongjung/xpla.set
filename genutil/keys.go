package genutil

import (
	"os"
	"path"

	"github.com/Moonyongjung/xpla-set/types"
	"github.com/Moonyongjung/xpla-set/util"
	"github.com/ethereum/go-ethereum/common"

	xtypes "github.com/Moonyongjung/xpla.go/types"
	xutil "github.com/Moonyongjung/xpla.go/util"
	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/go-bip39"
	"github.com/evmos/ethermint/crypto/hd"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type SaveMnemonic struct {
	Mnemonic string `json:"mnemonic"`
}

func KeysAdd(validator Validator, home string, valNodeName string) (keyring.Keyring, error) {
	clientCtx, err := xutil.NewClient()
	if err != nil {
		return nil, err
	}

	var keyring keyring.Keyring

	keys := validator.Keys
	if len(validator.Keys) == 0 {
		return nil, util.LogErr(types.ErrInvalidRequest, "key must be exist at least one")
	}

	for i, key := range keys {
		k, err := runAdd(key, clientCtx, home, validator.KeysOption)
		if err != nil {
			return nil, err
		}

		if i == 0 {
			keyring = k
		}
	}

	return keyring, nil
}

func runAdd(key Key, clientCtx cmclient.Context, home string, keysOption KeysOption) (keyring.Keyring, error) {
	name := key.Name
	util.LogKV("target keyname", name)
	if name == "" || key.KeyringBackend == "" {
		return nil, util.LogErr(types.ErrInvalidMsgType, "name and keyring_backend are manatory")
	}

	if !(key.KeyringBackend == keyring.BackendFile || key.KeyringBackend == keyring.BackendTest) {
		return nil, util.LogErr(types.ErrInvalidMsgType, "keyring_backend type must be file or test")
	}

	k, err := keyring.New(
		xtypes.XplaToolDefaultName,
		key.KeyringBackend,
		home,
		os.Stdin,
		hd.EthSecp256k1Option(),
	)
	if err != nil {
		return nil, util.LogErr(types.ErrInvalidMsgType, err)
	}

	hdPath := sdk.GetConfig().GetFullBIP44Path()

	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, err
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return nil, err
	}

	if keysOption.PrintMnemonic {
		util.LogKV("mnemonic words", mnemonic)
	}

	info, err := k.NewAccount(name, mnemonic, keyring.DefaultBIP39Passphrase, hdPath, hd.EthSecp256k1)
	if err != nil {
		return nil, util.LogErr(types.ErrInvalidRequest, err)
	}

	if !keysOption.NotSaveMnemonic {
		var saveMnemonic SaveMnemonic
		saveMnemonic.Mnemonic = mnemonic
		jsonByte, err := xutil.JsonMarshalData(saveMnemonic)
		if err != nil {
			return nil, util.LogErr(types.ErrParse, err)
		}

		var keyPath string
		if key.KeyringBackend == keyring.BackendFile {
			keyPath = "keyring-file"
		} else {
			keyPath = "keyring-test"
		}

		fPath := path.Join(home, keyPath)
		mnemonicFilePath := path.Join(fPath, name+"_mnemonic.json")

		f, err := os.Create(mnemonicFilePath)
		if err != nil {
			return nil, util.LogErr(types.ErrParse, err)
		}
		defer f.Close()

		if _, err = f.Write(jsonByte); err != nil {
			return nil, util.LogErr(types.ErrParse, err)
		}
	}

	balance := key.Balance
	if balance == "" {
		balance = "0"
	}
	util.LogKV("balance of the account", balance)

	balanceBigInt, err := util.ConvSdkInt(balance)
	if err != nil {
		return nil, util.LogErr(types.ErrInvalidRequest, err)
	}

	types.TotalSupply = types.TotalSupply.Add(balanceBigInt)
	coins := sdk.Coins{
		sdk.NewCoin(xtypes.XplaDenom, balanceBigInt),
	}

	types.GenBalances = append(types.GenBalances, banktypes.Balance{
		Address: info.GetAddress().String(),
		Coins:   coins.Sort(),
	})
	types.GenAccounts = append(types.GenAccounts, &ethermint.EthAccount{
		BaseAccount: authtypes.NewBaseAccount(info.GetAddress(), nil, 0, 0),
		CodeHash:    common.BytesToHash(evmtypes.EmptyCodeHash).Hex(),
	})

	util.LogInfo("create " + name + " is done")

	return k, nil
}
