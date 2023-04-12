package util

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Moonyongjung/xpla-set/types"
	"github.com/logrusorgru/aurora"
)

func LogInfo(log ...interface{}) {
	print := logTime() + G(" XPLA-SET   ") + ToStringTrim(log, "")
	fmt.Println(print)
}

func LogWarning(log ...interface{}) {
	print := logTime() + " " + BgR("WARNING") + "    " + ToStringTrim(log, "")
	fmt.Println(print)
}

func LogKV(key string, value string) {
	LogInfo(BB(key+"=") + value)
}

func LogErr(errType types.XGoError, errDesc ...interface{}) error {
	print := logErr("code", errType.Code(), ":", errType.Desc(), "-", errDesc)
	fmt.Println(print)

	return errors.New(ToStringTrim(errDesc, ""))
}

func logErr(log ...interface{}) string {
	return logTime() + R(" Error      ") + ToStringTrim(log, "")
}

func logTime() string {
	return B(time.Now().Format("2006-01-02 15:04:05"))
}

func G(str string) string {
	return aurora.Green(str).String()
}

func R(str string) string {
	return aurora.Red(str).String()
}

func BB(str string) string {
	return aurora.BrightBlue(str).String()
}

func B(str string) string {
	return aurora.Blue(str).String()
}

func Y(str string) string {
	return aurora.Yellow(str).String()
}

func BgR(str string) string {
	return aurora.BgRed(str).String()
}

func BgG(str string) string {
	return aurora.BgGreen(str).String()
}

func W(str string) string {
	return aurora.White(str).String()
}

func ToStringTrim(value interface{}, defaultValue string) string {
	s := fmt.Sprintf("%v", value)
	s = s[1 : len(s)-1]
	str := strings.TrimSpace(s)
	if str == "" {
		return defaultValue
	} else {
		return str
	}
}
