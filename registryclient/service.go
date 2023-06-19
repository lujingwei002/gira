package registryclient

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

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/options/service_options"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type service_registry struct {
	servicePrefix string // /service/<<ServiceName>>/      							可以根据服务名查找当前所在的服
	ctx           context.Context
	cancelFunc    context.CancelFunc
}

func newConfigServiceRegistry(r *RegistryClient) (*service_registry, error) {
	ctx, cancelFunc := context.WithCancel(r.ctx)
	self := &service_registry{
		servicePrefix: "/service/",
		ctx:           ctx,
		cancelFunc:    cancelFunc,
	}
	return self, nil
}

func (self *service_registry) NewServiceName(r *RegistryClient, serviceName string, opt ...service_options.RegisterOption) string {
	opts := service_options.RegisterOptions{}
	for _, v := range opt {
		v.ConfigRegisterOption(&opts)
	}
	if opts.AsAppService {
		serviceName = fmt.Sprintf("%s/%s_%d", serviceName, serviceName, r.appId)
	}
	return serviceName
}

// 查找服务位置
func (self *service_registry) WhereIsService(r *RegistryClient, serviceName string, opt ...service_options.WhereOption) (peers []*gira.Peer, err error) {
	opts := service_options.WhereOptions{}
	for _, v := range opt {
		v.ConfigWhereOption(&opts)
	}
	var getOpts []clientv3.OpOption

	client := r.client
	kv := clientv3.NewKV(client)
	var getResp *clientv3.GetResponse

	key := fmt.Sprintf("%s%s", self.servicePrefix, serviceName)
	if getResp, err = kv.Get(self.ctx, key, getOpts...); err != nil {
		return
	}
	for _, kv := range getResp.Kvs {
		peer := r.GetPeer(string(kv.Value))
		if peer != nil {
			peers = append(peers, peer)
		}
	}
	return

}
