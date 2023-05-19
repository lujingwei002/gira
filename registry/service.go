package registry

///
///
///
/// 注册表结构:
///   /unique_service/<<ServiceName>> => <<AppFullName>>
///   /service/<<ServiceGroupName>>/<<ServiceName>> => <<AppFullName>>
///   /peer_service/<<AppFullName>>/<<ServiceUniqueName>> => <<AppFullName>>
///
import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/options/service_options"
	mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type group_services struct {
	services sync.Map
}

type service_registry struct {
	peerPrefix         string   // /peer_service/<<AppFullName>><<ServiceName>>/			根据服务全名查找全部服务
	servicePrefix      string   // /service/<<ServiceName>>/      							可以根据服务名查找当前所在的服
	services           sync.Map // 全部service
	groupServices      sync.Map // service分组
	selfServices       sync.Map // 本节点上的service
	ctx                context.Context
	cancelFunc         context.CancelFunc
	watchStartRevision int64
}

func newConfigServiceRegistry(r *Registry) (*service_registry, error) {
	self := &service_registry{
		servicePrefix: "/service/",
		peerPrefix:    fmt.Sprintf("/peer_service/%s/", r.fullName),
	}
	return self, nil
}

func (self *service_registry) onStop(r *Registry) error {
	log.Debug("service registry on stop")
	if err := self.unregisterServices(r); err != nil {
		log.Info(err)
	}
	return nil
}

func (self *service_registry) OnStart(r *Registry) error {
	cancelCtx, cancelFunc := context.WithCancel(r.ctx)
	self.ctx = cancelCtx
	self.cancelFunc = cancelFunc
	if err := self.initServices(r); err != nil {
		return err
	}
	return nil
}

func (self *service_registry) Serve(r *Registry) error {
	return self.watchServices(r)
}

func (self *service_registry) notify(r *Registry) error {
	self.services.Range(func(k any, v any) bool {
		service := v.(*gira.ServiceName)
		self.onServiceAdd(r, service)
		return true
	})
	return nil
}

// on service add callback
func (self *service_registry) onServiceAdd(r *Registry, service *gira.ServiceName) {
	if r.isNotify == 0 {
		return
	}
	log.Infow("service registry on service add", "service_name", service.FullName, "peer", service.Peer.FullName)
	for _, fw := range r.application.Frameworks() {
		if handler, ok := fw.(gira.ServiceWatchHandler); ok {
			handler.OnServiceAdd(service)
		}
	}
	if handler, ok := r.application.(gira.ServiceWatchHandler); ok {
		handler.OnServiceAdd(service)
	}
}

// on service delete callback
func (self *service_registry) onServiceDelete(r *Registry, service *gira.ServiceName) {
	log.Infow("service registry on service delete", "service_name", service.FullName, "peer", service.Peer.FullName)
	for _, fw := range r.application.Frameworks() {
		if handler, ok := fw.(gira.ServiceWatchHandler); ok {
			handler.OnServiceDelete(service)
		}
	}
	if handler, ok := r.application.(gira.ServiceWatchHandler); ok {
		handler.OnServiceDelete(service)
	}
}

// func (self *service_registry) onServiceUpdate(r *Registry, service *gira.Service) {
// 	log.Infow("service update", "service_name", service.UniqueName)
// 	for _, fw := range r.application.Frameworks() {
// 		if handler, ok := fw.(gira.ServiceWatchHandler); ok {
// 			handler.OnServiceUpdate(service)
// 		}
// 	}
// 	if handler, ok := r.application.(gira.ServiceWatchHandler); ok {
// 		handler.OnServiceUpdate(service)
// 	}
// }

func (self *service_registry) onKvAdd(r *Registry, kv *mvccpb.KeyValue) error {
	pats := strings.Split(string(kv.Key), "/")
	var groupName string
	var serviceName string
	var name string
	if len(pats) == 4 {
		serviceName = fmt.Sprintf("%s/%s", pats[2], pats[3])
		groupName = pats[2]
		name = pats[3]
	} else if len(pats) == 3 {
		serviceName = pats[2]
	} else {
		log.Warnw("service registry got a invalid key", "key", string(kv.Key))
		return gira.ErrInvalidService
	}
	value := string(kv.Value)
	if lastValue, ok := self.services.Load(serviceName); ok {
		lastService := lastValue.(*gira.ServiceName)
		log.Warnw("service registry on kv add, but already exist", "service_name", serviceName, "peer", value, "last_peer", lastService.Peer.FullName)
	} else {
		// 新增service
		log.Infow("service registry on kv add", "service_name", serviceName, "peer", value)
		peer := r.GetPeer(value)
		if peer == nil {
			log.Warnw("service registry on kv add, but peer not found", "service_name", serviceName, "peer", value)
			return gira.ErrPeerNotFound
		}
		service := &gira.ServiceName{
			FullName:  serviceName,
			GroupName: groupName,
			Peer:      peer,
			Name:      name,
		}
		if peer == r.SelfPeer() {
			service.IsSelf = true
		}
		if len(service.GroupName) > 0 {
			service.IsGroup = true
		}
		self.services.Store(serviceName, service)
		if service.IsGroup {
			groupServices := &group_services{}
			if v, _ := self.groupServices.LoadOrStore(service.GroupName, groupServices); true {
				groupServices = v.(*group_services)
				groupServices.services.Store(serviceName, service)
			}
		}
		if service.IsSelf {
			self.selfServices.Store(serviceName, service)
		}
		self.onServiceAdd(r, service)
	}
	return nil
}

func (self *service_registry) onKvDelete(r *Registry, kv *mvccpb.KeyValue) error {
	pats := strings.Split(string(kv.Key), "/")
	var serviceName string
	if len(pats) == 4 {
		serviceName = fmt.Sprintf("%s/%s", pats[2], pats[3])
	} else if len(pats) == 3 {
		serviceName = pats[2]
	} else {
		log.Warnw("service registry got a invalid key", "key", string(kv.Key))
		return gira.ErrInvalidService
	}
	// value := string(kv.Value) value没有值
	if lastValue, ok := self.services.Load(serviceName); ok {
		lastService := lastValue.(*gira.ServiceName)
		log.Infow("service registry remove service", "service_name", serviceName, "last_peer", lastService.Peer.FullName)
		self.services.Delete(serviceName)
		self.onServiceDelete(r, lastService)
		if lastService.IsGroup {
			if v, ok := self.groupServices.Load(lastService.GroupName); ok {
				groupServices := v.(*group_services)
				groupServices.services.Delete(serviceName)
			}
		}
		if lastService.IsSelf {
			self.selfServices.Delete(serviceName)
		}
	} else {
		log.Warnw("service registry remove service, but service not found", "service_name", serviceName)
	}
	return nil
}

func (self *service_registry) initServices(r *Registry) error {
	client := r.client
	kv := clientv3.NewKV(client)
	var getResp *clientv3.GetResponse
	var err error
	// 删除自身之前注册，没清理干净的服务
	if getResp, err = kv.Get(self.ctx, self.peerPrefix, clientv3.WithPrefix()); err != nil {
		return err
	}
	for _, v := range getResp.Kvs {
		pats := strings.Split(string(v.Key), "/")
		var serviceName string
		if len(pats) == 5 {
			serviceName = fmt.Sprintf("%s/%s", pats[3], pats[4])
		} else if len(pats) == 4 {
			serviceName = pats[3]
		}
		txn := kv.Txn(self.ctx)
		serviceKey := fmt.Sprintf("%s%s", self.servicePrefix, serviceName)
		peerKey := fmt.Sprintf("%s%s", self.peerPrefix, serviceName)
		txn.If(clientv3.Compare(clientv3.CreateRevision(serviceKey), "!=", 0)).
			Then(clientv3.OpDelete(peerKey), clientv3.OpDelete(serviceKey))
		var txnResp *clientv3.TxnResponse
		if txnResp, err = txn.Commit(); err != nil {
			log.Errorw("service registry commit fail", "error", err)
			return err
		}
		if txnResp.Succeeded {
			log.Infow("service registry unregister", "service_key", serviceKey, "peer_key", peerKey, "service_name", serviceName)
		} else {
			log.Warnw("service registry unregister", "service_key", serviceKey, "peer_key", peerKey, "service_name", serviceName)
		}
	}
	// 初始化服务
	if getResp, err = kv.Get(self.ctx, self.servicePrefix, clientv3.WithPrefix()); err != nil {
		return err
	}
	for _, kv := range getResp.Kvs {
		if err := self.onKvAdd(r, kv); err != nil {
			return err
		}
	}
	self.watchStartRevision = getResp.Header.Revision + 1
	return nil
}

func (self *service_registry) watchServices(r *Registry) error {
	client := r.client
	watchStartRevision := self.watchStartRevision
	watcher := clientv3.NewWatcher(client)
	// r.application.Go(func() error {
	watchRespChan := watcher.Watch(self.ctx, self.servicePrefix, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix(), clientv3.WithPrevKV())
	log.Infow("service registry started", "service_prefix", self.servicePrefix, "watch_start_revision", watchStartRevision)
	for watchResp := range watchRespChan {
		// log.Info("etcd watch got events")
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				// log.Info("etcd got put event")
				if err := self.onKvAdd(r, event.Kv); err != nil {
					log.Warnw("service registry put event fail", "error", err)
				}
			case mvccpb.DELETE:
				// log.Info("etcd got delete event")
				if err := self.onKvDelete(r, event.Kv); err != nil {
					log.Warnw("service registry put event fail", "error", err)
				}
			}
		}
	}
	log.Info("service registry watch exit")
	return nil
	// })
}

func (self *service_registry) unregisterServices(r *Registry) error {
	client := r.client
	kv := clientv3.NewKV(client)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	log.Infow("service registry unregister", "peer_prefix", self.peerPrefix)

	var txnResp *clientv3.TxnResponse
	var err error
	self.selfServices.Range(func(serviceName any, v any) bool {
		txn := kv.Txn(ctx)
		serviceKey := fmt.Sprintf("%s%s", self.servicePrefix, serviceName)
		peerKey := fmt.Sprintf("%s%s", self.peerPrefix, serviceName)
		txn.If(clientv3.Compare(clientv3.CreateRevision(serviceKey), "!=", 0)).
			Then(clientv3.OpDelete(peerKey), clientv3.OpDelete(serviceKey))
		if txnResp, err = txn.Commit(); err != nil {
			log.Errorw("service registry commit fail", "error", err)
			return true
		}
		if txnResp.Succeeded {
			log.Infow("service registry unregister", "peer_key", peerKey)
		} else {
			log.Warnw("service registry unregister", "peer_key", peerKey)
		}
		return true
	})
	return nil
}

func (self *service_registry) NewServiceName(r *Registry, serviceName string, opt ...service_options.RegisterOption) string {
	opts := service_options.RegisterOptions{}
	for _, v := range opt {
		v.ConfigRegisterOption(&opts)
	}
	if opts.AsAppService {
		serviceName = fmt.Sprintf("%s/%s_%d", serviceName, serviceName, r.appId)
	}
	return serviceName
}

// 注册服务
func (self *service_registry) RegisterService(r *Registry, serviceName string, opt ...service_options.RegisterOption) (*gira.Peer, error) {
	serviceName = self.NewServiceName(r, serviceName, opt...)
	client := r.client
	serviceKey := fmt.Sprintf("%s%s", self.servicePrefix, serviceName)
	peerKey := fmt.Sprintf("%s%s", self.peerPrefix, serviceName)
	kv := clientv3.NewKV(client)
	var err error
	var txnResp *clientv3.TxnResponse
	txn := kv.Txn(self.ctx)
	log.Infow("service registry register", "service_name", serviceName, "peer_key", peerKey, "service_key", serviceKey)
	txn.If(clientv3.Compare(clientv3.CreateRevision(serviceKey), "=", 0)).
		Then(clientv3.OpPut(peerKey, r.fullName), clientv3.OpPut(serviceKey, r.fullName)).
		Else(clientv3.OpGet(serviceKey))
	if txnResp, err = txn.Commit(); err != nil {
		log.Errorw("service registry commit fail", "error", err)
		return nil, err
	}
	if txnResp.Succeeded {
		return nil, nil
	} else {
		log.Warnw("service registry register fail", "service_name", serviceName, "locked_by", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		fullName := string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
		peer := r.GetPeer(fullName)
		if peer == nil {
			return nil, gira.ErrServiceLocked.Trace()
		}
		return peer, gira.ErrServiceLocked.Trace()
	}
}

// 解锁
func (self *service_registry) UnregisterService(r *Registry, serviceName string) (*gira.Peer, error) {
	client := r.client
	serviceKey := fmt.Sprintf("%s%s", self.servicePrefix, serviceName)
	peerKey := fmt.Sprintf("%s%s", self.peerPrefix, serviceName)
	kv := clientv3.NewKV(client)
	var err error
	var txnResp *clientv3.TxnResponse
	txn := kv.Txn(self.ctx)
	log.Infow("service registry", "peer_key", peerKey, "service_key", serviceKey)
	txn.If(clientv3.Compare(clientv3.Value(serviceKey), "=", r.fullName), clientv3.Compare(clientv3.CreateRevision(serviceKey), "!=", 0)).
		Then(clientv3.OpDelete(peerKey), clientv3.OpDelete(serviceKey)).
		Else(clientv3.OpGet(serviceKey))
	if txnResp, err = txn.Commit(); err != nil {
		log.Errorw("service registry commit fail", "error", err)
		return nil, err
	}
	if txnResp.Succeeded {
		log.Infow("service registry unregister", "service_name", serviceName)
		return nil, nil
	} else {
		var fullName string
		if len(txnResp.Responses) > 0 && len(txnResp.Responses[0].GetResponseRange().Kvs) > 0 {
			fullName = string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
		}
		log.Warnw("service registry unregister fail", "service_name", serviceName, "locked_by", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		peer := r.GetPeer(fullName)
		if peer == nil {
			return nil, gira.ErrServiceLocked.Trace()
		}
		return peer, gira.ErrServiceLocked.Trace()
	}
}

// 查找服务位置
func (self *service_registry) WhereIsService(r *Registry, serviceName string, opt ...service_options.WhereOption) (peers []*gira.Peer, err error) {
	opts := service_options.WhereOptions{}
	for _, v := range opt {
		v.ConfigWhereOption(&opts)
	}
	pats := strings.Split(serviceName, "/")
	if len(pats) != 1 && len(pats) != 2 {
		err = gira.ErrInvalidArgs.Trace()
		return
	}
	if len(pats) <= 1 {
		peers = make([]*gira.Peer, 0)
		if value, ok := self.services.Load(serviceName); ok {
			service := value.(*gira.ServiceName)
			peers = append(peers, service.Peer)
		}
		return
	} else if pats[1] == "" {
		// 全部
		peers = make([]*gira.Peer, 0)
		multicastCount := opts.MaxCount
		if v, ok := self.groupServices.Load(pats[0]); ok {
			group := v.(*group_services)
			group.services.Range(func(key any, value any) bool {
				service := value.(*gira.ServiceName)
				peers = append(peers, service.Peer)
				// 多播指定数量
				if multicastCount > 0 {
					multicastCount--
					if multicastCount <= 0 {
						return false
					}
				}
				return true
			})
		}
		return
	} else {
		// 精确查找
		peers = make([]*gira.Peer, 0)
		if v, ok := self.groupServices.Load(pats[0]); ok {
			group := v.(*group_services)
			group.services.Range(func(key any, value any) bool {
				service := value.(*gira.ServiceName)
				if service.Name == pats[1] {
					peers = append(peers, service.Peer)
				}
				return true
			})
		}
		return
	}
}

func (self *service_registry) listServiceKvs(r *Registry) (kvs map[string][]string, err error) {
	client := r.client
	kv := clientv3.NewKV(client)
	var getResp *clientv3.GetResponse
	if getResp, err = kv.Get(self.ctx, self.servicePrefix, clientv3.WithPrefix()); err != nil {
		return
	}
	kvs = make(map[string][]string)
	for _, kv := range getResp.Kvs {
		pats := strings.Split(string(kv.Key), "/")
		var serviceName string
		if len(pats) == 4 {
			serviceName = fmt.Sprintf("%s/%s", pats[2], pats[3])
		} else if len(pats) == 3 {
			serviceName = pats[2]
		}
		peerFullName := string(kv.Value)
		kvs[peerFullName] = append(kvs[peerFullName], serviceName)
	}
	return
}
