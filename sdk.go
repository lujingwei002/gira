package gira

type SdkAccount struct {
	UnionId         string
	AccessToken     string // access token
	RefreshToken    string // refresh token
	TokenExpireTime int64  // token过期时间
	NickName        string // 昵称
	Gender          int32  // 性别
	SmallPic        string // 小头像地址
	LargePic        string // 大头像地址
}

type SdkComponent interface {
	// 登录
	Login(accountPlat string, openId string, token string) (*SdkAccount, error)
}
