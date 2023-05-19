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
	"google.golang.org/grpc"
)

func NewService(application gira.Application, proto gira.Proto, config Config, hallHandler HallHandler, playerHandler gira.ProtoHandler) (*HallService, error) {
	service := &HallService{
		facade:        application,
		proto:         proto,
		config:        config,
		hallHandler:   hallHandler,
		playerHandler: playerHandler,
		Actor:         actor.NewActor(16),
	}
	service.hallServer = &hall_server{
		hall: service,
	}
	return service, nil
}

func GetServiceName() string {
	return facade.NewServiceName(hall_grpc.HallServiceName, registry_options.WithRegisterAsGroupOption(true), registry_options.WithRegisterCatAppIdOption(true))
}

type HallService struct {
	facade     gira.Application
	hallServer *hall_server
	*actor.Actor
	ctx           context.Context
	cancelFunc    context.CancelFunc
	sessionDict   sync.Map
	sessionCount  int64
	hallHandler   HallHandler
	playerHandler gira.ProtoHandler
	proto         gira.Proto
	config        Config
	isDestory     bool
}

func (hall *HallService) OnStart() error {
	facade.RegisterGrpc(func(server *grpc.Server) error {
		hall_grpc.RegisterHallServer(server, hall.hallServer)
		return nil
	})
	go hall.serve()
	return nil
}

func (hall *HallService) OnStop() {
	hall.isDestory = true
	for {
		log.Infow("hall停止中", "session_count", hall.sessionCount)
		sessions := make([]*hall_sesssion, 0)
		hall.sessionDict.Range(func(key, value any) bool {
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

func (hall *HallService) serve() {
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

func (hall *HallService) SessionCount() int64 {
	return hall.sessionCount
}

func (hall *HallService) Kick(ctx context.Context, userId string, reason string) (err error) {
	if v, ok := hall.sessionDict.Load(userId); !ok {
		return gira.ErrUserNotFound
	} else {
		session := v.(*hall_sesssion)
		return session.Call_Kick(ctx, reason, actor.WithCallTimeOut(5*time.Second))
	}
}

// 推送消息给其他玩家
func (hall *HallService) Push(ctx context.Context, userId string, push gira.ProtoPush) (err error) {
	defer func() {
		log.Warnw("hall_server push panic", "user_id", userId, "route", push.GetPushName())
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
			err = e.(error)
		}
	}()

	if v, ok := hall.sessionDict.Load(userId); !ok {
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
func (hall *HallService) MustPush(ctx context.Context, userId string, push gira.ProtoPush) (err error) {
	defer func() {
		log.Warnw("hall_server must push panic", "user_id", userId, "route", push.GetPushName())
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
			err = e.(error)
		}
	}()
	if v, ok := hall.sessionDict.Load(userId); !ok {
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
