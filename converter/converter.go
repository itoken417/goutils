package converter

import (
	"github.com/ktnyt/go-moji"
)

func HK2ZK(str string) string {
	str = moji.Convert(str, moji.HK, moji.ZK)
	return str
}
func ZE2HE(str string) string {
	str = moji.Convert(str, moji.ZE, moji.HE)
	return str
}
