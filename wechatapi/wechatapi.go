package wechatapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lujingwei002/gira"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

type JsCode2SessionResponse struct {
	ErrCode    int32  `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	OpenId     string `json:"openid"`
}

const (
	JsCode2SessionGrantType_AuthorizationCode = "authorization_code "
)

// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
// 小程序登录
func JsCode2Session(appId string, secret string, jsCode string, grantType string) (*JsCode2SessionResponse, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	host := "https://api.weixin.qq.com"
	params := url.Values{}
	params.Set("appid", appId)
	params.Set("secret", secret)
	params.Set("js_code", jsCode)
	params.Set("grant_type", grantType)
	url := fmt.Sprintf("%s/sns/jscode2session?%s", host, params.Encode())
	var result *http.Response
	var err error
	result, err = client.Get(url)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	var body []byte
	body, err = io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	resp := &JsCode2SessionResponse{}
	if err = json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type GetAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

const (
	GetAccessTokenGrantType_AuthorizationCode = "client_credential"
)

// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-access-token/getAccessToken.html
// 小程序登录
func GetAccessToken(appId string, secret string, grantType string) (*GetAccessTokenResponse, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	host := "https://api.weixin.qq.com"
	params := url.Values{}
	params.Set("appid", appId)
	params.Set("secret", secret)
	params.Set("grant_type", grantType)
	url := fmt.Sprintf("%s/cgi-bin/token?%s", host, params.Encode())
	var result *http.Response
	var err error
	result, err = client.Get(url)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	var body []byte
	body, err = io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	resp := &GetAccessTokenResponse{}
	if err = json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	log.Println("f", resp)
	return resp, nil
}

type GetWxaCodeUnLimitResponse struct {
	ErrCode     int32  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	Buffer      []byte
	ContentType string
}

type GetWxaCodeUnLimitRequest struct {
	Scene      string `json:"scene"`
	Page       string `json:"page"`
	EnvVersion string `json:"env_version"`
}

// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/qr-code/getUnlimitedQRCode.html
// 小程序登录
func GetWxaCodeUnLimit(accessToken string, scene string, page string, envVersion string) (*GetWxaCodeUnLimitResponse, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	host := "https://api.weixin.qq.com"
	params := url.Values{}
	params.Set("access_token", accessToken)
	url := fmt.Sprintf("%s/wxa/getwxacodeunlimit?%s", host, params.Encode())
	req := &GetWxaCodeUnLimitRequest{
		Scene:      scene,
		Page:       page,
		EnvVersion: envVersion,
	}
	var httpReq *http.Request
	var result *http.Response
	var data []byte
	var err error
	data, err = json.Marshal(req)
	if err != nil {
		return nil, err
	}
	log.Println(string(data))
	httpReq, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	client := &http.Client{Transport: tr}
	// Send the HTTP request
	result, err = client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	var body []byte
	body, err = io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	resp := &GetWxaCodeUnLimitResponse{
		ContentType: result.Header.Get("Content-Type"),
	}
	log.Println(result.Header)
	log.Println(resp.ContentType)
	if strings.HasPrefix(resp.ContentType, "image/jpeg") {
		resp.Buffer = body
		return resp, nil
	} else if strings.HasPrefix(resp.ContentType, "application/json") {
		// log.Println(string(body))
		if err = json.Unmarshal(body, resp); err != nil {
			return nil, err
		}
		return resp, nil
	} else {
		return nil, gira.ErrTodo
	}
}

type StsGetCredentialResponse struct {
	ErrCode      int32  `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	TmpSecretID  string `json:"TmpSecretId"`
	TmpSecretKey string `json:"TmpSecretKey"`
	SessionToken string `json:"Token"`
}

// https://cloud.tencent.com/document/product/436/14048
// 获取临时密钥
// region ap-guangzhou
func StsGetCredential(appId string, secretId string, secretKey string, bucket string, region string) (*StsGetCredentialResponse, error) {
	response := &StsGetCredentialResponse{}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	c := sts.NewClient(
		secretId,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考https://cloud.tencent.com/document/product/598/37140
		secretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考https://cloud.tencent.com/document/product/598/37140
		client,
		// sts.Host("sts.internal.tencentcloudapi.com"), // 设置域名, 默认域名sts.tencentcloudapi.com
		// sts.Scheme("http"),      // 设置协议, 默认为https，公有云sts获取临时密钥不允许走http，特殊场景才需要设置http
	)
	// 策略概述 https://cloud.tencent.com/document/product/436/18023
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					// 密钥的权限列表。简单上传和分片需要以下的权限，其他权限列表请看 https://cloud.tencent.com/document/product/436/31923
					Action: []string{
						// 简单上传
						"name/cos:PostObject",
						"name/cos:PutObject",
						// 分片上传
						"name/cos:InitiateMultipartUpload",
						"name/cos:ListMultipartUploads",
						"name/cos:ListParts",
						"name/cos:UploadPart",
						"name/cos:CompleteMultipartUpload",
					},
					Effect: "allow",
					Resource: []string{
						"*",
						// 这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						// 存储桶的命名格式为 BucketName-APPID，此处填写的 bucket 必须为此格式
						//	"qcs::cos:ap-guangzhou:uid/" + appId + ":" + bucket + "/exampleobject",
					},
					// 开始构建生效条件 condition
					// 关于 condition 的详细设置规则和COS支持的condition类型可以参考https://cloud.tencent.com/document/product/436/71306
					// Condition: map[string]map[string]interface{}{
					// 	"ip_equal": map[string]interface{}{
					// 		"qcs:ip": []string{
					// 			"10.217.182.3/24",
					// 			"111.21.33.72/24",
					// 		},
					// 	},
					// },
				},
			},
		},
	}
	res, err := c.GetCredential(opt)
	if err != nil {
		return nil, err
	}
	response.TmpSecretID = res.Credentials.TmpSecretID
	response.SessionToken = res.Credentials.SessionToken
	response.TmpSecretKey = res.Credentials.TmpSecretKey
	return response, nil
}
