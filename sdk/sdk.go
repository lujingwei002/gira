package sdk

import (
	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
)

type sdk_interface interface {
	Login(accountPlat string, open_id string, token string) (*gira.SdkAccount, error)
}

type Sdk struct {
	TestSdk *TestSdk
	sdkDict map[string]sdk_interface
}

func NewSdk() *Sdk {
	self := &Sdk{}
	return self
}

func NewConfigSdk(config gira.SdkConfig) *Sdk {
	self := &Sdk{
		sdkDict: make(map[string]sdk_interface),
	}
	log.Info(config)
	if config.Test != nil {
		self.TestSdk = ConfigTestSdk(*config.Test)
		self.sdkDict["test"] = self.TestSdk
	}
	return self
}

func (self *Sdk) Login(accountPlat string, open_id string, token string) (*gira.SdkAccount, error) {
	if sdk, ok := self.sdkDict[accountPlat]; !ok {
		return nil, gira.ErrSdkNotImplement
	} else {
		return sdk.Login(accountPlat, open_id, token)
	}
}

type TestSdk struct {
	config gira.TestSdkConfig
}

func ConfigTestSdk(config gira.TestSdkConfig) *TestSdk {
	self := &TestSdk{
		config: config,
	}
	return self
}

func (self *TestSdk) Login(accountPlat string, open_id string, token string) (*gira.SdkAccount, error) {
	if token != self.config.Secret {
		return nil, gira.ErrInvalidSdkToken
	}
	result := &gira.SdkAccount{
		Nickname: open_id,
	}
	return result, nil
}
