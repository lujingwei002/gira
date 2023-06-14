package app

import (
	"context"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/framework/smallgame/game"
	"github.com/lujingwei002/gira/framework/smallgame/game/config"
	"github.com/lujingwei002/gira/framework/smallgame/game/hall"
	"github.com/lujingwei002/gira/log"
)

type Framework struct {
	hallService hall.HallService
	// 使用的协议，当中必须包括名为Login的协议
	proto         gira.Proto
	config        *config.Config
	hallHandler   game.HallHandler
	playerHandler gira.ProtoHandler
}

func NewFramework(proto gira.Proto, hallHandler game.HallHandler, playerHandler gira.ProtoHandler) game.GameFramework {
	return &Framework{
		proto:         proto,
		hallHandler:   hallHandler,
		playerHandler: playerHandler,
	}
}

func (framework *Framework) OnFrameworkCreate(application gira.Application) error {
	// 注册大厅服务
	if hallService, err := hall.NewService(
		framework.proto,
		framework.config.Framework.Game,
		framework.hallHandler,
		framework.playerHandler); err != nil {
		return err
	} else {
		framework.hallService = hallService
	}
	return nil
}

func (framework *Framework) OnFrameworkStart() error {
	if err := facade.StartService("hall", framework.hallService); err != nil {
		return err
	}
	return nil
}

func (framework *Framework) OnFrameworkStop() error {
	if err := facade.StopService(framework.hallService); err != nil {
		log.Warn(err)
	}
	return nil
}

func (framework *Framework) OnFrameworkConfigLoad(c *gira.Config) error {
	framework.config = &config.Config{}
	return framework.config.OnConfigLoad(c)
}

// implment hall.Hall
func (framework *Framework) GetConfig() *config.GameConfig {
	return &framework.config.Framework.Game
}

func (framework *Framework) SessionCount() int64 {
	return framework.hallService.SessionCount()
}

func (framework *Framework) Push(ctx context.Context, userId string, req gira.ProtoPush) error {
	return framework.hallService.Push(ctx, userId, req)
}

func (framework *Framework) MustPush(ctx context.Context, userId string, req gira.ProtoPush) (err error) {
	return framework.hallService.MustPush(ctx, userId, req)
}
