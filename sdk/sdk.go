package sdk

import (
	"github.com/lujingwei002/gira"
)

func NewConfigSdk(config gira.SdkConfig) *SdkComponent {
	self := &SdkComponent{
		sdkDict: make(map[string]sdk_server),
	}
	if config.Test != nil {
		self.testSdk = NewConfigTestSdk(*config.Test)
		self.sdkDict["test"] = self.testSdk
	}
	if config.Pwd != nil {
		self.pwdSdk = NewConfigGfSdk(*config.Pwd)
		self.sdkDict["pwd"] = self.pwdSdk
	}
	return self
}

func NewConfigTestSdk(config gira.TestSdkConfig) *TestSdk {
	self := &TestSdk{
		config: config,
	}
	return self
}

func NewConfigGfSdk(config gira.PwdSdkConfig) *PwdSdk {
	self := &PwdSdk{
		config: config,
	}
	return self
}

// 服务端sdk接口
type sdk_server interface {
	Login(accountPlat string, openId string, token string) (*gira.SdkAccount, error)
}

type SdkComponent struct {
	testSdk *TestSdk
	pwdSdk  *PwdSdk
	sdkDict map[string]sdk_server
}

func (self *SdkComponent) Login(accountPlat string, openId string, token string) (*gira.SdkAccount, error) {
	if sdk, ok := self.sdkDict[accountPlat]; !ok {
		return nil, gira.ErrSdkComponentNotImplement.Trace()
	} else {
		return sdk.Login(accountPlat, openId, token)
	}
}

type TestSdk struct {
	config gira.TestSdkConfig
}

func (self *TestSdk) Login(accountPlat string, openId string, token string) (*gira.SdkAccount, error) {
	if token != self.config.Secret {
		return nil, gira.ErrInvalidSdkToken.Trace()
	}
	result := &gira.SdkAccount{
		NickName: openId,
	}
	return result, nil
}

type PwdSdk struct {
	config gira.PwdSdkConfig
}

func (self *PwdSdk) Login(accountPlat string, openId string, token string) (*gira.SdkAccount, error) {
	result := &gira.SdkAccount{
		NickName:    openId,
		AccessToken: token,
	}
	return result, nil
}
