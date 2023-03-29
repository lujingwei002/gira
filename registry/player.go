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

type PlayerRegistry struct {
	PeerPrefix   string // /peer_player/<<Name>>/      根据服务类型查找全部玩家
	LocalPrefix  string // /local_player/<<FullName>>  根据服务全名查找全部玩家
	UserPrefix   string // /player/<<UserId>>/       可以根据user_id查找当前所在的服
	Peers        map[string]*gira.Peer
	LocalPlayers sync.Map
	SelfPeer     *gira.Peer
	cancelCtx    context.Context
	cancelFunc   context.CancelFunc
}

func newConfigPlayerRegistry(r *Registry) (*PlayerRegistry, error) {
	self := &PlayerRegistry{
		PeerPrefix:  fmt.Sprintf("/peer_user/%s/", r.Name),
		LocalPrefix: fmt.Sprintf("/local_user/%s/", r.FullName),
		UserPrefix:  fmt.Sprintf("/user/"),
		Peers:       make(map[string]*gira.Peer, 0),
	}
	self.cancelCtx, self.cancelFunc = context.WithCancel(r.cancelCtx)
	// 注册自己
	// if err := self.registerSelf(r); err != nil {
	// 	return nil, err
	// }
	// 侦听本服的player信息
	if err := self.watchLocalPlayers(r); err != nil {
		return nil, err
	}
	r.facade.Go(func() error {
		select {
		case <-r.facade.Done():
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

func (self *PlayerRegistry) notify(r *Registry) error {
	self.LocalPlayers.Range(func(k any, v any) bool {
		player := v.(*gira.LocalPlayer)
		self.onLocalPlayerAdd(r, player)
		return true
	})
	return nil
}

func (self *PlayerRegistry) onLocalPlayerAdd(r *Registry, player *gira.LocalPlayer) error {
	log.Infof("============ local user %s add ==================", player.UserId)
	if handler, ok := r.facade.(gira.LocalPlayerHandler); ok {
		handler.OnLocalPlayerAdd(player)
	} else {
		panic(gira.ErrPeerHandlerNotImplement)
	}
	return nil
}

func (self *PlayerRegistry) onLocalPeerDelete(r *Registry, player *gira.LocalPlayer) error {
	log.Infof("============ local user %s delete ==================", player.UserId)
	if handler, ok := r.facade.(gira.LocalPlayerHandler); ok {
		handler.OnLocalPlayerDelete(player)
	} else {
		panic(gira.ErrPeerHandlerNotImplement)
	}
	return nil
}

func (self *PlayerRegistry) onLocalPeerUpdate(r *Registry, player *gira.LocalPlayer) error {
	log.Infof("============ local user %s update ==================", player.UserId)
	if handler, ok := r.facade.(gira.LocalPlayerHandler); ok {
		handler.OnLocalPlayerUpdate(player)
	} else {
		panic(gira.ErrPeerHandlerNotImplement)
	}
	return nil
}

func (self *PlayerRegistry) onKvPut(r *Registry, kv *mvccpb.KeyValue) error {
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
	if lastValue, ok := self.LocalPlayers.Load(userId); ok {
		lastPlayer := lastValue.(*gira.LocalPlayer)
		log.Info("player registry add player, but already exist", userId, "=>", value, lastPlayer.LoginTime)
	} else {
		// 新增player
		log.Info("player registry add player", userId, "=>", value)
		player := &gira.LocalPlayer{
			LoginTime: loginTime,
			UserId:    userId,
		}
		self.LocalPlayers.Store(userId, player)
		self.onLocalPlayerAdd(r, player)
	}
	return nil
}

func (self *PlayerRegistry) onKvDelete(r *Registry, kv *mvccpb.KeyValue) error {
	pats := strings.Split(string(kv.Key), "/")
	if len(pats) != 4 {
		log.Info("player registry got a invalid player", string(kv.Key))
		return gira.ErrInvalidPeer
	}
	userId := pats[3]
	value := string(kv.Value)
	if lastPlayer, ok := self.LocalPlayers.Load(userId); ok {
		log.Info("player registry remove player", userId, "=>", value, lastPlayer)
		self.LocalPlayers.Delete(userId)
	} else {
		log.Info("player registry remote player, but player not found", userId, "=>", value)
	}
	return nil
}

// 只增加节点，但不通知handler, 等notify再通知
func (self *PlayerRegistry) onKvAdd(kv *mvccpb.KeyValue) error {
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
	if lastPlayer, ok := self.Peers[userId]; ok {
		log.Info("player registry add player, but already exist", userId, lastPlayer, value)
	} else {
		player := &gira.LocalPlayer{
			LoginTime: loginTime,
			UserId:    userId,
		}
		self.LocalPlayers.Store(userId, player)
		log.Info("player registry add player", userId, "=>", value)
	}
	return nil
}

func (self *PlayerRegistry) watchLocalPlayers(r *Registry) error {
	client := r.Client
	kv := clientv3.NewKV(client)
	var getResp *clientv3.GetResponse
	var err error
	if getResp, err = kv.Get(self.cancelCtx, self.LocalPrefix, clientv3.WithPrefix()); err != nil {
		return err
	}
	for _, kv := range getResp.Kvs {
		if err := self.onKvAdd(kv); err != nil {
			return err
		}
	}
	watchStartRevision := getResp.Header.Revision + 1
	watcher := clientv3.NewWatcher(client)
	r.facade.Go(func() error {
		watchRespChan := watcher.Watch(self.cancelCtx, self.LocalPrefix, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix(), clientv3.WithPrevKV())
		log.Info("etcd watch player started", self.LocalPrefix, watchStartRevision)
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

func (self *PlayerRegistry) unregisterSelf(r *Registry) error {
	client := r.Client
	kv := clientv3.NewKV(client)
	ctx, cancelFunc := context.WithTimeout(self.cancelCtx, 10*time.Second)
	defer cancelFunc()
	log.Info("etcd unregister self", self.LocalPrefix)

	var txnResp *clientv3.TxnResponse
	var err error
	self.LocalPlayers.Range(func(userId any, v any) bool {
		txn := kv.Txn(ctx)
		localKey := fmt.Sprintf("%s%s", self.LocalPrefix, userId)
		peerKey := fmt.Sprintf("%s%s", self.PeerPrefix, userId)
		userKey := fmt.Sprintf("%s%s/%s", self.UserPrefix, userId, r.FullName)
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

func (self *PlayerRegistry) LockLocalUser(r *Registry, userId string) (*gira.Peer, error) {
	if _, ok := self.LocalPlayers.Load(userId); ok {
		// return r.PeerRegistry.SelfPeer, nil
	}
	client := r.Client
	// 到etcd抢占localKey
	localKey := fmt.Sprintf("%s%s", self.LocalPrefix, userId)
	peerKey := fmt.Sprintf("%s%s", self.PeerPrefix, userId)
	userKey := fmt.Sprintf("%s%s/%s", self.UserPrefix, userId, r.FullName)
	value := fmt.Sprintf("%d", time.Now().Unix())
	kv := clientv3.NewKV(client)
	var err error
	var txnResp *clientv3.TxnResponse
	txn := kv.Txn(self.cancelCtx)
	log.Info("player registry local key", localKey)
	log.Info("player registry peer key", peerKey)
	log.Info("player registry user key", userKey)
	txn.If(clientv3.Compare(clientv3.CreateRevision(peerKey), "=", 0)).
		Then(clientv3.OpPut(localKey, value), clientv3.OpPut(peerKey, r.FullName), clientv3.OpPut(userKey, value)).
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

func (self *PlayerRegistry) UnlockLocalUser(r *Registry, userId string) (*gira.Peer, error) {
	client := r.Client
	localKey := fmt.Sprintf("%s%s", self.LocalPrefix, userId)
	peerKey := fmt.Sprintf("%s%s", self.PeerPrefix, userId)
	userKey := fmt.Sprintf("%s%s/%s", self.UserPrefix, userId, r.FullName)
	value := fmt.Sprintf("%d", time.Now().Unix())
	kv := clientv3.NewKV(client)
	var err error
	var txnResp *clientv3.TxnResponse
	txn := kv.Txn(self.cancelCtx)
	log.Info("player registry local key", localKey)
	log.Info("player registry peer key", peerKey)
	log.Info("player registry user key", userKey)
	txn.If(clientv3.Compare(clientv3.Value(peerKey), "=", r.FullName), clientv3.Compare(clientv3.CreateRevision(peerKey), "!=", 0)).
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

/*
func (self *PlayerRegistry) registerSelf(r *Registry) error {
	client := r.Client
	var err error
	var lease clientv3.Lease
	var leaseID clientv3.LeaseID
	if r.Config.LeaseTimeout > 0 {
		// 申请一个租约 lease
		lease = clientv3.Lease(client)
		var leaseGrantResp *clientv3.LeaseGrantResponse
		// 申请一个5s的租约
		if leaseGrantResp, err = lease.Grant(r.cancelCtx, 5); err != nil {
			return err
		}
		// 租约ID
		leaseID = leaseGrantResp.ID
	}

	// 需要同步的键值对
	advertises := make(map[string]string, 0)
	advertises[GRPC_KEY] = r.Config.Address
	for _, v := range r.Config.Advertise {
		advertises[v.Name] = v.Value
	}
	kv := clientv3.NewKV(client)
	for name, value := range advertises {
		var txnResp *clientv3.TxnResponse
		txn := kv.Txn(r.cancelCtx)
		key := fmt.Sprintf("%s%s", self.SelfPrefix, name)
		tx := txn.If(clientv3.Compare(clientv3.Value(key), "!=", value), clientv3.Compare(clientv3.CreateRevision(key), "!=", 0)).
			Then(clientv3.OpGet(key))
		if leaseID != 0 {
			tx.Else(clientv3.OpGet(key), clientv3.OpPut(key, value, clientv3.WithLease(leaseID)))
		} else {
			tx.Else(clientv3.OpGet(key), clientv3.OpPut(key, value))
		}
		if txnResp, err = txn.Commit(); err != nil {
			log.Info("txn err", err)
			return err
		}
		if txnResp.Succeeded {
			log.Info("etcd register", key, "=>", value, "failed", "lock by", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
			return gira.ErrRegisterServerFail
		} else {
			if len(txnResp.Responses[0].GetResponseRange().Kvs) == 0 {
				log.Info("etcd register create", key, "=>", value, "success")
			} else {
				log.Info("etcd register resume", key, "=>", value, "success")
			}
		}
	}
	if leaseID != 0 {
		var keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
		// 自动续租
		if keepRespChan, err = lease.KeepAlive(r.cancelCtx, leaseID); err != nil {
			log.Info("自动续租失败", err)
			return err
		}
		//判断续约应答的协程
		r.facade.Go(func() error {
			for {
				select {
				case keepResp := <-keepRespChan:
					if keepRespChan == nil {
						log.Info("租约已经失效了")
						goto END
					} else if keepResp == nil {
						log.Info("租约已经被取消")
						goto END
					} else {
						// KeepAlive每秒会续租一次,所以就会收到一次应答
						// log.Info("收到应答,租约ID是:", keepResp.ID)
					}
				case <-r.facade.Done():
					break
				}
			}
		END:
			return nil
		})
	}
	return nil
}
*/
