package gira

type SdkAccount struct {
	UnionId         string
	AccessToken     string
	RefreshToken    string
	TokenExpireTime int64
	NickName        string
	Gender          int32
	SmallPic        string
	LargePic        string
}

type Sdk interface {
	SdkLogin(accountPlat string, openId string, token string) (*SdkAccount, error)
}
