package registry

///
///
/// key不设置过期时间，程序正常退出时自动清理，非正常退出，要程序重启来解锁
///
/// 注册表结构:
///   /peer_type_user/<<AppName>>/<<UserId>> => <<AppFullName>>
///   /peer_user/<<AppFullName>>/<<UserId>> time
///   /user/<<UserId>> => <<AppFullName>>
///
import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
	mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type player_registry struct {
	peerTypePrefix     string   // /peer_type_user/<<AppName>>/      根据服务类型查找全部玩家
	peerPrefix         string   // /peer_user/<<AppFullName>>  根据服务全名查找全部玩家
	userPrefix         string   // /user/<<UserId>>/       	 	可以根据user_id查找当前所在的服
	localPlayers       sync.Map // 本节点上的用户
	ctx                context.Context
	cancelFunc         context.CancelFunc
	watchStartRevision int64
}

func newConfigPlayerRegistry(r *Registry) (*player_registry, error) {
	self := &player_registry{
		peerTypePrefix: fmt.Sprintf("/peer_type_user/%s/", r.name),
		peerPrefix:     fmt.Sprintf("/peer_user/%s/", r.fullName),
		userPrefix:     "/user/",
	}
	return self, nil
}

func (self *player_registry) onStop(r *Registry) error {
	log.Debug("player registry on stop")
	if err := self.unregisterLocalPlayers(r); err != nil {
		log.Info(err)
	}
	return nil
}

func (self *player_registry) OnStart(r *Registry) error {
	cancelCtx, cancelFunc := context.WithCancel(r.ctx)
	self.ctx = cancelCtx
	self.cancelFunc = cancelFunc
	if err := self.initSelfPeerPlayers(r); err != nil {
		return err
	}
	return nil
}

func (self *player_registry) Serve(r *Registry) error {
	return self.watchSelfPeerPlayers(r)
}

func (self *player_registry) notify(r *Registry) error {
	self.localPlayers.Range(func(k any, v any) bool {
		player := v.(*gira.LocalPlayer)
		self.onLocalPlayerAdd(r, player)
		return true
	})
	return nil
}

func (self *player_registry) onLocalPlayerAdd(r *Registry, player *gira.LocalPlayer) {
	if r.isNotify == 0 {
		return
	}
	log.Infow("player registry on local user add", "user_id", player.UserId)
	for _, fw := range r.application.Frameworks() {
		if handler, ok := fw.(gira.LocalPlayerWatchHandler); ok {
			handler.OnLocalPlayerAdd(player)
		}
	}
	if handler, ok := r.application.(gira.LocalPlayerWatchHandler); ok {
		handler.OnLocalPlayerAdd(player)
	}
}

func (self *player_registry) onLocalPlayerDelete(r *Registry, player *gira.LocalPlayer) {
	log.Infow("player registry on local user add", "delete", player.UserId)
	for _, fw := range r.application.Frameworks() {
		if handler, ok := fw.(gira.LocalPlayerWatchHandler); ok {
			handler.OnLocalPlayerDelete(player)
		}
	}
	if handler, ok := r.application.(gira.LocalPlayerWatchHandler); ok {
		handler.OnLocalPlayerDelete(player)
	}
}

// func (self *player_registry) onLocalPlayerUpdate(r *Registry, player *gira.LocalPlayer) {
// 	log.Infow("player registry on local user add", "update", player.UserId)
// 	for _, fw := range r.application.Frameworks() {
// 		if handler, ok := fw.(gira.LocalPlayerWatchHandler); ok {
// 			handler.OnLocalPlayerUpdate(player)
// 		}
// 	}
// 	if handler, ok := r.application.(gira.LocalPlayerWatchHandler); ok {
// 		handler.OnLocalPlayerUpdate(player)
// 	}
// }

func (self *player_registry) onLocalKvAdd(r *Registry, kv *mvccpb.KeyValue) error {
	pats := strings.Split(string(kv.Key), "/")
	if len(pats) != 4 {
		log.Warnw("player registry got a invalid key", "key", string(kv.Key))
		return gira.ErrInvalidPeer
	}
	userId := pats[3]
	value := string(kv.Value)
	loginTime, err := strconv.Atoi(value)
	if err != nil {
		log.Warnw("player registry got a invalid value", "value", string(kv.Value))
		return err
	}
	if _, ok := self.localPlayers.Load(userId); ok {
		log.Warnw("player registry add local player, but already exist", "user_id", userId)
	} else {
		// 新增player
		log.Infow("player registry add local player", "user_id", userId)
		player := &gira.LocalPlayer{
			LoginTime: loginTime,
			UserId:    userId,
		}
		self.localPlayers.Store(userId, player)
		self.onLocalPlayerAdd(r, player)
	}
	return nil
}

func (self *player_registry) onLocalKvDelete(r *Registry, kv *mvccpb.KeyValue) error {
	pats := strings.Split(string(kv.Key), "/")
	if len(pats) != 4 {
		log.Warnw("player registry got a invalid key", "key", string(kv.Key))
		return gira.ErrInvalidPeer
	}
	userId := pats[3]
	if lastValue, ok := self.localPlayers.Load(userId); ok {
		lastPlayer := lastValue.(*gira.LocalPlayer)
		log.Infow("player registry remove local player", "user_id", userId)
		self.localPlayers.Delete(userId)
		self.onLocalPlayerDelete(r, lastPlayer)
	} else {
		log.Warnw("player registry remove local player, but player not found", "user_id", userId)
	}
	return nil
}

func (self *player_registry) initSelfPeerPlayers(r *Registry) error {
	client := r.client
	kv := clientv3.NewKV(client)
	var getResp *clientv3.GetResponse
	var err error
	if getResp, err = kv.Get(self.ctx, self.peerPrefix, clientv3.WithPrefix()); err != nil {
		return err
	}
	for _, kv := range getResp.Kvs {
		if err := self.onLocalKvAdd(r, kv); err != nil {
			return err
		}
	}
	self.watchStartRevision = getResp.Header.Revision + 1
	return nil
}

func (self *player_registry) watchSelfPeerPlayers(r *Registry) error {
	client := r.client
	watchStartRevision := self.watchStartRevision
	watcher := clientv3.NewWatcher(client)
	// r.application.Go(func() error {
	watchRespChan := watcher.Watch(self.ctx, self.peerPrefix, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix(), clientv3.WithPrevKV())
	log.Infow("player registry started", "local_prefix", self.peerPrefix, "watch_start_revision", watchStartRevision)
	for watchResp := range watchRespChan {
		// log.Info("etcd watch got events")
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				// log.Info("etcd got put event")
				if err := self.onLocalKvAdd(r, event.Kv); err != nil {
					log.Warnw("player registry put event fail", "error", err)
				}
			case mvccpb.DELETE:
				// log.Info("etcd got delete event")
				if err := self.onLocalKvDelete(r, event.Kv); err != nil {
					log.Warnw("player registry put event fail", "error", err)
				}
			}
		}
	}
	log.Info("player registry watch shutdown")
	return nil
}

func (self *player_registry) unregisterLocalPlayers(r *Registry) error {
	client := r.client
	kv := clientv3.NewKV(client)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	log.Infow("player registry unregister self", "local_prefix", self.peerPrefix)

	var txnResp *clientv3.TxnResponse
	var err error
	self.localPlayers.Range(func(userId any, v any) bool {
		txn := kv.Txn(ctx)
		localKey := fmt.Sprintf("%s%s", self.peerPrefix, userId)
		peerKey := fmt.Sprintf("%s%s", self.peerTypePrefix, userId)
		userKey := fmt.Sprintf("%s%s", self.userPrefix, userId)
		log.Infow("player registry unregister self", "local_key", localKey, "peer_key", peerKey, "user_key", userKey)
		txn.If(clientv3.Compare(clientv3.CreateRevision(localKey), "!=", 0)).
			Then(clientv3.OpDelete(localKey), clientv3.OpDelete(peerKey), clientv3.OpDelete(userKey))

		if txnResp, err = txn.Commit(); err != nil {
			log.Errorw("player registry commit fail", "error", err)
			return true
		}
		if txnResp.Succeeded {
			log.Infow("player registry unregister", "user_id", userId)
		} else {
			log.Warnw("player registry unregister", "user_id", userId)
		}
		return true
	})
	return nil
}

// 锁定玩家
func (self *player_registry) LockLocalUser(r *Registry, userId string) (*gira.Peer, error) {
	//if _, ok := self.localPlayers.Load(userId); ok {
	//return r.peerRegistry.SelfPeer, nil
	//}
	client := r.client
	// 到etcd抢占localKey
	localKey := fmt.Sprintf("%s%s", self.peerPrefix, userId)
	peerKey := fmt.Sprintf("%s%s", self.peerTypePrefix, userId)
	userKey := fmt.Sprintf("%s%s", self.userPrefix, userId)
	value := fmt.Sprintf("%d", time.Now().Unix())
	kv := clientv3.NewKV(client)
	var err error
	var txnResp *clientv3.TxnResponse
	txn := kv.Txn(self.ctx)
	log.Infow("player registry", "local_key", localKey, "peer_key", peerKey, "user_key", userKey)
	txn.If(clientv3.Compare(clientv3.CreateRevision(peerKey), "=", 0)).
		Then(clientv3.OpPut(localKey, value), clientv3.OpPut(peerKey, r.fullName), clientv3.OpPut(userKey, r.fullName)).
		Else(clientv3.OpGet(peerKey))
	if txnResp, err = txn.Commit(); err != nil {
		log.Errorw("player registry commit fail", "error", err)
		return nil, err
	}
	if txnResp.Succeeded {
		log.Infow("player registry register success", "local_key", localKey)
		return nil, nil
	} else {
		var fullName string
		if len(txnResp.Responses) > 0 && len(txnResp.Responses[0].GetResponseRange().Kvs) > 0 {
			fullName = string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
		}
		log.Warnw("player registry register", localKey, "=>", value, "failed", "lock by", fullName)
		peer := r.GetPeer(fullName)
		if peer == nil {
			return nil, gira.ErrUserLocked
		}
		return peer, gira.ErrUserLocked
	}
}

// 解锁
func (self *player_registry) UnlockLocalUser(r *Registry, userId string) (*gira.Peer, error) {
	client := r.client
	localKey := fmt.Sprintf("%s%s", self.peerPrefix, userId)
	peerKey := fmt.Sprintf("%s%s", self.peerTypePrefix, userId)
	userKey := fmt.Sprintf("%s%s", self.userPrefix, userId)
	kv := clientv3.NewKV(client)
	var err error
	var txnResp *clientv3.TxnResponse
	txn := kv.Txn(self.ctx)
	log.Infow("player registry unregister", "local_key", localKey, "peer_key", peerKey, "user_key", userKey)
	txn.If(clientv3.Compare(clientv3.Value(peerKey), "=", r.fullName), clientv3.Compare(clientv3.CreateRevision(peerKey), "!=", 0)).
		Then(clientv3.OpDelete(localKey), clientv3.OpDelete(peerKey), clientv3.OpDelete(userKey)).
		Else(clientv3.OpGet(peerKey))
	if txnResp, err = txn.Commit(); err != nil {
		log.Errorw("player registry commit fail", "error", err)
		return nil, err
	}
	if txnResp.Succeeded {
		log.Infow("player registry unregister success", "local_key", localKey)
		return nil, nil
	} else {
		var fullName string
		if len(txnResp.Responses) > 0 && len(txnResp.Responses[0].GetResponseRange().Kvs) > 0 {
			fullName = string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
		}
		log.Warnw("player registry unregister fail", "local_key", localKey, "locked_by", fullName)
		peer := r.GetPeer(fullName)
		if peer == nil {
			return nil, gira.ErrUserLocked
		}
		return peer, gira.ErrUserLocked
	}
}

// 查找玩家位置
func (self *player_registry) WhereIsUser(r *Registry, userId string) (*gira.Peer, error) {
	if _, ok := self.localPlayers.Load(userId); ok {
		return r.peerRegistry.SelfPeer, nil
	}
	client := r.client
	// 到etcd抢占localKey
	userKey := fmt.Sprintf("%s%s", self.userPrefix, userId)
	kv := clientv3.NewKV(client)
	var err error
	getResp, err := kv.Get(self.ctx, userKey)
	if err != nil {
		return nil, err
	}
	if len(getResp.Kvs) <= 0 {
		return nil, gira.ErrUserNotFound.Trace()
	}
	fullName := string(getResp.Kvs[0].Value)
	peer := r.GetPeer(fullName)
	if peer == nil {
		return nil, gira.ErrPeerNotFound.Trace()
	} else {
		return peer, nil
	}
}
