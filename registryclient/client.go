package registryclient

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/options/service_options"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

const GRPC_KEY string = "grpc"

type RegistryClient struct {
	config     gira.EtcdClientConfig
	appId      int32
	fullName   string // 节点全名
	name       string // 节点名
	client     *clientv3.Client
	ctx        context.Context
	cancelFunc context.CancelFunc

	peerRegistry    *peer_registry
	playerRegistry  *player_registry
	serviceRegistry *service_registry
}

func (r *RegistryClient) StartAsClient() error {
	return nil
}

func NewConfigRegistryClient(ctx context.Context, config *gira.EtcdClientConfig) (*RegistryClient, error) {
	r := &RegistryClient{
		config:   *config,
		fullName: facade.GetAppFullName(),
		appId:    facade.GetAppId(),
		name:     facade.GetAppType(),
	}
	r.ctx, r.cancelFunc = context.WithCancel(ctx)
	// 配置endpoints
	endpoints := make([]string, 0)
	for _, v := range r.config.Endpoints {
		endpoints = append(endpoints, fmt.Sprintf("http://%s:%d", v.Host, v.Port))
	}
	var client *clientv3.Client
	var err error
	c := clientv3.Config{
		Endpoints:   endpoints,                                         // 节点信息
		DialTimeout: time.Duration(r.config.DialTimeout) * time.Second, // 超时时间
		DialOptions: []grpc.DialOption{grpc.WithBlock()},               // 使用阻塞模式，确认启动时etcd是可用的
		Username:    r.config.Username,
		Password:    r.config.Password,
		Context:     r.ctx,
	}
	// 建立连接
	if client, err = clientv3.New(c); err != nil {
		log.Errorw("connect to etcd fail", "error", err)
		return nil, err
	}
	r.client = client
	// log.Infow("connect registry success", "endpoints", endpoints)
	if v, err := newConfigPeerRegistry(r); err != nil {
		return nil, err
	} else {
		r.peerRegistry = v
	}
	if v, err := newConfigPlayerRegistry(r); err != nil {
		return nil, err
	} else {
		r.playerRegistry = v
	}
	if v, err := newConfigServiceRegistry(r); err != nil {
		return nil, err
	} else {
		r.serviceRegistry = v
	}
	return r, nil
}

// 查找玩家
func (r *RegistryClient) WhereIsUser(userId string) (*gira.Peer, error) {
	return r.playerRegistry.WhereIsUser(r, userId)
}

// 查找服务
func (r *RegistryClient) WhereIsService(serviceName string, opt ...service_options.WhereOption) ([]*gira.Peer, error) {
	return r.serviceRegistry.WhereIsService(r, serviceName, opt...)
}

func (r *RegistryClient) NewServiceName(serviceName string, opt ...service_options.RegisterOption) string {
	return r.serviceRegistry.NewServiceName(r, serviceName, opt...)
}

func (r *RegistryClient) GetPeer(fullName string) *gira.Peer {
	return r.peerRegistry.getPeer(r, fullName)
}

func explodeServerFullName(fullName string) (name string, id int32, err error) {
	pats := strings.Split(string(fullName), "_")
	if len(pats) != 4 {
		err = gira.ErrInvalidPeer
		return
	}
	name = pats[0]
	var v int
	if v, err = strconv.Atoi(pats[3]); err == nil {
		id = int32(v)
	}
	return
}
