package app

import (
	"context"

	"github.com/lujingwei002/gira"
)

type BaseFacade struct {
	application *Application
	gateHandler gira.GateHandler
}

type AdminClient interface {
	BroadcastReloadResource(ctx context.Context, name string) error
}

func (self *BaseFacade) GetAppId() int32 {
	return self.application.appId
}

func (self *BaseFacade) GetAppType() string {
	return self.application.appType
}

func (self *BaseFacade) GetAppName() string {
	return self.application.appName
}

func (self *BaseFacade) GetAppFullName() string {
	return self.application.appFullName
}

func (self *BaseFacade) Go(f func() error) {
	self.application.errGroup.Go(f)
}
func (self *BaseFacade) Done() <-chan struct{} {
	return self.application.ctx.Done()
}

func (self *BaseFacade) Quit() {
	self.application.cancelFunc()
}

func (self *BaseFacade) Context() context.Context {
	return self.application.ctx
}

func (self *BaseFacade) SetApplication(application *Application) {
	self.application = application
	gira.SetApp(application.Facade)
}

func (self *BaseFacade) GetWorkDir() string {
	return self.application.WorkDir
}
func (self *BaseFacade) GetLogDir() string {
	return self.application.LogDir
}
func (self *BaseFacade) Wait() error {
	return self.application.wait()
}

func (self *BaseFacade) GetAccountDbClient() gira.MongoClient {
	return self.application.AccountDbClient
}

func (self *BaseFacade) GetGameDbClient() gira.MongoClient {
	return self.application.GameDbClient
}

func (self *BaseFacade) GetStatDbClient() gira.MongoClient {
	return self.application.StatDbClient
}

func (self *BaseFacade) GetAccountCacheClient() gira.RedisClient {
	return self.application.AccountCacheClient
}

func (self *BaseFacade) GetResourceDbClient() gira.MongoClient {
	return self.application.ResourceDbClient
}

func (self *BaseFacade) GetAdminDbClient() gira.MysqlClient {
	return self.application.adminDbClient
}

func (self *BaseFacade) SdkLogin(accountPlat string, openId string, token string) (*gira.SdkAccount, error) {
	return self.application.Sdk.Login(accountPlat, openId, token)
}

func (self *BaseFacade) LockLocalMember(memberId string) (*gira.Peer, error) {
	return self.application.Registry.LockLocalMember(memberId)
}

func (self *BaseFacade) UnlockLocalMember(memberId string) (*gira.Peer, error) {
	return self.application.Registry.UnlockLocalMember(memberId)
}

func (self *BaseFacade) OnLocalPlayerAdd(player *gira.LocalPlayer) {

}
func (self *BaseFacade) OnLocalPlayerDelete(player *gira.LocalPlayer) {

}
func (self *BaseFacade) OnLocalPlayerUpdate(player *gira.LocalPlayer) {

}

func (self *BaseFacade) ReloadResource() error {
	if self.application.resourceLoader == nil {
		return gira.ErrResourceLoaderNotImplement
	}
	return self.application.resourceLoader.ReloadResource("resource")
}

func (self *BaseFacade) RangePeers(f func(k any, v any) bool) {
	self.application.Registry.RangePeers(f)
}

func (self *BaseFacade) BroadcastReloadResource(ctx context.Context, name string) error {
	return self.application.adminClient.BroadcastReloadResource(ctx, name)
}
