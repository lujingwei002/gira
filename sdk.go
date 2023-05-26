package gira

type SdkAccount struct {
	UnionId         string
	OpenId          string
	AccessToken     string // access token
	RefreshToken    string // refresh token
	TokenExpireTime int64  // token过期时间
	NickName        string // 昵称
	Gender          int32  // 性别
	SmallPic        string // 小头像地址
	LargePic        string // 大头像地址
}

type SdkPayOrder struct {
	CporderId  string
	OrderId    string
	PayTime    int64
	Amount     int64 // 单位 分
	OpenId     string
	UserDefine string
}
type SdkComponent interface {
	// 登录
	Login(accountPlat string, openId string, token string, authUrl string, appId string, appSecret string) (*SdkAccount, error)
	PayOrderCheck(accountPlat string, data string, sign string, paySecret string) (*SdkPayOrder, error)
}
