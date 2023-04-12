package util

import (
	"fmt"
	"net"
	"os"

	"github.com/Moonyongjung/xpla-set/types"
	xutil "github.com/Moonyongjung/xpla.go/util"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetConfirm() bool {
	for {
		var s string
		fmt.Scan(&s)

		if s == "y" {
			return true
		} else if s == "N" {
			return false
		} else {
			LogErr(types.ErrInvalidRequest, "Input correct string [y/N]")
			LogWarning("try again")
		}
	}
}

func IsExistPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func GenFilePath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func GetIP(startingIPAddr string) (ip string, err error) {
	if len(startingIPAddr) == 0 {
		ip, err = server.ExternalIP()
		if err != nil {
			return "", err
		}
		return ip, nil
	}
	return calculateIP(startingIPAddr)
}

func calculateIP(ip string) (string, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return "", fmt.Errorf("%v: non ipv4 address", ip)
	}

	return ipv4.String(), nil
}

func ConvSdkInt(value string) (sdk.Int, error) {
	bigInt, err := xutil.FromStringToBigInt(value)
	if err != nil {
		return sdk.NewInt(0), LogErr(types.ErrInvalidRequest, err)
	}

	return sdk.NewIntFromBigInt(bigInt), nil
}
