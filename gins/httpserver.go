package gins

/// 参考 https://juejin.cn/post/6844903833273892871

import (
	"context"
	"net/http"
	"time"

	"github.com/lujingwei002/gira/log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lujingwei002/gira"
)

type BaseJsonResponse struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

func (self BaseJsonResponse) SetCode(v int32) {
	self.Code = v
}

func (self BaseJsonResponse) SetMsg(v string) {
	self.Msg = v
}

type JsonResponse interface {
	SetCode(v int32)
	SetMsg(v string)
}

func HttpJsonResponse(g *gin.Context, httpCode int, err error, data JsonResponse) {
	errorCode := gira.ErrCode(err)
	errorMsg := gira.ErrMsg(err)
	if data == nil {
		resp := BaseJsonResponse{
			Code: errorCode,
			Msg:  errorMsg,
		}
		g.JSON(httpCode, resp)
	} else {
		data.SetCode(errorCode)
		data.SetMsg(errorMsg)
		g.JSON(httpCode, data)
	}
}

// 获取远程客户端ip要设置SetTrustedProxies
// 参考 https://www.cnblogs.com/mayanan/p/15703234.html

func JWT() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		var claims *Claims
		token := g.GetHeader("Authorization")
		if token == "" {
			err = gira.ErrInvalidJwt
		} else {
			//TODO
			claims, err = ParseJwtToken1(token, "app.Config.Jwt.Secret")
			if err != nil {
				err = gira.ErrInvalidJwt
			} else if time.Now().Unix() > claims.ExpiresAt {
				err = gira.ErrJwtExpire
			} else if claims.MemberId == "" {
				err = gira.ErrInvalidJwt
			} else {
				g.Set("MemberId", claims.MemberId)
				g.Set("Claims", claims)
			}
		}
		if err != nil {
			HttpJsonResponse(g, http.StatusUnauthorized, err, nil)
			g.Abort()
			return
		}
		g.Next()
	}
}

type HttpServer struct {
	Config     gira.HttpConfig
	Handler    http.Handler
	server     *http.Server
	cancelCtx  context.Context
	cancelFunc context.CancelFunc
}

func NewConfigHttpServer(facade gira.Application, config gira.HttpConfig, router http.Handler) (*HttpServer, error) {
	h := &HttpServer{
		Config:  config,
		Handler: router,
	}
	h.cancelCtx, h.cancelFunc = context.WithCancel(facade.Context())
	server := &http.Server{
		Addr:           h.Config.Addr,
		Handler:        h.Handler,
		ReadTimeout:    time.Duration(h.Config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(h.Config.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	h.server = server
	httpFunc := func() error {
		if config.Ssl && len(config.CertFile) > 0 && len(config.KeyFile) > 0 {
			if err := server.ListenAndServeTLS(config.CertFile, config.KeyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("https server shutdown, err: %s\n", err)
				return err
			} else {
				log.Infof("https server shutdown")
			}
		} else {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Infof("http server shutdown, err: %s\n", err)
				return err
			} else {
				log.Infof("http server shutdown")
			}
		}
		return nil
	}
	httpCtl := func() error {
		select {
		case <-facade.Done():
			log.Info("http server recv down")
			h.server.Close()
		}
		return nil
	}
	facade.Go(httpFunc)
	facade.Go(httpCtl)
	log.Infof("http server started, addr=%s\n", h.Config.Addr)
	return h, nil
}

type Claims struct {
	jwt.StandardClaims
	MemberId    string `json:"member_id"`
	AccountPlat string `json:"account_plat"`
	PhonePlat   int32  `json:"phone_plat"`
	DeviceId    string `json:"device_id"`
	Channel     int64  `json:"channel"`
}

func GenerateJwtToken1(secret string, duration time.Duration, memberId string, channel int64, accountPlat string, phonePlat int32, deviceID string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(duration)
	claims := Claims{
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gira",
		},
		memberId,
		accountPlat,
		phonePlat,
		deviceID,
		channel,
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(secret))
	return token, err
}

func GenerateJwtRefreshToken1(secret string, duration time.Duration) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(duration)
	claims := Claims{
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gira",
		},
		"", // MemberId
		"", // AccountPlat
		0,  // PhonePlat
		"", // DeviceId
		0,  // Channel
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(secret))
	return token, err
}

func ParseJwtToken1(token string, secret string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
