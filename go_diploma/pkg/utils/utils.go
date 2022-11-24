package utils

import (
	"go_diploma/pkg/config"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

func LoadSettings() *config.Config {
	conf := &config.Config{}
	if _, err := toml.DecodeFile("../settings.toml", conf); err != nil {
		panic(err)
	}
	conf.LoadCoutryCodes()
	conf.SetCoutryCodes(Sorting(conf.GetCoutryCodes()))
	conf.SetEmailProviders(Sorting(conf.GetEmailProviders()))
	conf.SetSmsProviders(Sorting(conf.GetSmsProviders()))
	conf.SetVoiceProviders(Sorting(conf.GetVoiceProviders()))
	if len(strings.Trim(conf.Emu.Path, " ")) == 0 {
		panic("Please fulfill field EmuPath in settings.toml!")
	}
	return conf
}

func SliceContains(s []string, searchterm string) bool {
	/*s - sorted slice*/
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

func Sorting(arr []string) []string {
	bsSize := len(arr)
	for i := 1; i < bsSize; i++ {
		x := arr[i]
		j := i
		for ; j > 0 && arr[j-1] > x; j-- {
			arr[j] = arr[j-1]
		}
		arr[j] = x
	}
	return arr
}
