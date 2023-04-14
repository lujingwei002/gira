package registry

///
///
/// key不设置过期时间，程序正常退出时自动清理，非正常退出，要程序重启来解锁
///
/// 注册表结构:
///   /peer_player/<<Name>>/<<UserId>> <<FullName>>
///   /local_player/<<FullName>>/<<UserId>> time
///   /player/<<UserId>>/<<FullName>> time
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
	peerPrefix   string // /peer_player/<<Name>>/      根据服务类型查找全部玩家
	localPrefix  string // /local_player/<<FullName>>  根据服务全名查找全部玩家
	userPrefix   string // /user/<<UserId>>/       可以根据user_id查找当前所在的服
	peers        map[string]*gira.Peer
	localPlayers sync.Map
	selfPeer     *gira.Peer
	ctx          context.Context
	cancelFunc   context.CancelFunc
}

func newConfigPlayerRegistry(r *Registry) (*player_registry, error) {
	self := &player_registry{
		peerPrefix:  fmt.Sprintf("/peer_user/%s/", r.name),
		localPrefix: fmt.Sprintf("/local_user/%s/", r.fullName),
		userPrefix:  fmt.Sprintf("/user/"),
		peers:       make(map[string]*gira.Peer, 0),
	}
	self.ctx, self.cancelFunc = context.WithCancel(r.cancelCtx)
	// 侦听本服的player信息
	if err := self.watchLocalPlayers(r); err != nil {
		return nil, err
	}
	r.application.Go(func() error {
		select {
		case <-r.application.Done():
			{
				log.Info("player registry recv down")
			}
		}
		if err := self.unregisterSelf(r); err != nil {
			log.Info(err)
		}
		self.cancelFunc()
		return nil
	})
	return self, nil
}

func (self *player_registry) notify(r *Registry) error {
	self.localPlayers.Range(func(k any, v any) bool {
		player := v.(*gira.LocalPlayer)
		self.onLocalPlayerAdd(r, player)
		return true
	})
	return nil
}

func (self *player_registry) onLocalPlayerAdd(r *Registry, player *gira.LocalPlayer) error {
	log.Infof("============ local user %s add ==================", player.UserId)
	for _, fw := range r.application.Frameworks() {
		if handler, ok := fw.(gira.LocalPlayerHandler); ok {
			handler.OnLocalPlayerAdd(player)
		}
	}
	if handler, ok := r.application.(gira.LocalPlayerHandler); ok {
		handler.OnLocalPlayerAdd(player)
	}
	return nil
}

func (self *player_registry) onLocalPeerDelete(r *Registry, player *gira.LocalPlayer) error {
	log.Infof("============ local user %s delete ==================", player.UserId)
	for _, fw := range r.application.Frameworks() {
		if handler, ok := fw.(gira.LocalPlayerHandler); ok {
			handler.OnLocalPlayerDelete(player)
		}
	}
	if handler, ok := r.application.(gira.LocalPlayerHandler); ok {
		handler.OnLocalPlayerDelete(player)
	}
	return nil
}

func (self *player_registry) onLocalPeerUpdate(r *Registry, player *gira.LocalPlayer) error {
	log.Infof("============ local user %s update ==================", player.UserId)
	for _, fw := range r.application.Frameworks() {
		if handler, ok := fw.(gira.LocalPlayerHandler); ok {
			handler.OnLocalPlayerUpdate(player)
		}
	}
	if handler, ok := r.application.(gira.LocalPlayerHandler); ok {
		handler.OnLocalPlayerUpdate(player)
	}
	return nil
}

func (self *player_registry) onKvPut(r *Registry, kv *mvccpb.KeyValue) error {
	pats := strings.Split(string(kv.Key), "/")
	if len(pats) != 4 {
		log.Info("player registry got a invalid player", string(kv.Key))
		return gira.ErrInvalidPeer
	}
	userId := pats[3]
	value := string(kv.Value)
	loginTime, err := strconv.Atoi(value)
	if err != nil {
		log.Info("player registry got a invalid player", string(kv.Value))
		return err
	}
	if lastValue, ok := self.localPlayers.Load(userId); ok {
		lastPlayer := lastValue.(*gira.LocalPlayer)
		log.Info("player registry add player, but already exist", userId, "=>", value, lastPlayer.LoginTime)
	} else {
		// 新增player
		log.Info("player registry add player", userId, "=>", value)
		player := &gira.LocalPlayer{
			LoginTime: loginTime,
			UserId:    userId,
		}
		self.localPlayers.Store(userId, player)
		self.onLocalPlayerAdd(r, player)
	}
	return nil
}

func (self *player_registry) onKvDelete(r *Registry, kv *mvccpb.KeyValue) error {
	pats := strings.Split(string(kv.Key), "/")
	if len(pats) != 4 {
		log.Info("player registry got a invalid player", string(kv.Key))
		return gira.ErrInvalidPeer
	}
	userId := pats[3]
	value := string(kv.Value)
	if lastPlayer, ok := self.localPlayers.Load(userId); ok {
		log.Info("player registry remove player", userId, "=>", value, lastPlayer)
		self.localPlayers.Delete(userId)
	} else {
		log.Info("player registry remote player, but player not found", userId, "=>", value)
	}
	return nil
}

// 只增加节点，但不通知handler, 等notify再通知
func (self *player_registry) onKvAdd(kv *mvccpb.KeyValue) error {
	pats := strings.Split(string(kv.Key), "/")
	if len(pats) != 4 {
		log.Info("player registry got a invalid player", string(kv.Key))
		return gira.ErrInvalidPeer
	}
	userId := pats[3]
	value := string(kv.Value)
	loginTime, err := strconv.Atoi(value)
	if err != nil {
		log.Info("player registry got a invalid player", string(kv.Value))
		return err
	}
	if lastPlayer, ok := self.peers[userId]; ok {
		log.Info("player registry add player, but already exist", userId, lastPlayer, value)
	} else {
		player := &gira.LocalPlayer{
			LoginTime: loginTime,
			UserId:    userId,
		}
		self.localPlayers.Store(userId, player)
		log.Info("player registry add player", userId, "=>", value)
	}
	return nil
}

func (self *player_registry) watchLocalPlayers(r *Registry) error {
	client := r.client
	kv := clientv3.NewKV(client)
	var getResp *clientv3.GetResponse
	var err error
	if getResp, err = kv.Get(self.ctx, self.localPrefix, clientv3.WithPrefix()); err != nil {
		return err
	}
	for _, kv := range getResp.Kvs {
		if err := self.onKvAdd(kv); err != nil {
			return err
		}
	}
	watchStartRevision := getResp.Header.Revision + 1
	watcher := clientv3.NewWatcher(client)
	r.application.Go(func() error {
		watchRespChan := watcher.Watch(self.ctx, self.localPrefix, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix(), clientv3.WithPrevKV())
		log.Info("etcd watch player started", self.localPrefix, watchStartRevision)
		for watchResp := range watchRespChan {
			// log.Info("etcd watch got events")
			for _, event := range watchResp.Events {
				switch event.Type {
				case mvccpb.PUT:
					// log.Info("etcd got put event")
					if err := self.onKvPut(r, event.Kv); err != nil {
						log.Info("player registry put event error", err)
					}
				case mvccpb.DELETE:
					// log.Info("etcd got delete event")
					if err := self.onKvDelete(r, event.Kv); err != nil {
						log.Info("player registry put event error", err)
					}
				}
			}
		}
		log.Info("player registry watch shutdown")
		return nil
	})
	return nil
}

func (self *player_registry) unregisterSelf(r *Registry) error {
	client := r.client
	kv := clientv3.NewKV(client)
	ctx, cancelFunc := context.WithTimeout(self.ctx, 10*time.Second)
	defer cancelFunc()
	log.Info("etcd unregister self", self.localPrefix)

	var txnResp *clientv3.TxnResponse
	var err error
	self.localPlayers.Range(func(userId any, v any) bool {
		txn := kv.Txn(ctx)
		localKey := fmt.Sprintf("%s%s", self.localPrefix, userId)
		peerKey := fmt.Sprintf("%s%s", self.peerPrefix, userId)
		userKey := fmt.Sprintf("%s%s", self.userPrefix, userId)
		txn.If(clientv3.Compare(clientv3.CreateRevision(localKey), "!=", 0)).
			Then(clientv3.OpDelete(localKey), clientv3.OpDelete(peerKey), clientv3.OpDelete(userKey))

		if txnResp, err = txn.Commit(); err != nil {
			log.Info("txn err", err)
			return true
		}
		if txnResp.Succeeded {
			log.Info("local user", userId, "delete")
		} else {
			log.Info("local user", userId, "delete, but not found")
		}
		return true
	})
	return nil
}

// 锁定玩家
func (self *player_registry) LockLocalUser(r *Registry, userId string) (*gira.Peer, error) {
	if _, ok := self.localPlayers.Load(userId); ok {
		//return r.peerRegistry.SelfPeer, nil
	}
	client := r.client
	// 到etcd抢占localKey
	localKey := fmt.Sprintf("%s%s", self.localPrefix, userId)
	peerKey := fmt.Sprintf("%s%s", self.peerPrefix, userId)
	userKey := fmt.Sprintf("%s%s", self.userPrefix, userId)
	value := fmt.Sprintf("%d", time.Now().Unix())
	kv := clientv3.NewKV(client)
	var err error
	var txnResp *clientv3.TxnResponse
	txn := kv.Txn(self.ctx)
	log.Info("player registry local key", localKey)
	log.Info("player registry peer key", peerKey)
	log.Info("player registry user key", userKey)
	txn.If(clientv3.Compare(clientv3.CreateRevision(peerKey), "=", 0)).
		Then(clientv3.OpPut(localKey, value), clientv3.OpPut(peerKey, r.fullName), clientv3.OpPut(userKey, r.fullName)).
		Else(clientv3.OpGet(peerKey))
	if txnResp, err = txn.Commit(); err != nil {
		log.Info("txn err", err)
		return nil, err
	}
	if txnResp.Succeeded {
		log.Info("player registry register", localKey, "=>", value, "success")
		return nil, nil
	} else {
		log.Info("player registry register", localKey, "=>", value, "failed", "lock by", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		fullName := string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
		peer := r.GetPeer(fullName)
		if peer == nil {
			return nil, gira.ErrPeerNotFound
		}
		return peer, gira.ErrUserLocked
	}
}

// 解锁
func (self *player_registry) UnlockLocalUser(r *Registry, userId string) (*gira.Peer, error) {
	client := r.client
	localKey := fmt.Sprintf("%s%s", self.localPrefix, userId)
	peerKey := fmt.Sprintf("%s%s", self.peerPrefix, userId)
	userKey := fmt.Sprintf("%s%s", self.userPrefix, userId)
	value := fmt.Sprintf("%d", time.Now().Unix())
	kv := clientv3.NewKV(client)
	var err error
	var txnResp *clientv3.TxnResponse
	txn := kv.Txn(self.ctx)
	log.Info("player registry local key", localKey)
	log.Info("player registry peer key", peerKey)
	log.Info("player registry user key", userKey)
	txn.If(clientv3.Compare(clientv3.Value(peerKey), "=", r.fullName), clientv3.Compare(clientv3.CreateRevision(peerKey), "!=", 0)).
		Then(clientv3.OpDelete(localKey), clientv3.OpDelete(peerKey), clientv3.OpDelete(userKey)).
		Else(clientv3.OpGet(peerKey))
	if txnResp, err = txn.Commit(); err != nil {
		log.Info("txn err", err)
		return nil, err
	}
	if txnResp.Succeeded {
		log.Info("player registry unregister", localKey, "=>", value, "success")
		return nil, nil
	} else {
		var fullName string
		if len(txnResp.Responses) > 0 && len(txnResp.Responses[0].GetResponseRange().Kvs) > 0 {
			fullName = string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
		}
		log.Info("player registry unregister", localKey, "=>", value, "failed", "lock by", fullName)
		peer := r.GetPeer(fullName)
		if peer == nil {
			return nil, gira.ErrPeerNotFound
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
		return nil, gira.ErrPeerNotFound
	}
	fullName := string(getResp.Kvs[0].Value)
	peer := r.GetPeer(fullName)
	if peer == nil {
		return nil, gira.ErrPeerNotFound
	} else {
		return peer, nil
	}
}
