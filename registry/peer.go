package registry

/// ## server key
/// server key 会设置过期时间。
///	如果server异常关闭，key就没法续期，其他server就会被通知到
/// 但如果由于网络波动等原因，导致key没法续期，导致server key会被非正常delete掉, 其他server会被通知，等网络恢复正常，当前的server如果watch到后，应该要偿试重新注册自己。
///	所以这里有两种原因，会watch到自己被反注册，一个是意外情况，一个是正常关闭。重新注册应该区分开来分别处理。

/// ## go routinue
/// watch routine 监听peer变化
/// lease routine 定时续期
/// ctrl routine 监听程序结束

/// ## context
// app
//  |
//  ctx --- cancelCtx
//  |    |
//  |    --errCtx
//
//
// background
//  |
//  -----watchCtx---unregister self
//
// 程序结束后，手动关闭watchCtx, 等成功从注册表删除了自己再结束

// ## 特殊情况
// 1.程序退出时，如果server key被*自己*占用着，则要等删除了，再退出
// 2.程序退出时，如果server key没被*自己*占用着，则可以直接退出了
// 3.如果由于网络等异常原因，导致server key过期，则要重新抢占，还要设置重试(TODO)

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type PeerRegistry struct {
	SelfPrefix string // /peer/<<FullName>>
	Prefix     string // /peer/
	// TODO 要加锁
	Peers                  sync.Map //map[string]*gira.Peer
	isNormalUnregisterSelf bool
	SelfPeer               *gira.Peer
	cancelCtx              context.Context
	cancelFunc             context.CancelFunc
}

func newConfigPeerRegistry(r *Registry) (*PeerRegistry, error) {
	self := &PeerRegistry{
		Prefix:     "/peer/",
		SelfPrefix: fmt.Sprintf("/peer/%s/", r.FullName),
	}
	self.cancelCtx, self.cancelFunc = context.WithCancel(r.cancelCtx)
	// 注册自己
	if err := self.registerSelf(r); err != nil {
		return nil, err
	}
	// 侦听伙伴信息
	if err := self.watchPeers(r); err != nil {
		return nil, err
	}
	r.application.Go(func() error {
		select {
		case <-r.application.Done():
			{
				log.Infow("peer registry recv down")
			}
		}
		if err := self.unregisterSelf(r); err != nil {
			log.Errorw("unregister self fail", "error", err)
		}

		return nil
	})
	return self, nil
}
func (self *PeerRegistry) RangePeers(f func(k any, v any) bool) {
	self.Peers.Range(f)
}
func (self *PeerRegistry) getPeer(r *Registry, fullName string) *gira.Peer {
	if lastValue, ok := self.Peers.Load(fullName); ok {
		lastPeer := lastValue.(*gira.Peer)
		return lastPeer
	}
	return nil
}

func (self *PeerRegistry) notify(r *Registry) error {
	self.Peers.Range(func(k any, v any) bool {
		peer := v.(*gira.Peer)
		self.onPeerAdd(r, peer)
		return true
	})
	return nil
}

func (self *PeerRegistry) onPeerAdd(r *Registry, peer *gira.Peer) error {
	log.Debugw("============ peer add ==================", "full_name", peer.FullName)
	if peer.FullName == r.FullName {
		self.SelfPeer = peer
	}
	for _, fw := range r.application.Frameworks() {
		if handler, ok := fw.(gira.PeerHandler); ok {
			handler.OnPeerAdd(peer)
		}
	}
	if handler, ok := r.application.(gira.PeerHandler); ok {
		handler.OnPeerAdd(peer)
	}
	return nil
}

func (self *PeerRegistry) onPeerDelete(r *Registry, peer *gira.Peer) error {
	log.Debugw("============ peer delete ==================", "full_name", peer.FullName)
	for _, fw := range r.application.Frameworks() {
		if handler, ok := fw.(gira.PeerHandler); ok {
			handler.OnPeerDelete(peer)
		}
	}
	if handler, ok := r.application.(gira.PeerHandler); ok {
		handler.OnPeerDelete(peer)
	}
	if peer.FullName == r.FullName {
		log.Errorw("delete myself???", "is_normal", self.isNormalUnregisterSelf)
		if self.isNormalUnregisterSelf {
			self.cancelFunc()
		} else {
			log.Errorw("重新注册自己")
			if err := self.registerSelf(r); err != nil {
				log.Errorw("register self fail", "error", err)
				return err
			}
		}
	}
	return nil
}

func (self *PeerRegistry) onPeerUpdate(r *Registry, peer *gira.Peer) error {
	log.Debugw("============ peer update ==================", "full_name", peer.FullName)
	for _, fw := range r.application.Frameworks() {
		if handler, ok := fw.(gira.PeerHandler); ok {
			handler.OnPeerUpdate(peer)
		}
	}
	if handler, ok := r.application.(gira.PeerHandler); ok {
		handler.OnPeerUpdate(peer)
	}
	return nil
}

func (self *PeerRegistry) onKvPut(r *Registry, kv *mvccpb.KeyValue) error {
	pats := strings.Split(string(kv.Key), "/")
	if len(pats) != 4 {
		log.Errorw("peer registry got a invalid peer", "kv", string(kv.Key))
		return gira.ErrInvalidPeer
	}
	fullName := pats[2]
	attrName := pats[3]
	name, serverId, err := explodeServerFullName(fullName)
	if err != nil {
		log.Errorw("peer registry got a invalid peer", "full_name", fullName)
		return err
	}
	attrValue := string(kv.Value)
	if lastValue, ok := self.Peers.Load(fullName); ok {
		lastPeer := lastValue.(*gira.Peer)
		if attrName == GRPC_KEY {
			if lastPeer.GrpcAddr == "" {
				// 新增节点
				lastPeer.GrpcAddr = attrValue
				log.Infow("etcd add peer", "full_name", fullName, "attr_value", attrValue)
				self.onPeerAdd(r, lastPeer)
			} else if attrValue != lastPeer.GrpcAddr {
				// 节点地址改变
				lastPeer.GrpcAddr = attrValue
				self.onPeerUpdate(r, lastPeer)
				log.Info("etcd peer", fullName, "address change from", lastPeer.GrpcAddr, "to", attrValue)
			} else {
				log.Info("etcd peer", fullName, "address not change", lastPeer.GrpcAddr)
			}
		} else {
			if lastAttrValue, ok := lastPeer.Kvs[attrName]; ok {
				if lastAttrValue != attrValue {
					log.Info("etcd peer", fullName, "attr change", attrName, lastAttrValue, "=>", attrValue)
				} else {
					log.Info("etcd peer", fullName, "attr not change", attrName, lastAttrValue)
				}
			} else {
				log.Info("etcd peer", fullName, "add attr", attrName, "=>", attrValue)
			}
			lastPeer.Kvs[attrName] = attrValue
		}
	} else {
		peer := &gira.Peer{
			Id:       serverId,
			Name:     name,
			FullName: fullName,
			Kvs:      make(map[string]string),
		}
		self.Peers.Store(fullName, peer)
		if attrName == GRPC_KEY {
			// 新增节点
			log.Info("etcd add peer", fullName, attrValue)
			peer.GrpcAddr = attrValue
			self.onPeerAdd(r, peer)
		} else {
			log.Info("etcd peer", fullName, "add attr", attrName, "=>", attrValue)
			peer.Kvs[attrName] = attrValue
		}
	}
	return nil
}

func (self *PeerRegistry) onKvDelete(r *Registry, kv *mvccpb.KeyValue) error {
	pats := strings.Split(string(kv.Key), "/")
	if len(pats) != 4 {
		log.Info("peer registry got a invalid peer", string(kv.Key))
		return gira.ErrInvalidPeer
	}
	fullName := pats[2]
	attrName := pats[3]
	if lastValue, ok := self.Peers.Load(fullName); ok {
		lastPeer := lastValue.(*gira.Peer)
		if attrName == GRPC_KEY {
			//删除节点
			log.Info("etcd remove peer", fullName, lastPeer.GrpcAddr)
			lastPeer.GrpcAddr = ""
			self.onPeerDelete(r, lastPeer)
		} else {
			if lastAttrValue, ok := lastPeer.Kvs[attrName]; ok {
				delete(lastPeer.Kvs, attrName)
				log.Info("etcd peer", fullName, "remove attr", attrName, "=>", lastAttrValue)
			} else {
				log.Info("etcd peer", fullName, "remove attr, but attr not found!!!!!", fullName, "=>", attrName)
			}
		}
	} else {
		log.Info("etcd peer", fullName, "remove attr, but peer not found!!!!!", fullName, "=>", attrName)
	}
	return nil
}

// 只增加节点，但不通知handler, 等notify再通知
func (r *PeerRegistry) onKvAdd(kv *mvccpb.KeyValue) error {
	pats := strings.Split(string(kv.Key), "/")
	if len(pats) != 4 {
		log.Info("peer registry got a invalid peer", string(kv.Key))
		return gira.ErrInvalidPeer
	}
	fullName := pats[2]
	attrName := pats[3]
	name, serverId, err := explodeServerFullName(fullName)
	if err != nil {
		log.Info("peer registry got a invalid peer", fullName)
		return err
	}
	attrValue := string(kv.Value)
	if lastValue, ok := r.Peers.Load(fullName); ok {
		lastPeer := lastValue.(*gira.Peer)
		if attrName == GRPC_KEY {
			if lastPeer.GrpcAddr == "" {
				// 新增节点
				lastPeer.GrpcAddr = attrValue
				log.Info("etcd add peer", fullName, attrValue)
			} else if attrValue != lastPeer.GrpcAddr {
				// 节点地址改变
				lastPeer.GrpcAddr = attrValue
				log.Info("etcd peer", fullName, "address change from", lastPeer.GrpcAddr, "to", attrValue)
			} else {
				log.Info("etcd peer", fullName, "address not change", lastPeer.GrpcAddr)
			}
		} else {
			if lastAttrValue, ok := lastPeer.Kvs[attrName]; ok {
				if lastAttrValue != attrValue {
					log.Info("etcd peer", fullName, "attr change", attrName, lastAttrValue, "=>", attrValue)
				} else {
					log.Info("etcd peer", fullName, "attr not change", attrName, lastAttrValue)
				}
			} else {
				log.Info("etcd peer", fullName, "add attr", attrName, "=>", attrValue)
			}
			lastPeer.Kvs[attrName] = attrValue
		}
	} else {
		peer := &gira.Peer{
			Id:       serverId,
			Name:     name,
			FullName: fullName,
			Kvs:      make(map[string]string),
		}
		r.Peers.Store(fullName, peer)
		if attrName == GRPC_KEY {
			peer.GrpcAddr = attrValue
			// 新增节点
			log.Info("etcd add peer", fullName, attrValue)
		} else {
			peer.Kvs[attrName] = attrValue
			log.Info("etcd peer", fullName, "add attr", attrName, "=>", attrValue)
		}
	}
	return nil
}

func (self *PeerRegistry) watchPeers(r *Registry) error {
	client := r.Client
	kv := clientv3.NewKV(client)
	var getResp *clientv3.GetResponse
	var err error
	if getResp, err = kv.Get(self.cancelCtx, self.Prefix, clientv3.WithPrefix()); err != nil {
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
		watchRespChan := watcher.Watch(self.cancelCtx, self.Prefix, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix(), clientv3.WithPrevKV())
		log.Info("etcd watch peer started", self.Prefix, watchStartRevision)
		for watchResp := range watchRespChan {
			// log.Info("etcd watch got events")
			for _, event := range watchResp.Events {
				switch event.Type {
				case mvccpb.PUT:
					// log.Info("etcd got put event")
					if err := self.onKvPut(r, event.Kv); err != nil {
						log.Info("peer registry put event error", err)
					}
				case mvccpb.DELETE:
					// log.Info("etcd got delete event")
					if err := self.onKvDelete(r, event.Kv); err != nil {
						log.Info("peer registry put event error", err)
					}
				}
			}
		}
		log.Info("peer registry watch shutdown")

		return nil
	})
	return nil
}

func (self *PeerRegistry) unregisterSelf(r *Registry) error {
	client := r.Client
	kv := clientv3.NewKV(client)
	ctx, cancelFunc := context.WithTimeout(self.cancelCtx, 10*time.Second)
	defer cancelFunc()
	log.Info("etcd unregister self", self.SelfPrefix)
	self.isNormalUnregisterSelf = true

	// txn.If(clientv3.Compare(clientv3.Value(key), "!=", value), clientv3.Compare(clientv3.CreateRevision(key), "!=", 0))

	var txnResp *clientv3.TxnResponse
	var err error
	txn := kv.Txn(ctx)
	key := fmt.Sprintf("%s%s", self.SelfPrefix, GRPC_KEY)
	value := r.Config.Address

	txn.If(clientv3.Compare(clientv3.Value(key), "!=", value), clientv3.Compare(clientv3.CreateRevision(key), "!=", 0)).
		Then(clientv3.OpGet(key)).
		Else(clientv3.OpGet(key), clientv3.OpDelete(self.SelfPrefix, clientv3.WithPrefix()))

	if txnResp, err = txn.Commit(); err != nil {
		log.Info("txn err", err)
		return err
	}

	if txnResp.Succeeded {
		// key 被其他程序占用着，不用管，直接退出
		log.Info(" key 被其他程序占用着，不用管，直接退出")
		self.cancelFunc()
	} else {
		if len(txnResp.Responses[0].GetResponseRange().Kvs) == 0 {
			// key 没被任何程序占用，不用管，直接退出
			log.Info("key 没被任何程序占用，不用管，直接退出")
			self.cancelFunc()
		} else {
			// 等监听到清理完成后再退出
			log.Info("等监听到清理完成后再退出")

		}
	}
	return nil
}

func (self *PeerRegistry) registerSelf(r *Registry) error {
	client := r.Client
	var err error
	var lease clientv3.Lease
	var leaseID clientv3.LeaseID
	if r.Config.LeaseTimeout > 0 {
		// 申请一个租约 lease
		lease = clientv3.Lease(client)
		var leaseGrantResp *clientv3.LeaseGrantResponse
		// 申请一个5s的租约
		if leaseGrantResp, err = lease.Grant(self.cancelCtx, 5); err != nil {
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
		txn := kv.Txn(self.cancelCtx)
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
		if keepRespChan, err = lease.KeepAlive(self.cancelCtx, leaseID); err != nil {
			log.Info("自动续租失败", err)
			return err
		}
		//判断续约应答的协程
		r.application.Go(func() error {
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
				case <-r.application.Done():
					break
				}
			}
		END:
			return nil
		})
	}
	return nil
}
