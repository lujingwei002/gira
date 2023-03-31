package game

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/sproto"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/lujingwei002/gira/app"
	"google.golang.org/grpc"
)

// 玩家接口
type Player interface {
	Load(ctx context.Context, memberId string, userId string) error
	Logout(ctx context.Context) error
}

// 大厅
type Hall interface {
	// 将消息推送给userId
	Push(ctx context.Context, userId string, req sproto.SprotoPush) error
}

type Session interface {
	//  将消息推送给当前玩家
	Push(resp sproto.SprotoPush) (err error)
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

type HallApplication struct {
	app.BaseFacade
	hall *hall

	Proto         *sproto.Sproto
	PlayerHandler *sproto.SprotoHandler
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
		app.hall = newHall(app.Proto, handler, app.PlayerHandler)
		app.Hall = app.hall
		return nil
	}
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

func (self *HallApplication) Push(ctx context.Context, userId string, req sproto.SprotoPush) error {
	return self.hall.Push(ctx, userId, req)
}
