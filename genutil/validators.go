package genutil

import (
	"path"

	"github.com/Moonyongjung/xpla.go/client"
	xutil "github.com/Moonyongjung/xpla.go/util"
	"github.com/Moonyongjung/xpla.set/types"
	"github.com/Moonyongjung/xpla.set/util"
)

// Setting validator node.
// Handle all processes in order to create validator node such as creating key, initializing config file and etc.
func Validators(validators []Validator, home string, xplac *client.XplaClient) error {
	for i, validator := range ConfigFile().Get().XplaGen.Validators {
		valNodeName := "validator" + xutil.FromIntToString(i)
		valPath := path.Join(home, valNodeName)
		util.GuideCreateValidator(valNodeName)

		if validator.Moniker == "" {
			validator.Moniker = valNodeName
		}

		util.LogInfo(util.B("→ start create config files..."))
		// init config/app files
		err := SetInit(valPath)
		if err != nil {
			return err
		}

		// ready to generate sentry nodes if sentry infos are included in the config.yaml file
		if len(validator.Sentries.IpAddress) != 0 {
			for i := 0; i < len(validator.Sentries.IpAddress); i++ {
				sentryPath := path.Join(valPath, "sentry"+xutil.FromIntToString(i))
				err = util.GenFilePath(sentryPath)
				if err != nil {
					return util.LogErr(types.ErrInvalidRequest, err)
				}

				err = SetInit(sentryPath)
				if err != nil {
					return util.LogErr(types.ErrInvalidRequest, err)
				}

			}
		}
		util.LogInfo(util.B("← success create config files of the " + validator.Moniker))

		// create keys are validator keys and genesis accounts
		util.LogInfo(util.B("→ start create keys..."))
		keyring, err := KeysAdd(validator, valPath, valNodeName)
		if err != nil {
			return err
		}
		util.LogInfo(util.B("← success create keys of the " + validator.Moniker))

		types.GentxsDirPath = path.Join(home, "gentxs")
		err = util.GenFilePath(types.GentxsDirPath)
		if err != nil {
			return err
		}

		// set node
		err = Node(validator, valPath, i, keyring, xplac)
		if err != nil {
			return err
		}

		util.LogInfo(util.G("success node setting and create gentx of the " + validator.Moniker))
	}

	return nil
}
