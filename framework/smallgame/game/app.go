package game

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/common/rpc"
	"github.com/lujingwei002/gira/log"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/lujingwei002/gira/app"
	"google.golang.org/grpc"
)

// 玩家接口
type Player interface {
	Load(ctx context.Context, memberId string, userId string) error
	Logout(ctx context.Context) error
	Update()
	Save(ctx context.Context) error
}

// 大厅
type Hall interface {
	// 将消息推送给玩家, 不保证消息已经被处理，如果玩家当前不在线，消息将会推送失败，但不会返回错误
	Push(ctx context.Context, userId string, req gira.ProtoPush) error
	// 会将消息推送到玩家的消息队列中，但不等待结果，如果玩家不在线，会返回错误
	MustPush(ctx context.Context, userId string, resp gira.ProtoPush) (err error)
}

type Session interface {
	Push(resp gira.ProtoPush) (err error)
}

type HallHandler interface {
	// 由memberId创建账号
	NewUser(ctx context.Context, memberId string) (avatar UserAvatar, err error)
	// 创建player
	NewPlayer(ctx context.Context, session Session, memberId string, avatar UserAvatar) (player Player, err error)
}

type UserAvatar interface {
	GetUserId() primitive.ObjectID
}

// 应用
type HallApplication struct {
	app.BaseFacade
	hall   *hall
	Config *Config

	Proto         gira.Proto
	PlayerHandler gira.ProtoHandler
	Hall          Hall
}

func (self *HallApplication) Create() error {
	log.Info("create")
	return nil
}

func (app *HallApplication) OnFrameworkAwake(facade gira.ApplicationFacade) error {
	if handler, ok := facade.(HallHandler); !ok {
		return gira.ErrGateHandlerNotImplement
	} else {
		hall := newHall(facade, app.Proto, app.Config, handler, app.PlayerHandler)
		app.hall = hall
		app.Hall = hall
		rpc.OnAwake(facade)
		return nil
	}
}

func (self *HallApplication) OnFrameworkConfigLoad(c *gira.Config) error {
	self.Config = &Config{}
	return self.Config.OnConfigLoad(c)
}

func (self *HallApplication) OnGrpcServerStart(server *grpc.Server) error {
	return nil
}

func (self *HallApplication) OnFrameworkGrpcServerStart(server *grpc.Server) error {
	if err := self.hall.Register(server); err != nil {
		return err
	}
	return nil
}
