package registry

///
///
///
/// 注册表结构:
///   /unique_service/<<ServiceName>> => <<AppFullName>>
///   /service/<<ServiceGroupName>>/<<ServiceName>> => <<AppFullName>>
///   /local_service/<<AppFullName>>/<<ServiceUniqueName>> => <<AppFullName>>
///
import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/registry/service/options"
	mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type group_services struct {
	services sync.Map
}

type service_registry struct {
	localPrefix   string // /local_service/<<AppFullName>><<ServiceName>>/			根据服务全名查找全部服务
	servicePrefix string // /service/<<ServiceName>>/      							可以根据服务名查找当前所在的服
	services      sync.Map
	localServices sync.Map
	groupServices sync.Map
	ctx           context.Context
	cancelFunc    context.CancelFunc
}

func newConfigServiceRegistry(r *Registry) (*service_registry, error) {
	self := &service_registry{
		servicePrefix: "/service/",
		localPrefix:   fmt.Sprintf("/local_service/%s/", r.fullName),
	}
	self.ctx, self.cancelFunc = context.WithCancel(r.ctx)
	// r.application.Go(func() error {
	// 	select {
	// 	case <-r.application.Done():
	// 		{
	// 			log.Info("service registry recv down")
	// 		}
	// 	}
	// 	if err := self.unregisterSelf(r); err != nil {
	// 		log.Info(err)
	// 	}
	// 	self.cancelFunc()
	// 	return nil
	// })
	return self, nil
}

func (self *service_registry) onDestory(r *Registry) error {
	if err := self.unregisterSelf(r); err != nil {
		log.Info(err)
	}
	self.cancelFunc()
	return nil
}

func (self *service_registry) onStart(r *Registry) error {
	if err := self.watchServices(r); err != nil {
		return err
	}
	if err := self.notify(r); err != nil {
		return err
	}
	return nil
}

func (self *service_registry) notify(r *Registry) error {
	self.services.Range(func(k any, v any) bool {
		service := v.(*gira.Service)
		self.onServiceAdd(r, service)
		return true
	})
	return nil
}

func (self *service_registry) onServiceAdd(r *Registry, service *gira.Service) {
	log.Infow("service registry on service add", "service_name", service.FullName)
	for _, fw := range r.application.Frameworks() {
		if handler, ok := fw.(gira.ServiceWatchHandler); ok {
			handler.OnServiceAdd(service)
		}
	}
	if handler, ok := r.application.(gira.ServiceWatchHandler); ok {
		handler.OnServiceAdd(service)
	}
}

func (self *service_registry) onServiceDelete(r *Registry, service *gira.Service) {
	log.Infow("service registry on service delete", "service_name", service.FullName)
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

func (self *service_registry) onKvPut(r *Registry, kv *mvccpb.KeyValue) error {
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
		lastService := lastValue.(*gira.Service)
		log.Warnw("service registry add service, but already exist", "service_name", serviceName, "peer", value, "last_peer", lastService.Peer.FullName)
	} else {
		// 新增service
		log.Infow("service registry add service", "service_name", serviceName, "peer", value)
		peer := r.GetPeer(value)
		if peer == nil {
			log.Warnw("service registry add service, but peer not found", "service_name", serviceName, "peer", value)
			return gira.ErrPeerNotFound
		}
		service := &gira.Service{
			FullName:  serviceName,
			GroupName: groupName,
			Peer:      peer,
			Name:      name,
		}
		if peer == r.SelfPeer() {
			service.IsLocal = true
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
		if service.IsLocal {
			self.localServices.Store(serviceName, service)
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
		lastService := lastValue.(*gira.Service)
		log.Infow("service registry remove service", "service_name", serviceName, "last_peer", lastService.Peer.FullName)
		self.services.Delete(serviceName)
		self.onServiceDelete(r, lastService)
		if lastService.IsGroup {
			if v, ok := self.groupServices.Load(lastService.GroupName); ok {
				groupServices := v.(*group_services)
				groupServices.services.Delete(serviceName)
			}
		}
		if lastService.IsLocal {
			self.localServices.Delete(serviceName)
		}
	} else {
		log.Warnw("service registry remove service, but service not found", "service_name", serviceName)
	}
	return nil
}

func (self *service_registry) onKvAdd(r *Registry, kv *mvccpb.KeyValue) error {
	pats := strings.Split(string(kv.Key), "/")
	var serviceName string
	var groupName string
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
		lastService := lastValue.(*gira.Service)
		log.Warnw("service registry add service, but already exist", "service_name", serviceName, "locked_by", lastService.Peer.FullName)
	} else {
		peer := r.GetPeer(value)
		if peer == nil {
			log.Warnw("service registry add service, but peer not found", "service_name", serviceName, "peer", value)
			return gira.ErrPeerNotFound
		}
		service := &gira.Service{
			FullName:  serviceName,
			GroupName: groupName,
			Name:      name,
			Peer:      peer,
		}
		self.services.Store(serviceName, service)
		if peer == r.SelfPeer() {
			service.IsLocal = true
		}
		if len(service.GroupName) > 0 {
			service.IsGroup = true
		}
		if service.IsGroup {
			groupServices := &group_services{}
			if v, _ := self.groupServices.LoadOrStore(service.GroupName, groupServices); true {
				groupServices = v.(*group_services)
				groupServices.services.Store(serviceName, service)
			}
		}
		if service.IsLocal {
			self.localServices.Store(serviceName, service)
		}
		log.Infow("service registry add service", "service_name", serviceName, "peer", value)
	}
	return nil
}

func (self *service_registry) watchServices(r *Registry) error {
	client := r.client
	kv := clientv3.NewKV(client)
	var getResp *clientv3.GetResponse
	var err error
	// 删除自身之前注册，没清理干净的服务
	if getResp, err = kv.Get(self.ctx, self.localPrefix, clientv3.WithPrefix()); err != nil {
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
		log.Println(pats)
		txn := kv.Txn(self.ctx)
		serviceKey := fmt.Sprintf("%s%s", self.servicePrefix, serviceName)
		localKey := fmt.Sprintf("%s%s", self.localPrefix, serviceName)
		txn.If(clientv3.Compare(clientv3.CreateRevision(serviceKey), "!=", 0)).
			Then(clientv3.OpDelete(localKey), clientv3.OpDelete(serviceKey))
		var txnResp *clientv3.TxnResponse
		if txnResp, err = txn.Commit(); err != nil {
			log.Errorw("service registry commit fail", "error", err)
			return err
		}
		if txnResp.Succeeded {
			log.Infow("service registry unregister self", "serviceKey", serviceKey, "local_key", localKey, "service_name", serviceName)
		} else {
			log.Warnw("service registry unregister self", "serviceKey", serviceKey, "local_key", localKey, "service_name", serviceName)
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
	watchStartRevision := getResp.Header.Revision + 1
	watcher := clientv3.NewWatcher(client)
	r.application.Go(func() error {
		watchRespChan := watcher.Watch(self.ctx, self.servicePrefix, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix(), clientv3.WithPrevKV())
		log.Infow("service registry started", "service_prefix", self.servicePrefix, "watch_start_revision", watchStartRevision)
		for watchResp := range watchRespChan {
			// log.Info("etcd watch got events")
			for _, event := range watchResp.Events {
				switch event.Type {
				case mvccpb.PUT:
					// log.Info("etcd got put event")
					if err := self.onKvPut(r, event.Kv); err != nil {
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
		log.Info("service registry watch shutdown")
		return nil
	})
	return nil
}

func (self *service_registry) unregisterSelf(r *Registry) error {
	client := r.client
	kv := clientv3.NewKV(client)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	log.Infow("service registry unregister self", "local_prefix", self.localPrefix)

	var txnResp *clientv3.TxnResponse
	var err error
	self.localServices.Range(func(serviceName any, v any) bool {
		txn := kv.Txn(ctx)
		serviceKey := fmt.Sprintf("%s%s", self.servicePrefix, serviceName)
		localKey := fmt.Sprintf("%s%s", self.localPrefix, serviceName)
		txn.If(clientv3.Compare(clientv3.CreateRevision(serviceKey), "!=", 0)).
			Then(clientv3.OpDelete(localKey), clientv3.OpDelete(serviceKey))
		if txnResp, err = txn.Commit(); err != nil {
			log.Errorw("service registry commit fail", "error", err)
			return true
		}
		if txnResp.Succeeded {
			log.Infow("service registry unregister self", "serviceKey", serviceKey, "local_key", localKey, "service_name", serviceName)
		} else {
			log.Warnw("service registry unregister self", "serviceKey", serviceKey, "local_key", localKey, "service_name", serviceName)
		}
		return true
	})
	return nil
}

// 注册服务
func (self *service_registry) RegisterService(r *Registry, serviceName string, opt ...options.RegisterOption) (*gira.Peer, error) {
	opts := options.RegisterOptions{}
	for _, v := range opt {
		v.ConfigRegisterOption(&opts)
	}
	if opts.AsGroup {
		if opts.CatAppId {
			serviceName = fmt.Sprintf("%s/%s_%d", serviceName, serviceName, r.appId)
		}
	}
	client := r.client
	serviceKey := fmt.Sprintf("%s%s", self.servicePrefix, serviceName)
	localKey := fmt.Sprintf("%s%s", self.localPrefix, serviceName)
	kv := clientv3.NewKV(client)
	var err error
	var txnResp *clientv3.TxnResponse
	txn := kv.Txn(self.ctx)
	log.Infow("service registry register service", "local_key", localKey, "service_key", serviceKey)
	txn.If(clientv3.Compare(clientv3.CreateRevision(serviceKey), "=", 0)).
		Then(clientv3.OpPut(localKey, r.fullName), clientv3.OpPut(serviceKey, r.fullName)).
		Else(clientv3.OpGet(serviceKey))
	if txnResp, err = txn.Commit(); err != nil {
		log.Errorw("service registry commit fail", "error", err)
		return nil, err
	}
	if txnResp.Succeeded {
		log.Infow("service registry register success", "service_name", serviceName)
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
	localKey := fmt.Sprintf("%s%s", self.localPrefix, serviceName)
	kv := clientv3.NewKV(client)
	var err error
	var txnResp *clientv3.TxnResponse
	txn := kv.Txn(self.ctx)
	log.Infow("service registry", "local_key", localKey, "service_key", serviceKey)
	txn.If(clientv3.Compare(clientv3.Value(serviceKey), "=", r.fullName), clientv3.Compare(clientv3.CreateRevision(serviceKey), "!=", 0)).
		Then(clientv3.OpDelete(localKey), clientv3.OpDelete(serviceKey)).
		Else(clientv3.OpGet(serviceKey))
	if txnResp, err = txn.Commit(); err != nil {
		log.Errorw("service registry commit fail", "error", err)
		return nil, err
	}
	if txnResp.Succeeded {
		log.Infow("service registry unregister success", "service_name", serviceName)
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
func (self *service_registry) WhereIsService(r *Registry, serviceName string, opt ...options.WhereOption) (peers []*gira.Peer, err error) {
	opts := options.WhereOptions{}
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
			service := value.(*gira.Service)
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
				service := value.(*gira.Service)
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
				service := value.(*gira.Service)
				if service.Name == pats[1] {
					peers = append(peers, service.Peer)
				}
				return true
			})
		}
		return
	}
}
