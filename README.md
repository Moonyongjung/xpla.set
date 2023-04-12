# xpla.set
Xpla.set is the validator node setting tool for xpla members such as new validators that want to participate xpla network, tester and xpla developers. This tool is able to handle from chain initialization to collecting gentxs of validators. Multiple validator sets can be created by inputting the desired number of validator information, and necessary information to operate the validator, such as config file, keyring, genesis file, sentry nodes, and gentx, is created in each set.

## Prerequisites
- Make the `xplad` binary file by running [xpla](https://github.com/xpladev/xpla) core.

## Start xpla.set
### Configuration
To start xpla.set, write the config.yaml file

```yaml
xpla_gen:
  chain_id: cube_47-5
  home: $HOME/.xplaset
  validators:
    - validator:
      moniker: validator_node_0
      ip_address: 172.26.0.10 
      del_amount: 500000000000000000000000 
      min_self_delegation: 1 

      commission_option: 
        rate: 0.1
        max_rate: 0.2
        max_change_rate: 0.01

      validator_option: 
        website: "https://xpla.io"
        identity: "IDENTITY"
        security_contact: "contact target"
        details: "details of the validator"
      
      local_keys: 
        - name: myKey1 
          keyring_backend: test 
          balance: 1000000000000000000000000 
        - name: myKey2 
          keyring_backend: test 
          balance: 12345 
        - name: myKey3 
          keyring_backend: test 

      keys_option: 
        not_save_mnemonic: false 
        print_mnemonic: false

      sentries: 
        ip_address: 
          - 172.26.0.30
          - 172.26.0.31
```
Each param is seperated to mandatory or optional `(M/O)`. 

- `chain_id`: The chain ID of generating chain. `(M)`
- `home`: Root directory of validator nodes. `(O, default "$HOME/.xplaset")`
- `validators`: The list of validators. If you want to create multiple validators, write `validator` parameters in `validators` according to each required setting. `(M)`
  - `validator`: Information of a validator. `(M)`
    - `moniker`: The validator's moniker. `(O, default "validator0, 1, 2...")`
    - `ip_address`: IP address of the validator node. `(O, default "detected node IP")`
    - `del_amount`: Self-delegated amount of the validator. It must be less than validator balance and defalut denom is `axpla`. `(M)`
    - `min_self_delegation`: Minimum self delegation. `(M)`
    - `commission_option`: Commission options of the validator. `(O)`
      - `rate`: Commission rate. `(O, default "0.1")`
      - `max_rate`: Commission max rate. `(O, default "0.2")`
      - `max_change_rate`: Commission max change rate. `(O, default "0.01")`
    - `validator_option`: The validator option. `(O)`
      - `website`: Introduce website of the validator. `(O, default "")`
      - `identity`: Identity of the validator. `(O, default "")`
      - `security_contact`: The contact point of the validator. `(O, default "")`
      - `details`: Detail information of the validator. `(O, default "")`
    - `local_keys`: Keyring files that are validator key or genesis account key. The first key is always validator key which is madatory, but genesis account keys that are not first info are optional. `(M)`
      - `name`: Key name. `(M)`
      - `keyring_backend`: Keyring backend type. Select only 'file' or 'test.`(M)`
      - `balance`: Balance amount of validator or genesis account. In case of genesis account, it is optional `(M)`
    - `keys_option`: The option of all keys in the validator node. `(O)`
      - `not_save_mnemonic`: The option that do not save the key_seed file. `(O, default false)`
      - `print_mnemonic`: The option that print mnemonic when generate key. `(O, default false)`
    - `sentries`: The option to create sentry nodes of the validator. These are same configuration of the validator except for persisten peer list and etc. If the validator has not this option, xpla.set only creates the validator. `(O)`
      - `ip_address`: IP address of the each sentry node. `(O, default "")`

### Gen files
Run xpla.set.

```shell
./start
```

### Run xpla node
After run xpla.set, required files to run the validator are created in the home directory. Copy created directories to root directory of each validator. The tree as below indicates the configuration of the validator's files. If sentry node option is exist, config files of the sentry node such as `sentry0` and `sentry1` should be moved to the root directory of each sentry node.

```bash
.xplaset
├── gentxs
│   ├── validator_node_0.json
│   └── validator_node_1.json
├── validator0
│   ├── config
│   │   ├── app.toml
│   │   ├── config.toml
│   │   ├── genesis.json
│   │   ├── gentx
│   │   │   └── gentx-f94b6f0cebc48468cfef95084bdbf1d12181ba44.json
│   │   ├── node_key.json
│   │   └── priv_validator_key.json
│   ├── data
│   │   └── priv_validator_state.json
│   ├── keyring-file
│   │   ├── 2d419fc05ab8486e496627b88b9e8fd0d4dac0f3.address
│   │   ├── 34585d96409d8af38b2db5894b170f03ec6ace76.address
│   │   ├── 8c64277064884e40b848911f3f0130d9645ee68b.address
│   │   ├── keyhash
│   │   ├── myKey1.info
│   │   ├── myKey1_mnemonic.json
│   │   ├── myKey2.info
│   │   ├── myKey2_mnemonic.json
│   │   ├── myKey3.info
│   │   └── myKey3_mnemonic.json
│   ├── sentry0
│   │   ├── config
│   │   │   ├── app.toml
│   │   │   ├── config.toml
│   │   │   ├── genesis.json
│   │   │   ├── node_key.json
│   │   │   └── priv_validator_key.json
│   │   └── data
│   │       └── priv_validator_state.json
│   └── sentry1
│       ├── config
│       │   ├── app.toml
│       │   ├── config.toml
│       │   ├── genesis.json
│       │   ├── node_key.json
│       │   └── priv_validator_key.json
│       └── data
│           └── priv_validator_state.json
└── validator1
    ├── ...
```
Then, run `xplad`.
```shell
xplad start
```


