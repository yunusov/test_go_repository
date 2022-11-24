package config

import (
	"fmt"
)

type Config struct {
	Emu struct {
		CountryCodes    []string
		CountryCodesMap map[string]string
		Path            string
	}
	Datafiles struct {
		Billing string
		Email   string
		Sms     string
		Voice   string
	}
	Providers struct {
		Email []string
		Sms   []string
		Voice []string
	}
	ServerAddresses struct {
		Incindent string
		Mms       string
		Service   string
		Support   string
	}
}

func (conf *Config) ToString() string {
	return fmt.Sprintf("config = %v", conf)
}

func (conf *Config) GetCoutryCodes() []string {
	return conf.Emu.CountryCodes
}

func (conf *Config) GetCoutryCodesMap() map[string]string {
	return conf.Emu.CountryCodesMap
}

func (conf *Config) SetCoutryCodes(countryCodes []string) {
	conf.Emu.CountryCodes = countryCodes
}

func (conf *Config) LoadCoutryCodes() {
	ccMap := conf.GetCoutryCodesMap()
	keys := make([]string, 0, len(ccMap))
	for k := range ccMap {
		keys = append(keys, k)
	}
	conf.SetCoutryCodes(keys)
}

func (conf *Config) GetEmailProviders() []string {
	return conf.Providers.Email
}

func (conf *Config) SetEmailProviders(emailProviders []string) {
	conf.Providers.Email = emailProviders
}

func (conf *Config) GetSmsProviders() []string {
	return conf.Providers.Sms
}

func (conf *Config) SetSmsProviders(smsProviders []string) {
	conf.Providers.Sms = smsProviders
}

func (conf *Config) GetVoiceProviders() []string {
	return conf.Providers.Voice
}

func (conf *Config) SetVoiceProviders(voiceCodes []string) {
	conf.Providers.Voice = voiceCodes
}

func (conf *Config) GetServiceServerAddress() string {
	return conf.ServerAddresses.Service
}

func (conf *Config) GetIncindentServerAddress() string {
	return conf.ServerAddresses.Incindent
}

func (conf *Config) GetEmuPath() string {
	return conf.Emu.Path
}

func (conf *Config) GetBillingDataFile() string {
	return conf.Datafiles.Billing
}

func (conf *Config) GetEmailDataFile() string {
	return conf.Datafiles.Email
}

func (conf *Config) GetSmsDataFile() string {
	return conf.Datafiles.Sms
}

func (conf *Config) GetVoiceDataFile() string {
	return conf.Datafiles.Voice
}

func (conf *Config) GetMmsServerAddress() string {
	return conf.ServerAddresses.Mms
}

func (conf *Config) GetSupportServerAddress() string {
	return conf.ServerAddresses.Support
}
