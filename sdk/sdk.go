package sdk

import (
	"context"
	"fmt"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/sdk/ultra"
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
	if config.Ultra != nil {
		self.ultraSdk = NewConfigUltraSdk(*config.Ultra)
		self.sdkDict["ultra"] = self.ultraSdk
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

func NewConfigUltraSdk(config gira.UltraSdkConfig) *UltraSdk {
	self := &UltraSdk{
		config: config,
	}
	return self
}

// 服务端sdk接口
type sdk_server interface {
	Login(accountPlat string, openId string, token string, authUrl string, appId string, appSecret string) (*gira.SdkAccount, error)
	PayOrderCheck(accountPlat string, args map[string]interface{}, paySecret string) (*gira.SdkPayOrder, error)
}

type SdkComponent struct {
	testSdk  *TestSdk
	pwdSdk   *PwdSdk
	ultraSdk *UltraSdk
	sdkDict  map[string]sdk_server
}

func (self *SdkComponent) Login(accountPlat string, openId string, token string, authUrl string, appId string, appSecret string) (*gira.SdkAccount, error) {
	if sdk, ok := self.sdkDict[accountPlat]; !ok {
		return nil, gira.ErrSdkComponentNotImplement.Trace()
	} else {
		return sdk.Login(accountPlat, openId, token, authUrl, appId, appSecret)
	}
}

func (self *SdkComponent) PayOrderCheck(accountPlat string, args map[string]interface{}, paySecret string) (*gira.SdkPayOrder, error) {
	if sdk, ok := self.sdkDict[accountPlat]; !ok {
		return nil, gira.ErrSdkComponentNotImplement.Trace()
	} else {
		return sdk.PayOrderCheck(accountPlat, args, paySecret)
	}
}

type TestSdk struct {
	config gira.TestSdkConfig
}

func (self *TestSdk) Login(accountPlat string, openId string, token string, authUrl string, appId string, appSecret string) (*gira.SdkAccount, error) {
	if token != appSecret {
		return nil, gira.ErrInvalidSdkToken.Trace()
	}
	result := &gira.SdkAccount{
		OpenId:   openId,
		NickName: openId,
	}
	return result, nil
}

func (self *TestSdk) PayOrderCheck(accountPlat string, args map[string]interface{}, paySecret string) (*gira.SdkPayOrder, error) {
	if v, ok := args["cporder_id"]; !ok {
		return nil, gira.ErrSdkPayOrderCheckArgsInvalid
	} else if cporderId, ok := v.(string); !ok {
		return nil, gira.ErrSdkPayOrderCheckArgsInvalid
	} else if v, ok := args["amount"]; !ok {
		return nil, gira.ErrSdkPayOrderCheckArgsInvalid
	} else if amount, ok := v.(int64); !ok {
		return nil, gira.ErrSdkPayOrderCheckArgsInvalid
	} else {
		result := &gira.SdkPayOrder{
			CporderId: cporderId,
			OrderId:   cporderId,
			// PayTime: payInfo.PayTime,
			Amount:   amount,
			Response: "SUCCESS",
		}
		return result, nil
	}
}

type PwdSdk struct {
	config gira.PwdSdkConfig
}

func (self *PwdSdk) Login(accountPlat string, openId string, token string, authUrl string, appId string, appSecret string) (*gira.SdkAccount, error) {
	result := &gira.SdkAccount{
		OpenId:      openId,
		NickName:    openId,
		AccessToken: token, // 登录密码
	}
	return result, nil
}

func (self *PwdSdk) PayOrderCheck(accountPlat string, args map[string]interface{}, paySecret string) (*gira.SdkPayOrder, error) {
	return nil, gira.ErrSdkPayOrderCheckMethodNotImplement
}

type UltraSdk struct {
	config   gira.UltraSdkConfig
	loginSdk *ultra.USDK
}

func (self *UltraSdk) Login(accountPlat string, openId string, token string, authUrl string, appId string, appSecret string) (*gira.SdkAccount, error) {
	log.Infow("ultra sdk login", "open_id", openId, "token", token, "app_id", appId, "app_secret", appSecret, "auth_url", authUrl)
	if self.loginSdk == nil {
		if sdk, err := ultra.NewUSDK(context.Background(), authUrl, appId, appSecret, "", ultra.AreaGlobal); err != nil {
			return nil, err
		} else {
			self.loginSdk = sdk
		}
	}
	if _, resp, err := self.loginSdk.Login(openId, "", token); err != nil {
		return nil, err
	} else if resp.Code != 0 {
		return nil, gira.NewError(int32(resp.Code), resp.Msg)
	} else {
		result := &gira.SdkAccount{
			OpenId:   fmt.Sprintf("%v_%v", resp.ChannelId, resp.CUid),
			NickName: resp.CName,
		}
		return result, nil
	}
}

func (self *UltraSdk) PayOrderCheck(accountPlat string, args map[string]interface{}, paySecret string) (*gira.SdkPayOrder, error) {
	if v, ok := args["data"]; !ok {
		return nil, gira.ErrSdkPayOrderCheckArgsInvalid
	} else if data, ok := v.(string); !ok {
		return nil, gira.ErrSdkPayOrderCheckArgsInvalid
	} else if v, ok := args["sign"]; !ok {
		return nil, gira.ErrSdkPayOrderCheckArgsInvalid
	} else if sign, ok := v.(string); !ok {
		return nil, gira.ErrSdkPayOrderCheckArgsInvalid
	} else if sdk, err := ultra.NewUSDK(context.Background(), "", "", "", paySecret, ultra.AreaGlobal); err != nil {
		return nil, err
	} else if response, _, payInfo, err := sdk.PayCallback(data, sign); err != nil {
		return nil, err
	} else {
		log.Printf("支付回调 gameOrder: %v\t游戏的订单号", payInfo.GameOrder)
		log.Printf("支付回调 selfDefine: %v\t游戏的自定义信息", payInfo.SelfDefine)
		log.Printf("支付回调 amount: %v\t付款金额", payInfo.Amount)
		log.Printf("支付回调 payTime: %v\t支付时间", payInfo.PayTime)
		log.Printf("支付回调 goodsId: %v\t支付的道具ID", payInfo.GoodsId)
		log.Printf("支付回调 channelUid: %v\t渠道用户ID", payInfo.ChannelUid)
		log.Printf("支付回调 channel: %v\t渠道码", payInfo.Channel)
		log.Printf("支付回调 channelId: %v\t渠道ID", payInfo.ChannelId)
		log.Printf("支付回调 orderNo: %v\tUSDK的订单号", payInfo.OrderNo)
		result := &gira.SdkPayOrder{
			CporderId: payInfo.GameOrder,
			OrderId:   payInfo.OrderNo,
			// PayTime: payInfo.PayTime,
			Amount:   int64(payInfo.Amount),
			Response: response,
		}
		return result, nil
	}
}
