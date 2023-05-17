package gins

/// 参考 https://juejin.cn/post/6844903833273892871

import (
	"context"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/log"
)

type BaseJsonResponse struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

func (self *BaseJsonResponse) SetCode(v int32) {
	self.Code = v
}

func (self *BaseJsonResponse) SetMsg(v string) {
	self.Msg = v
}

type JsonResponse interface {
	SetCode(v int32)
	SetMsg(v string)
}

// 返回json response
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
	config     gira.HttpConfig
	Handler    http.Handler
	server     *http.Server
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewConfigHttpServer(ctx context.Context, config gira.HttpConfig, router http.Handler) (*HttpServer, error) {
	cancelCtx, cancelFunc := context.WithCancel(ctx)

	h := &HttpServer{
		config:     config,
		Handler:    router,
		ctx:        cancelCtx,
		cancelFunc: cancelFunc,
	}

	server := &http.Server{
		Addr:           h.config.Addr,
		Handler:        h.Handler,
		ReadTimeout:    time.Duration(h.config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(h.config.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	h.server = server
	return h, nil
}

func (h *HttpServer) Start() error {
	facade.Go(h.serve)
	facade.Go(func() error {
		select {
		case <-h.ctx.Done():
			h.server.Close()
		}
		return nil
	})
	return nil
}

func (h *HttpServer) serve() error {
	log.Debugw("http server started", "addr", h.config.Addr)
	if h.config.Ssl && len(h.config.CertFile) > 0 && len(h.config.KeyFile) > 0 {
		if err := h.server.ListenAndServeTLS(h.config.CertFile, h.config.KeyFile); err != nil && err != http.ErrServerClosed {
			return err
		}
	} else {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
	}
	h.onStop()
	return nil
}

func (h *HttpServer) onStop() {
	log.Debugw("http server on stop")
}

func (h *HttpServer) Stop() error {
	log.Debugw("http server stop")
	h.cancelFunc()
	return nil
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
