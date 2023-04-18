package hall

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/framework/smallgame/common/rpc"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// 应用
type Framework struct {
	Config *Config
	Hall   Hall
	// 由application设置
	Proto         gira.Proto
	PlayerHandler gira.ProtoHandler
}

func (framework *Framework) OnFrameworkAwake(application gira.Application) error {
	if handler, ok := application.(HallHandler); !ok {
		return gira.ErrGateHandlerNotImplement
	} else {
		hall := newHall(framework, framework.Proto, framework.Config, handler, framework.PlayerHandler)
		framework.Hall = hall
		if err := hall.OnAwake(); err != nil {
			return err
		}
		rpc.OnAwake()
		return nil
	}
}

// 框架启动
func (framework *Framework) OnFrameworkStart() error {
	return nil
}

func (framework *Framework) OnFrameworkConfigLoad(c *gira.Config) error {
	framework.Config = &Config{}
	return framework.Config.OnConfigLoad(c)
}
