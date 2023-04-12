package types

import (
	"os"
	"path/filepath"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var (
	DefaultHome     = filepath.Join(os.Getenv("HOME"), ".xplaset")
	GenFiles        []string
	GenAccounts     []authtypes.GenesisAccount
	GenBalances     []banktypes.Balance
	NodeIds         []string
	ValPubkeys      []cryptotypes.PubKey
	GentxsDirPath   string
	TotalSupply     = sdk.NewInt(0)
	SentryPeersList []string
	SentryInfos     = make(map[string][]SentryInfo)
)

type SentryInfo struct {
	NodeId string
	Ip     string
}
