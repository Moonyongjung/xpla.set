package genutil

type ConfigType struct {
	XplaGen XplaGen `yaml:"xpla_gen"`
}

type XplaGen struct {
	ChainId    string      `yaml:"chain_id"`
	Home       string      `yaml:"home"`
	Validators []Validator `yaml:"validators"`
}

type Validator struct {
	Moniker           string           `yaml:"moniker"`
	IpAddress         string           `yaml:"ip_address"`
	DelAmount         string           `yaml:"del_amount"`
	MinSelfDelegation string           `yaml:"min_self_delegation"`
	CommissionOption  CommissionOption `yaml:"commission_option"`
	ValidatorOption   ValidatorOption  `yaml:"validator_option"`
	Keys              []Key            `yaml:"local_keys"`
	KeysOption        KeysOption       `yaml:"keys_option"`
	Sentries          Sentries         `yaml:"sentries"`
}

type CommissionOption struct {
	Rate          string `yaml:"rate"`
	MaxRate       string `yaml:"max_rate"`
	MaxChangeRate string `yaml:"max_change_rate"`
}

type ValidatorOption struct {
	Website         string `yaml:"website"`
	Identity        string `yaml:"identity"`
	SecurityContact string `yaml:"security_contact"`
	Details         string `yaml:"details"`
}

type Key struct {
	Name           string `yaml:"name"`
	KeyringBackend string `yaml:"keyring_backend"`
	Balance        string `yaml:"balance"`
}

type KeysOption struct {
	NotSaveMnemonic bool `yaml:"not_save_mnemonic"`
	PrintMnemonic   bool `yaml:"print_mnemonic"`
}

type Sentries struct {
	IpAddress []string `yaml:"ip_address"`
}
