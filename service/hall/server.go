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
	"github.com/lujingwei002/gira/options/service_options"
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
	return facade.NewServiceName(hall_grpc.HallServiceName, service_options.WithAsAppServiceOption(true))
}

type HallService struct {
	*actor.Actor
	facade        gira.Application
	hallServer    *hall_server
	ctx           context.Context
	cancelCtx     context.Context
	cancelFunc    context.CancelFunc
	sessionDict   sync.Map
	sessionCount  int64
	hallHandler   HallHandler
	playerHandler gira.ProtoHandler
	proto         gira.Proto
	config        Config
	isDestory     bool
}

func (hall *HallService) OnStart(ctx context.Context) error {
	hall.ctx = ctx
	// 后台运行
	hall.cancelCtx, hall.cancelFunc = context.WithCancel(context.Background())
	if err := facade.RegisterGrpc(func(server *grpc.Server) error {
		hall_grpc.RegisterHallServer(server, hall.hallServer)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (hall *HallService) OnStop() error {
	return nil
}

func (hall *HallService) Serve() error {
	//
	// 1.服务器关闭时保存数据后再退出
	// 2.处理actor调用
	//
	// hall.ctx, hall.cancelFunc = context.WithCancel(facade.Context())
	// defer hall.cancelFunc()
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
			goto TAG_CLEAN_UP
		}
	}
TAG_CLEAN_UP:
	hall.isDestory = true
	for {
		log.Infow("hall on stop------------", "session_count", hall.sessionCount)
		sessions := make([]*hall_sesssion, 0)
		hall.sessionDict.Range(func(key, value any) bool {
			session, _ := value.(*hall_sesssion)
			sessions = append(sessions, session)
			return true
		})
		for _, session := range sessions {
			session.Call_close(hall.cancelCtx, actor.WithCallTimeOut(5*time.Second))
		}
		log.Infow("hall on stop++++++++", "session_count", hall.sessionCount)
		if hall.sessionCount <= 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	hall.cancelFunc()
	return nil
}

func (hall *HallService) SessionCount() int64 {
	return hall.sessionCount
}

// 踢用户下线
// 协程程安全
func (hall *HallService) Kick(ctx context.Context, userId string, reason string) (err error) {
	if v, ok := hall.sessionDict.Load(userId); !ok {
		return gira.ErrUserNotFound
	} else {
		session := v.(*hall_sesssion)
		return session.Call_Kick(ctx, reason, actor.WithCallTimeOut(5*time.Second))
	}
}

// 顶号下线
// 协程程安全
func (hall *HallService) Instead(ctx context.Context, userId string, reason string) error {
	log.Infow("user instead", "user_id", userId)
	if v, ok := hall.sessionDict.Load(userId); !ok {
		// 偿试解锁
		peer, err := facade.UnlockLocalUser(userId)
		log.Infow("unlock local user return", "user_id", userId, "peer", peer, "err", err)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		session := v.(*hall_sesssion)
		timeoutCtx, timeoutFunc := context.WithTimeout(ctx, 10*time.Second)
		defer timeoutFunc()
		if err := session.Call_instead(timeoutCtx, reason, actor.WithCallTimeOut(1*time.Second)); err != nil {
			log.Errorw("user instead fail", "user_id", userId, "error", err)
			return err
		} else {
			log.Infow("user instead success", "user_id", userId)
			return nil
		}
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
