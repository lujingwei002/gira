package hall

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/actor"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/options/registry_options"
	"github.com/lujingwei002/gira/service/hall/hall_grpc"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

// 玩家接口
type Player interface {
	// 加载数据
	Load(ctx context.Context, memberId string, userId string) error
	// 断开连接,登出
	Logout(ctx context.Context) error
	// 保存数据
	Save(ctx context.Context) error
	Update()
}

// 大厅接口
type Hall interface {
	SessionCount() int64
	// 将消息推送给玩家, 不保证消息已经被处理，如果玩家当前不在线，消息将会推送失败，但不会返回错误
	Push(ctx context.Context, userId string, req gira.ProtoPush) error
	// 会将消息推送到玩家的消息队列中，但不等待结果，如果玩家不在线，会返回错误
	MustPush(ctx context.Context, userId string, resp gira.ProtoPush) (err error)
}

type Session interface {
	// 推送消息给当前session
	Push(resp gira.ProtoPush) (err error)
	Kick(ctx context.Context, reason string) (err error)
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

func Register(application gira.Application, proto gira.Proto, config Config, hallHandler HallHandler, playerHandler gira.ProtoHandler) (Hall, error) {
	service := &hall_service{
		facade:        application,
		proto:         proto,
		config:        config,
		hallHandler:   hallHandler,
		playerHandler: playerHandler,
		Actor:         actor.NewActor(16),
	}
	if err := service.register(); err != nil {
		return nil, err
	}
	go service.serve()
	return service, nil
}

type hall_service struct {
	facade gira.Application
	Config *Config
	Hall   Hall
	// 由application设置
	hallServer    *hall_server
	Proto         gira.Proto
	PlayerHandler gira.ProtoHandler
	*actor.Actor
	ctx           context.Context
	cancelFunc    context.CancelFunc
	SessionDict   sync.Map
	sessionCount  int64
	hallHandler   HallHandler
	playerHandler gira.ProtoHandler
	proto         gira.Proto
	config        Config
	isDestory     bool
}

func (hall *hall_service) register() error {
	facade.RegisterGrpc(func(server *grpc.Server) error {
		hall_grpc.RegisterHallServer(server, &hall_server{
			application: hall.facade,
			hall:        hall,
		})
		return nil
	})
	if _, err := facade.RegisterService(hall_grpc.HallServiceName, registry_options.WithRegisterAsGroupOption(true), registry_options.WithRegisterCatAppIdOption(true)); err != nil {
		return err
	}
	return nil
}

func (hall *hall_service) onDestory() {
	hall.isDestory = true
	for {
		log.Infow("hall停止中", "session_count", hall.sessionCount)
		sessions := make([]*hall_sesssion, 0)
		hall.SessionDict.Range(func(key, value any) bool {
			session, _ := value.(*hall_sesssion)
			sessions = append(sessions, session)
			return true
		})
		for _, session := range sessions {
			session.Call_close(hall.ctx, actor.WithCallTimeOut(5*time.Second))
		}
		log.Infow("hall停止中", "session_count", hall.sessionCount)
		if hall.sessionCount <= 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (hall *hall_service) serve() {
	//
	// 1.服务器关闭时保存数据后再退出
	// 2.处理actor调用
	//
	hall.ctx, hall.cancelFunc = context.WithCancel(facade.Context())
	defer hall.cancelFunc()
	// facade.Go(func() error {
	// 	select {
	// 	case <-hall.ctx.Done():
	// 		// 服务器停止，保存数据
	// 		for {
	// 			log.Infow("hall停止中", "session_count", hall.sessionCount)
	// 			if hall.sessionCount <= 0 {
	// 				break
	// 			}
	// 			time.Sleep(100 * time.Millisecond)
	// 		}
	// 	}
	// 	return nil
	// })
	for {
		select {
		// 处理actor请求
		case r := <-hall.Inbox():
			r.Call()
		case <-hall.ctx.Done():
			log.Infow("hall exit")
			return
		}
	}
}

func (hall *hall_service) SessionCount() int64 {
	return hall.sessionCount
}

// 推送消息给其他玩家
func (hall *hall_service) Push(ctx context.Context, userId string, push gira.ProtoPush) (err error) {
	defer func() {
		log.Warnw("hall_server push panic", "user_id", userId, "route", push.GetPushName())
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
			err = e.(error)
		}
	}()

	if v, ok := hall.SessionDict.Load(userId); !ok {
		var data []byte
		var peer *gira.Peer
		data, err = hall.proto.PushEncode(push)
		if err != nil {
			return
		}
		peer, err = facade.WhereIsUser(userId)
		if err != nil {
			return
		}
		if _, err = hall_grpc.DefaultHallClients.WithUnicast().WithPeer(peer).MustPush(ctx, &hall_grpc.MustPushRequest{
			UserId: userId,
			Data:   data,
		}); err != nil {
			return
		}
		return
	} else {
		session, _ := v.(*hall_sesssion)
		// WARN: chPush有可能已经关闭
		session.chPeerPush <- push
		return
	}
}

// 推送消息给其他玩家
func (hall *hall_service) MustPush(ctx context.Context, userId string, push gira.ProtoPush) (err error) {
	defer func() {
		log.Warnw("hall_server must push panic", "user_id", userId, "route", push.GetPushName())
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
			err = e.(error)
		}
	}()
	if v, ok := hall.SessionDict.Load(userId); !ok {
		var data []byte
		var peer *gira.Peer
		data, err = hall.proto.PushEncode(push)
		if err != nil {
			return
		}
		peer, err = facade.WhereIsUser(userId)
		if err != nil {
			return
		}
		if _, err = hall_grpc.DefaultHallClients.WithUnicast().WithPeer(peer).MustPush(ctx, &hall_grpc.MustPushRequest{
			UserId: userId,
			Data:   data,
		}); err != nil {
			return
		}
		return
	} else {
		session, _ := v.(*hall_sesssion)
		// WARN: chPush有可能已经关闭
		session.chPeerPush <- push
		return
	}
}
