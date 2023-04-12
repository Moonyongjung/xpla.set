package genutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Moonyongjung/xpla.set/types"
	"github.com/Moonyongjung/xpla.set/util"

	"github.com/Moonyongjung/xpla.go/client"
	xtypes "github.com/Moonyongjung/xpla.go/types"
	xutil "github.com/Moonyongjung/xpla.go/util"
	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	tmcfg "github.com/tendermint/tendermint/config"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tendermint "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	"github.com/xpladev/xpla/app"
)

// Initialize genesis file.
// Record params to optimize xpla network.
func InitGenFile(xplac *client.XplaClient) (tendermint.GenesisDoc, error) {
	util.GuideMakeGenesis()
	util.LogInfo("initialize genesis file with optimize to xpla network...")
	zeroDec, _ := sdk.NewDecFromStr("0.000000000000000000")
	codec := xplac.EncodingConfig.Marshaler
	appGenState := app.ModuleBasics.DefaultGenesis(codec)

	// auth
	var authGenState authtypes.GenesisState
	codec.MustUnmarshalJSON(appGenState[authtypes.ModuleName], &authGenState)

	accounts, err := authtypes.PackAccounts(types.GenAccounts)
	if err != nil {
		return tendermint.GenesisDoc{}, util.LogErr(types.ErrParse, err)
	}

	authGenState.Accounts = accounts
	appGenState[authtypes.ModuleName] = codec.MustMarshalJSON(&authGenState)
	util.LogKV("optimize auth module", "done")

	// bank
	var bankGenState banktypes.GenesisState
	codec.MustUnmarshalJSON(appGenState[banktypes.ModuleName], &bankGenState)

	bankGenState.Balances = types.GenBalances
	bankGenState.Supply = sdk.Coins{
		sdk.NewCoin(xtypes.XplaDenom, types.TotalSupply),
	}
	appGenState[banktypes.ModuleName] = codec.MustMarshalJSON(&bankGenState)
	util.LogKV("optimize bank module", "done")

	// staking
	var stakingGenState stakingtypes.GenesisState
	codec.MustUnmarshalJSON(appGenState[stakingtypes.ModuleName], &stakingGenState)

	stakingGenState.Params.BondDenom = xtypes.XplaDenom
	stakingGenState.Params.MaxValidators = 8

	appGenState[stakingtypes.ModuleName] = codec.MustMarshalJSON(&stakingGenState)
	util.LogKV("optimize staking module", "done")

	// gov
	var govGenState govtypes.GenesisState
	codec.MustUnmarshalJSON(appGenState[govtypes.ModuleName], &govGenState)

	govGenState.DepositParams.MinDeposit[0].Denom = xtypes.XplaDenom
	govAmount, err := util.ConvSdkInt("10000000000000000000")
	if err != nil {
		return tendermint.GenesisDoc{}, util.LogErr(types.ErrParse, err)
	}
	govGenState.DepositParams.MinDeposit[0].Amount = govAmount
	govGenState.VotingParams.VotingPeriod = time.Second * 604800
	appGenState[govtypes.ModuleName] = codec.MustMarshalJSON(&govGenState)
	util.LogKV("optimize gov module", "done")

	// mint
	var mintGenState minttypes.GenesisState
	codec.MustUnmarshalJSON(appGenState[minttypes.ModuleName], &mintGenState)

	mintGenState.Minter.Inflation = zeroDec
	mintGenState.Minter.AnnualProvisions = zeroDec
	mintGenState.Params.MintDenom = xtypes.XplaDenom
	mintGenState.Params.InflationRateChange = zeroDec
	mintGenState.Params.InflationMax = zeroDec

	appGenState[minttypes.ModuleName] = codec.MustMarshalJSON(&mintGenState)
	util.LogKV("optimize mint module", "done")

	// crisis
	var crisisGenState crisistypes.GenesisState
	codec.MustUnmarshalJSON(appGenState[crisistypes.ModuleName], &crisisGenState)

	crisisGenState.ConstantFee.Denom = xtypes.XplaDenom
	crisisAmount, err := util.ConvSdkInt("1000000000000000")
	if err != nil {
		return tendermint.GenesisDoc{}, util.LogErr(types.ErrParse, err)
	}
	crisisGenState.ConstantFee.Amount = crisisAmount
	appGenState[crisistypes.ModuleName] = codec.MustMarshalJSON(&crisisGenState)
	util.LogKV("optimize crisis module", "done")

	// evm
	var evmGenState evmtypes.GenesisState
	codec.MustUnmarshalJSON(appGenState[evmtypes.ModuleName], &evmGenState)

	evmGenState.Params.EvmDenom = xtypes.XplaDenom
	appGenState[evmtypes.ModuleName] = codec.MustMarshalJSON(&evmGenState)
	util.LogKV("optimize evm module", "done")

	// slashing
	var slashingGenState slashingtypes.GenesisState
	codec.MustUnmarshalJSON(appGenState[slashingtypes.ModuleName], &slashingGenState)

	slashingGenState.Params.SignedBlocksWindow = 10000
	dec, _ := sdk.NewDecFromStr("0.005000000000000000")
	slashingGenState.Params.MinSignedPerWindow = dec
	dec, _ = sdk.NewDecFromStr("0.000100000000000000")
	slashingGenState.Params.SlashFractionDowntime = dec

	appGenState[slashingtypes.ModuleName] = codec.MustMarshalJSON(&slashingGenState)
	util.LogKV("optimize slashing module", "done")

	// distribution
	var distGenState disttypes.GenesisState
	codec.MustUnmarshalJSON(appGenState[disttypes.ModuleName], &distGenState)

	distGenState.Params.CommunityTax = zeroDec
	appGenState[disttypes.ModuleName] = codec.MustMarshalJSON(&distGenState)
	util.LogKV("optimize distribution module", "done")

	appGenStateJson, err := xutil.JsonMarshalData(appGenState)
	if err != nil {
		util.LogErr(types.ErrParse, err)
	}

	// tendermint consensus params
	consensusParams := tmproto.ConsensusParams{
		Block: tmproto.BlockParams{
			MaxBytes:   1000000,
			MaxGas:     100000000,
			TimeIotaMs: 1000,
		},
		Evidence: tmproto.EvidenceParams{
			MaxBytes:        1000000,
			MaxAgeNumBlocks: 100000,
			MaxAgeDuration:  172800000000000,
		},
		Validator: tmproto.ValidatorParams{
			PubKeyTypes: []string{tendermint.ABCIPubKeyTypeEd25519},
		},
	}

	genDoc := tendermint.GenesisDoc{
		ConsensusParams: &consensusParams,
		ChainID:         xplac.GetChainId(),
		AppState:        appGenStateJson,
		Validators:      nil,
		GenesisTime:     tmtime.Now(),
	}

	for _, genFile := range types.GenFiles {
		if err := genDoc.SaveAs(genFile); err != nil {
			util.LogErr(types.ErrParse, err)
		}
	}

	util.LogInfo(util.G("success to initialize genesis file"))

	return genDoc, nil
}

// Collect gentxs.
// Input gentxs which is create validator message to already created genesis file.
func CollectGenFiles(validators []Validator, home string, xplac *client.XplaClient, genDoc tendermint.GenesisDoc) error {
	util.GuideCollectGetxs()
	serverCtx := server.NewDefaultContext()
	nodeConfig := serverCtx.Config
	chainId := ConfigFile().Get().XplaGen.ChainId

	for i, validator := range validators {
		valNodeName := "validator" + xutil.FromIntToString(i)
		if validator.Moniker == "" {
			validator.Moniker = valNodeName
		}
		util.LogInfo("collect gentxs of", validator.Moniker+"...")

		valPath := path.Join(home, valNodeName)
		nodeConfig.SetRoot(valPath)
		nodeConfig.Moniker = validator.Moniker

		nodeId := types.NodeIds[i]
		valPubkey := types.ValPubkeys[i]

		initCfg := genutiltypes.NewInitConfig(chainId, types.GentxsDirPath, nodeId, valPubkey)
		nodeAppState, err := genAppStateFromConfig("", "", xplac.EncodingConfig.Marshaler, xplac.EncodingConfig.TxConfig, nodeConfig, initCfg, genDoc, banktypes.GenesisBalancesIterator{})
		if err != nil {
			return util.LogErr(types.ErrParse, err)
		}

		genDoc.AppState = nodeAppState
		if err := genDoc.SaveAs(types.GenFiles[i]); err != nil {
			util.LogErr(types.ErrParse, err)
		}

		// set sentry nodes
		if len(validator.Sentries.IpAddress) != 0 {
			for j, sentryInfo := range types.SentryInfos[nodeId] {
				sentryPath := path.Join(valPath, "sentry"+xutil.FromIntToString(j))
				nodeConfig.SetRoot(sentryPath)
				nodeConfig.Moniker = validator.Moniker + "-sentry" + xutil.FromIntToString(j)

				sentryNodeId := sentryInfo.NodeId

				initCfg := genutiltypes.NewInitConfig(chainId, types.GentxsDirPath, sentryNodeId, nil)
				nodeAppState, err := genAppStateFromConfig(nodeId, sentryNodeId, xplac.EncodingConfig.Marshaler, xplac.EncodingConfig.TxConfig, nodeConfig, initCfg, genDoc, banktypes.GenesisBalancesIterator{})
				if err != nil {
					return util.LogErr(types.ErrParse, err)
				}

				genDoc.AppState = nodeAppState
				if err := genDoc.SaveAs(path.Join(sentryPath, "config", "genesis.json")); err != nil {
					util.LogErr(types.ErrParse, err)
				}

			}
		}
	}

	util.LogInfo(util.G("success to collect gentxs"))

	return nil
}

// GenAppStateFromConfig gets the genesis app state from the config
func genAppStateFromConfig(valNodeId string, sentryNodeId string, cdc codec.JSONCodec, txEncodingConfig cmclient.TxEncodingConfig,
	config *tmcfg.Config, initCfg genutiltypes.InitConfig, genDoc tendermint.GenesisDoc, genBalIterator genutiltypes.GenesisBalancesIterator,
) (appState json.RawMessage, err error) {
	// process genesis transactions, else create default genesis.json
	appGenTxs, persistentPeers, err := genutil.CollectTxs(
		cdc, txEncodingConfig.TxJSONDecoder(), config.Moniker, initCfg.GenTxsDir, genDoc, genBalIterator,
	)
	if err != nil {
		return appState, err
	}

	// in case of the validator
	if valNodeId == "" {
		sentryInfos := types.SentryInfos[initCfg.NodeID]
		if len(sentryInfos) == 0 {
			config.P2P.PersistentPeers = persistentPeers
			config.P2P.UnconditionalPeerIDs = strings.Join(types.NodeIds, ",")

		} else {
			var sentryPeers []string
			var sentryNodeIds []string

			for _, sentryInfo := range sentryInfos {
				peer := fmt.Sprintf("%s@%s:26656", sentryInfo.NodeId, sentryInfo.Ip)
				sentryPeers = append(sentryPeers, peer)
				sentryNodeIds = append(sentryNodeIds, sentryInfo.NodeId)
			}
			sentries := strings.Join(sentryPeers, ",")
			config.P2P.PersistentPeers = persistentPeers + "," + sentries

			valIds := strings.Join(types.NodeIds, ",")
			sentryIds := strings.Join(sentryNodeIds, ",")
			config.P2P.UnconditionalPeerIDs = valIds + "," + sentryIds
		}
		config.P2P.AddrBookStrict = false
		config.P2P.MaxNumInboundPeers = 0
		config.P2P.MaxNumOutboundPeers = 0
		config.P2P.PexReactor = false
		config.P2P.PrivatePeerIDs = ""

		// in case of the sentry
	} else {
		var newPersistentPeers []string
		var newUnconditionalPeersIds []string

		splitedPeers := strings.Split(persistentPeers, ",")
		for _, peer := range splitedPeers {
			if strings.Contains(peer, valNodeId) {
				newPersistentPeers = append(newPersistentPeers, peer)
				newUnconditionalPeersIds = append(newUnconditionalPeersIds, valNodeId)
				break
			}
		}

		for _, sentryPeer := range types.SentryPeersList {
			id := strings.Split(sentryPeer, "@")
			if id[0] != sentryNodeId {
				newPersistentPeers = append(newPersistentPeers, sentryPeer)
				newUnconditionalPeersIds = append(newUnconditionalPeersIds, id[0])
			}
		}

		config.P2P.PersistentPeers = strings.Join(newPersistentPeers, ",")
		config.P2P.UnconditionalPeerIDs = strings.Join(newUnconditionalPeersIds, ",")
		config.P2P.AddrBookStrict = false
		config.P2P.PexReactor = true
		config.P2P.PrivatePeerIDs = valNodeId
	}

	config.RPC.PprofListenAddress = "localhost:6060"
	config.Consensus.TimeoutCommit = time.Second * 5

	tmcfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

	// if there are no gen txs to be processed, return the default empty state
	if len(appGenTxs) == 0 {
		return appState, errors.New("there must be at least one genesis tx")
	}

	// create the app state
	appGenesisState, err := genutiltypes.GenesisStateFromGenDoc(genDoc)
	if err != nil {
		return appState, err
	}

	appGenesisState, err = genutil.SetGenTxsInAppGenesisState(cdc, txEncodingConfig.TxJSONEncoder(), appGenesisState, appGenTxs)
	if err != nil {
		return appState, err
	}

	appState, err = json.MarshalIndent(appGenesisState, "", "  ")
	if err != nil {
		return appState, err
	}

	genDoc.AppState = appState
	err = genutil.ExportGenesisFile(&genDoc, config.GenesisFile())

	return appState, err
}
